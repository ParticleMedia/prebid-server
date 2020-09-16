package vrtcal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type VrtcalAdapter struct {
	endpoint string
}

func processImp(imp *openrtb.Imp) error {
	// get the Vrtcal extension
	var ext adapters.ExtImpBidder
	var vrtcalExt openrtb_ext.ExtImpVrtcal
	if err := json.Unmarshal(imp.Ext, &ext); err != nil {
		return err
	}
	if err := json.Unmarshal(ext.Bidder, &vrtcalExt); err != nil {
		return err
	}

	// set impression bid floor
	imp.BidFloor = vrtcalExt.BidFloor
	// no error
	return nil
}

func (a *VrtcalAdapter) MakeRequests(request *openrtb.BidRequest, reqInfo *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	var errs []error
	var validImps []openrtb.Imp
	var adapterRequests []*adapters.RequestData

	for _, imp := range request.Imp {
		if err := processImp(&imp); err == nil {
			validImps = append(validImps, imp)
		} else {
			errs = append(errs, err)
		}
	}

	if len(validImps) == 0 {
		return nil, errs
	}

	request.Imp = validImps

	reqJSON, err := json.Marshal(request)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")

	reqData := adapters.RequestData{
		Method:  "POST",
		Uri:     a.endpoint,
		Body:    reqJSON,
		Headers: headers,
	}

	adapterRequests = append(adapterRequests, &reqData)

	return adapterRequests, errs
}

// MakeBids make the bids for the bid response.
func (a *VrtcalAdapter) MakeBids(internalRequest *openrtb.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {

	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, []error{&errortypes.BadInput{
			Message: fmt.Sprintf("Unexpected status code: %d. Run with request.debug = 1 for more info", response.StatusCode),
		}}
	}

	if response.StatusCode != http.StatusOK {
		return nil, []error{&errortypes.BadServerResponse{
			Message: fmt.Sprintf("Unexpected status code: %d. Run with request.debug = 1 for more info", response.StatusCode),
		}}
	}

	var bidResp openrtb.BidResponse

	if err := json.Unmarshal(response.Body, &bidResp); err != nil {
		return nil, []error{err}
	}

	bidResponse := adapters.NewBidderResponseWithBidsCapacity(1)

	for _, sb := range bidResp.SeatBid {
		for i := range sb.Bid {
			bidResponse.Bids = append(bidResponse.Bids, &adapters.TypedBid{
				Bid:     &sb.Bid[i],
				BidType: "banner",
			})
		}
	}
	return bidResponse, nil

}

func NewVrtcalBidder(endpoint string) *VrtcalAdapter {
	return &VrtcalAdapter{
		endpoint: endpoint,
	}
}
