package applovin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/prebid/openrtb/v17/openrtb2"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type ApplovinAdapter struct {
	endpoint string
	token    string
}

func (a *ApplovinAdapter) MakeRequests(request *openrtb2.BidRequest, _ *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")
	headers.Add("X-Account-Key", a.token)
	impressions := request.Imp
	result := make([]*adapters.RequestData, 0, len(impressions))
	errs := make([]error, 0, len(impressions))

	for _, impression := range impressions {
		if impression.Banner == nil && impression.Video == nil && impression.Native == nil {
			errs = append(errs, &errortypes.BadInput{
				Message: "Applovin only supports banner, video or native ads",
			})
			continue
		}
		if impression.Banner != nil {
			if impression.Banner.W == nil || impression.Banner.H == nil || *impression.Banner.W == 0 || *impression.Banner.H == 0 {
				if len(impression.Banner.Format) == 0 {
					errs = append(errs, &errortypes.BadInput{
						Message: "banner size information missing",
					})
					continue
				}
				banner := *impression.Banner
				banner.W = openrtb2.Int64Ptr(banner.Format[0].W)
				banner.H = openrtb2.Int64Ptr(banner.Format[0].H)
				impression.Banner = &banner
			}
		}
		if len(impression.Ext) == 0 {
			errs = append(errs, errors.New("impression extensions required"))
			continue
		}
		var bidderExt adapters.ExtImpBidder
		err := json.Unmarshal(impression.Ext, &bidderExt)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(bidderExt.Bidder) == 0 {
			errs = append(errs, errors.New("bidder required"))
			continue
		}
		var impressionExt openrtb_ext.ExtImpApplovin
		err = json.Unmarshal(bidderExt.Bidder, &impressionExt)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if impressionExt.TagID == "" {
			errs = append(errs, errors.New("Applovin token required"))
			continue
		}
		request.Imp = []openrtb2.Imp{impression}
		body, err := json.Marshal(request)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		result = append(result, &adapters.RequestData{
			Method:  "POST",
			Uri:     a.endpoint,
			Body:    body,
			Headers: headers,
		})
	}

	request.Imp = impressions

	if len(result) == 0 {
		return nil, errs
	}
	return result, errs
}

func (a *ApplovinAdapter) MakeBids(bidReq *openrtb2.BidRequest, adapterReq *adapters.RequestData, adapterResp *adapters.ResponseData) (*adapters.BidderResponse, []error) {
	if adapterResp.StatusCode != http.StatusOK {
		if adapterResp.StatusCode == http.StatusNoContent {
			return nil, nil
		}
		if adapterResp.StatusCode == http.StatusBadRequest {
			return nil, []error{&errortypes.BadInput{
				Message: fmt.Sprintf("Unexpected status code: %d", adapterResp.StatusCode),
			}}
		}
		return nil, []error{&errortypes.BadServerResponse{
			Message: fmt.Sprintf("Unexpected status code: %d", adapterResp.StatusCode),
		}}
	}

	var bidResp openrtb2.BidResponse
	if err := json.Unmarshal(adapterResp.Body, &bidResp); err != nil {
		return nil, []error{&errortypes.BadServerResponse{
			Message: fmt.Sprintf("Failed to unmarshal bid response: %s", err.Error()),
		}}
	}

	bidderResp := adapters.NewBidderResponseWithBidsCapacity(len(bidReq.Imp))
	var errors []error

	for _, seatbid := range bidResp.SeatBid {
		for _, bid := range seatbid.Bid {
			bid.ImpID = bidReq.Imp[0].ID
			bidderResp.Bids = append(bidderResp.Bids, &adapters.TypedBid{
				Bid:     &bid,
				BidType: openrtb_ext.BidTypeBanner,
			})
			break
		}
	}

	if bidResp.Cur != "" {
		bidderResp.Currency = bidResp.Cur
	}
	return bidderResp, errors
}

func Builder(bidderName openrtb_ext.BidderName, config config.Adapter, server config.Server) (adapters.Bidder, error) {
	if config.AppSecret == "" {
		return nil, errors.New("AppSecret is not configured. Did you set adapters.applovin.app_secret in the app config?")
	}
	bidder := &ApplovinAdapter{
		endpoint: config.Endpoint,
		token:    config.AppSecret,
	}
	return bidder, nil
}
