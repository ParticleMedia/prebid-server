package exchange

import (
	"encoding/json"
	"time"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
	"github.com/prebid/openrtb/v20/openrtb2"
	"github.com/prebid/prebid-server/v3/openrtb_ext"
	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

type MSBExt struct {
	MSB MSBConfig `json:"msb"`
}

type MSBConfig struct {
	LastPeek MSBLastPeekConfig `json:"last_peek"`
}

type MSBLastPeekConfig struct {
	// start peek available bids after this
	PeekStartTimeMilliSeconds int64 `json:"peek_start_time_miliseconds"`
	// list of bidder names that are in last peek tier
	PeekBidders []string `json:"peek_bidders"`
}

/*
MSB feature controller(req.ext) example:
{
    "ext": {
        "msb": {
            "last_peek": {
                "peek_start_time_miliseconds": 900,
                "peek_bidders": ["msp_google", "msp_nova"]
            }
        }
    }
}
MSP server is response for adding above MSB info to requests for certain traffic/placement/exps.
In the above example
	there are two peek tiers for bidders:
	1. last peek tier (bidders specified in peek_bidders)
	2. the rest(normal tier)
	after all reponses are ready for normal tier or timeout=peek_start_time_miliseconds(900ms), for bidders in last peek tier peek available responses from normal tier,
	get max_available_bid_prices_for_normal_tier and the winner bidder, then write peek_winner_price and peek_winner_bidder to req.imp.ext.bidder
*/

type MSPFloor struct {
	Floor            float64 `json:"floor"`
	PeekWinnerBidder string  `json:"peek_winner_bidder"`
	PeekWinnerPrice  float64 `json:"peek_winner_price"`
}

var mspBidders = map[openrtb_ext.BidderName]int{
	openrtb_ext.BidderMspGoogle:  1,
	openrtb_ext.BidderMspFbAlpha: 1,
	openrtb_ext.BidderMspFbBeta:  1,
	openrtb_ext.BidderMspFbGamma: 1,
	openrtb_ext.BidderMspNova:    1,
}

// peek max available(ready) bid price and winner bidder from channel without comsuming
func peekChannelAvailableMaxBidPriceWithinTimeout(peekTier string, chBids chan *bidResponseWrapper, lastPeekConfig MSBLastPeekConfig, totalNormalBidders int) (float64, string) {
	timeout := time.After(time.Duration(lastPeekConfig.PeekStartTimeMilliSeconds) * time.Millisecond)
	maxPrice := 0.0
	winnerBidder := ""
	hasData := true
	peekedRespList := []*bidResponseWrapper{}
	availableBidders := []string{}
	// keep consuming message from channel until all normal bidder requests are collected or timeout reaches
	for hasData {
		select {
		case resp, ok := <-chBids:
			if !ok {
				hasData = false
			} else {
				for _, bids := range resp.adapterSeatBids {
					for _, bid := range bids.Bids {
						if bid.Bid != nil {
							if bid.Bid.Price > maxPrice {
								maxPrice = bid.Bid.Price
								winnerBidder = resp.bidder.String()
							}
						}
					}
				}
				peekedRespList = append(peekedRespList, resp)
				availableBidders = append(availableBidders, resp.bidder.String())
				if len(peekedRespList) == totalNormalBidders {
					hasData = false
				}
			}

		case <-timeout:
			hasData = false
		}

	}
	// push message back
	for _, resp := range peekedRespList {
		chBids <- resp
	}
	glog.Infof("MSB tier %s, peeked from available bidders %v, current max bid price: %f, winner bidder: %s", peekTier, availableBidders, maxPrice, winnerBidder)
	return maxPrice, winnerBidder
}

func mspUpdateLastPeekBiddersRequest(
	chBids chan *bidResponseWrapper,
	lastPeekBidderRequests []BidderRequest,
	lastPeekConfig MSBLastPeekConfig,
	totalNormalBidders int,
) []BidderRequest {
	maxPrice, winnerBidder := peekChannelAvailableMaxBidPriceWithinTimeout("lastPeek", chBids, lastPeekConfig, totalNormalBidders)
	for reqIdx := range lastPeekBidderRequests {
		for idx := range lastPeekBidderRequests[reqIdx].BidRequest.Imp {
			bidder := &lastPeekBidderRequests[reqIdx]
			// for msp bidders, update req.imp.ext.bidder with peek_winner_price and peek_winner_bidder
			// which will be used later by msp module stage: https://github.com/ParticleMedia/msp/blob/master/pkg/modules/dam_buckets/module/hook_bidder_request.go#L69
			if _, found := mspBidders[bidder.BidderName]; found {
				extBytes, err := jsonObject(bidder.BidRequest.Imp[idx].Ext, "bidder")
				if err == nil {
					var impExt MSPFloor
					err = json.Unmarshal(extBytes, &impExt)
					if err == nil {
						impExt.PeekWinnerPrice = maxPrice
						impExt.PeekWinnerBidder = winnerBidder
						updatedBytes, _ := json.Marshal(impExt)
						updatedBidderBytes, _ := jsonpatch.MergePatch(extBytes, updatedBytes)
						updatedExtBytes, _ := jsonparser.Set(bidder.BidRequest.Imp[idx].Ext, updatedBidderBytes, "bidder")
						bidder.BidRequest.Imp[idx].Ext = updatedExtBytes
					}
				}
			}
		}
	}
	return lastPeekBidderRequests
}

func jsonObject(data []byte, keys ...string) ([]byte, error) {
	if result, dataType, _, err := jsonparser.Get(data, keys...); err == nil && dataType == jsonparser.Object {
		return result, nil
	} else {
		return nil, err
	}
}

func extractMSBInfoBidders(reqList []BidderRequest) MSBConfig {
	if len(reqList) > 0 {
		return ExtractMSBInfoReq(reqList[0].BidRequest)
	}

	return MSBConfig{}
}

func ExtractMSBInfoReq(req *openrtb2.BidRequest) MSBConfig {
	var config MSBExt
	if req != nil {
		err := json.Unmarshal(req.Ext, &config)
		if err != nil {
			glog.Error("MSB extract bidder config:", err)
		}
	}
	return config.MSB
}
