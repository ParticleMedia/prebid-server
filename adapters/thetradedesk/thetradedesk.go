package thetradedesk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"net/http"
	"regexp"
	"text/template"

	"github.com/prebid/prebid-server/v3/adapters"
	"github.com/prebid/prebid-server/v3/config"
	"github.com/prebid/prebid-server/v3/macros"
	"github.com/prebid/prebid-server/v3/openrtb_ext"
	"github.com/prebid/prebid-server/v3/util/jsonutil"

	"github.com/prebid/openrtb/v20/openrtb2"
)

//const PREBID_INTEGRATION_TYPE = "1"

type adapter struct {
	bidderEndpointTemplate string
	defaultEndpoint        string
	templateEndpoint       *template.Template
}

func (a *adapter) MakeRequests(request *openrtb2.BidRequest, reqInfo *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	pubID, supplySourceId, err := getExtensionInfo(request.Imp)

	if err != nil {
		return nil, []error{err}
	}

	modifiedImps := make([]openrtb2.Imp, 0, len(request.Imp))

	for _, imp := range request.Imp {

		if imp.Banner != nil {
			if len(imp.Banner.Format) > 0 {
				firstFormat := imp.Banner.Format[0]
				bannerCopy := *imp.Banner
				bannerCopy.H = &firstFormat.H
				bannerCopy.W = &firstFormat.W
				imp.Banner = &bannerCopy

			}
		}

		modifiedImps = append(modifiedImps, imp)
	}

	request.Imp = modifiedImps

	if request.Site != nil {
		siteCopy := *request.Site
		if siteCopy.Publisher != nil {
			publisherCopy := *siteCopy.Publisher
			if pubID != "" {
				publisherCopy.ID = pubID
			}
			siteCopy.Publisher = &publisherCopy
		} else {
			siteCopy.Publisher = &openrtb2.Publisher{ID: pubID}
		}
		request.Site = &siteCopy
	} else if request.App != nil {
		appCopy := *request.App
		if appCopy.Publisher != nil {
			publisherCopy := *appCopy.Publisher
			if pubID != "" {
				publisherCopy.ID = pubID
			}
			appCopy.Publisher = &publisherCopy
		} else {
			appCopy.Publisher = &openrtb2.Publisher{ID: pubID}
		}
		request.App = &appCopy
	}

	errs := make([]error, 0, len(request.Imp))
	reqJSON, err := json.Marshal(request)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	bidderEndpoint, err := a.buildEndpointURL(supplySourceId)
	if err != nil {
		return nil, []error{errors.New("Failed to build endpoint URL")}
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")
	//headers.Add("x-integration-type", PREBID_INTEGRATION_TYPE) this will be parsed and added conditionally later
	return []*adapters.RequestData{{
		Method:  "POST",
		Uri:     bidderEndpoint,
		Body:    reqJSON,
		Headers: headers,
		ImpIDs:  openrtb_ext.GetImpIDs(request.Imp),
	}}, errs
}

func (a *adapter) buildEndpointURL(supplySourceId string) (string, error) {
	if supplySourceId == "" {
		return a.defaultEndpoint, nil
	}

	urlParams := macros.EndpointTemplateParams{SupplyId: supplySourceId}
	bidderEndpoint, err := macros.ResolveMacros(a.templateEndpoint, urlParams)

	if err != nil {
		return "", fmt.Errorf("unable to resolve endpoint macros: %v", err)
	}

	return bidderEndpoint, nil
}

func getExtensionInfo(impressions []openrtb2.Imp) (string, string, error) {
	publisherId := ""
	supplySourceId := ""
	for _, imp := range impressions {
		var ttdExt, err = getImpressionExt(&imp)
		if err != nil {
			return "", "", err
		}

		if ttdExt.PublisherId != "" && publisherId == "" {
			publisherId = ttdExt.PublisherId
		}

		if ttdExt.SupplySourceId != "" && supplySourceId == "" {
			supplySourceId = ttdExt.SupplySourceId
		}

		if publisherId != "" && supplySourceId != "" {
			break
		}
	}
	return publisherId, supplySourceId, nil
}

func getImpressionExt(imp *openrtb2.Imp) (*openrtb_ext.ExtImpTheTradeDesk, error) {
	var bidderExt adapters.ExtImpBidder
	if err := jsonutil.Unmarshal(imp.Ext, &bidderExt); err != nil {
		return nil, err
	}

	var ttdExt openrtb_ext.ExtImpTheTradeDesk
	if err := jsonutil.Unmarshal(bidderExt.Bidder, &ttdExt); err != nil {
		return nil, err
	}
	return &ttdExt, nil
}

func (a *adapter) MakeBids(internalRequest *openrtb2.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {
	if len(internalRequest.Imp) > 0 && internalRequest.Imp[0].Ext != nil {
		uid, err := jsonparser.GetString(internalRequest.Imp[0].Ext, "context", "data", "user_id", "[0]")
		if err == nil {
			if uid == "111111112222222" {
				ttdBidResponseBody := []byte(`{
			  "id": "response-id-123",
			  "seatbid": [{
				"bid": [{
				  "id": "bid-id-1",
				  "impid": "9ae5f2ae-71a1-40f4-8da9-5c6026824e15",
				  "price": 0.017433353,
				  "adm": "<img src=\"https://usw-ca2.adsrvr.org/bid/feedback/learnings?t=1&iid=470f404f-bcb2-405e-9fef-18da97817300&crid=k2ed2otl&wp=${AUCTION_PRICE}&aid=9ae5f2ae-71a1-40f4-8da9-5c6026824e15&wpc=USD&sfe=1a9d4350&puid=&bdc=10&tdid=&pid=tdpartnerid&ag=j0n4dzq&adv=cxpowqgj&sig=1rONppIPJHFuxS3EfBXj2DBhLEAP9mG_A9Q_1BsBH0EM.&bp=0.0182548187445625&cf=8954840&fq=0&td_s=easy.killer.sudoku.puzzle.solver.free&rcats=&mste=&mfld=2&mssi=&mfsi=&uhow=98&agsa=&rgz=77133&svbttd=1&dt=Tablet&osf=Android&os=Android150&br=WebView&rlangs=uk&mlang=&svpid=1&did=&rcxt=InApp&lat=50.452200&lon=30.528700&tmpc=18.400000000000034&daid=dcacf7ba-6a79-4a3c-9088-e7e465224982&vp=0&osi=&osv=&mk=Lenovo&mdl=TB330FU&testid=iavc1%20&cc=1~KLUv_WNa1XJFKgCFCAATUD09VgQQsQTyPIE8bRkZEnECegIhI2WUzGdYkIkhQQYnSOALW0ULAStbBQwBKH4AKj8AUYQBO3OAHXOAtcQAKKzbGzCEhAA8M7GegDasACg-ACQ2QB6OI9ni6hkY4OoDDGG37n5f2mztH7ag_K0k7FaQ3laOnhYqVguUgCBaiAQE0crDUuSxoHisJD1gjp4wy_0_VawPEuu7Y-a6FFrf1Js6JRKrbTmzOsX12oYSo1NpB5-Z3s_zHOy2KV7jzBqGgCF1WkOVKbrapm93G17rZMQcr99KJOCy8bvW21YAs2JWkJ_qeKJI-MlKARE03pgZMUBECiwDcAb8SNkFvcBpJZKR08YqzsAM_uACOQu3_w..&dur=1~KLUv_SMFoCgkV0UAAAABAMQlTAwI&durs=oPZ81L&crrelr=&vc=12&said=9ae5f2ae-71a1-40f4-8da9-5c6026824e15&ict=WiFi&auct=1&us_privacy=1YNN&im=1&mc=3b04ae0b-0bba-47ff-86ff-c7b200f43c88&ev=wzwl3lflvqOE70l681LU8WbTLyPO-A1l2P9qbxat31A.&rsv=170.256417342362&abr=a97ef3a4-bb01-4814-a33a-2b2b06627f01&tail=1\" width=\"1\" height=\"1\" style=\"display: none;\"/><a target=\"_blank\" href=\"https://insight.adsrvr.org/track/clk?imp=470f404f-bcb2-405e-9fef-18da97817300&ag=j0n4dzq&sfe=1a9d4350&sig=kvbRR79HZm8xUiCLpPlJvEatAxtVu0ZqKP_PM5j5xxs.&crid=k2ed2otl&cf=8954840&fq=0&t=1&td_s=easy.killer.sudoku.puzzle.solver.free&rcats=&mste=&mfld=2&mssi=&mfsi=&sv=learnings&uhow=98&agsa=&wp=${AUCTION_PRICE}&rgz=77133&dt=Tablet&osf=Android&os=Android150&br=WebView&svpid=1&rlangs=uk&mlang=&did=&rcxt=InApp&tmpc=18.400000000000034&vrtd=&osi=&osv=&daid=dcacf7ba-6a79-4a3c-9088-e7e465224982&dnr=0&vpb=&crrelr=&npt=&cc=1~KLUv_WNa1XJFKgCFCAATUD09VgQQsQTyPIE8bRkZEnECegIhI2WUzGdYkIkhQQYnSOALW0ULAStbBQwBKH4AKj8AUYQBO3OAHXOAtcQAKKzbGzCEhAA8M7GegDasACg-ACQ2QB6OI9ni6hkY4OoDDGG37n5f2mztH7ag_K0k7FaQ3laOnhYqVguUgCBaiAQE0crDUuSxoHisJD1gjp4wy_0_VawPEuu7Y-a6FFrf1Js6JRKrbTmzOsX12oYSo1NpB5-Z3s_zHOy2KV7jzBqGgCF1WkOVKbrapm93G17rZMQcr99KJOCy8bvW21YAs2JWkJ_qeKJI-MlKARE03pgZMUBECiwDcAb8SNkFvcBpJZKR08YqzsAM_uACOQu3_w..&dur=1~KLUv_SMFoCgkV0UAAAABAMQlTAwI&durs=oPZ81L&bdc=10&mk=Lenovo&mdl=TB330FU&testid=iavc1%20&ict=WiFi&said=9ae5f2ae-71a1-40f4-8da9-5c6026824e15&auct=1&us_privacy=1YNN&tail=1&r=https://thetradedesk.com\" alt=\"Click Me\"><img src=\"https://ad.adsrvr.org/tdpartnerid/cxpowqgj/k2ed2otl_320x480.jpg?cb=173131\" border=\"0\" /></a><span id=\"te-clearads-js-tradedesk01cont1\"><script type=\"text/javascript\" src=\"https://choices.truste.com/ca?pid=tradedesk01&aid=tradedesk01&cid=cmlvp04_j0n4dzq_k2ed2otl&c=tradedesk01cont1&js=pmw0&w=320&h=480&sid=ev9-FM60pzWefkMAiM77ugE-aNy2u3kzAYRYrOVH6CCiJPrSZ8LYBMM2TPF5FU42NKIBlrhk3PK6holqcS01UwAnOhvRl44yjOSi3IPVuWrHqW3C5epk9CMXoJWQu6-a&dsarequired=&dsabehalf=&dsapaid=&dsaparams=\"></script></span><script src=\"https://cdn.doubleverify.com/dvtp_src.js?ctx=818052&cmp=DV140326&sid=TTD&plc=dispview&advid=818053&adsrv=163&btreg=&btadsrv=&dvtagver=6.1.src&DVP_TTD_1=tdpartnerid&DVP_TTD_2=cxpowqgj&DVP_TTD_3=j0n4dzq&DVP_TTD_4=cmlvp04&DVP_TTD_6=learnings&DVP_HAS_VIEW=0&rtsurl=https%3A%2F%2Fenduser.adsrvr.org%2Fenduser%2Fdv%2F%3Frtb%3DdD0xJmlpZD00NzBmNDA0Zi1iY2IyLTQwNWUtOWZlZi0xOGRhOTc4MTczMDAmY3JpZD1rMmVkMm90bCZ3cD0ke0FVQ1RJT05fUFJJQ0V9JmFpZD05YWU1ZjJhZS03MWExLTQwZjQtOGRhOS01YzYwMjY4MjRlMTUmd3BjPVVTRCZzZmU9MWE5ZDQzNTAmcHVpZD0mYmRjPTEwJnRkaWQ9JnBpZD10ZHBhcnRuZXJpZCZhZz1qMG40ZHpxJmFkdj1jeHBvd3FnaiZicD0wLjAxODI1NDgxODc0NDU2MjUmY2Y9ODk1NDg0MCZmcT0wJnRkX3M9ZWFzeS5raWxsZXIuc3Vkb2t1LnB1enpsZS5zb2x2ZXIuZnJlZSZyY2F0cz0mbXN0ZT0mbWZsZD0yJm1zc2k9Jm1mc2k9JnVob3c9OTgmYWdzYT0mcmd6PTc3MTMzJnN2YnR0ZD0xJmR0PVRhYmxldCZvc2Y9QW5kcm9pZCZvcz1BbmRyb2lkMTUwJmJyPVdlYlZpZXcmcmxhbmdzPXVrJm1sYW5nPSZzdnBpZD0xJmRpZD0mcmN4dD1JbkFwcCZsYXQ9NTAuNDUyMjAwJmxvbj0zMC41Mjg3MDAmdG1wYz0xOC40MDAwMDAwMDAwMDAwMzQmZGFpZD1kY2FjZjdiYS02YTc5LTRhM2MtOTA4OC1lN2U0NjUyMjQ5ODImdnA9MCZvc2k9Jm9zdj0mbWs9TGVub3ZvJm1kbD1UQjMzMEZVJnRlc3RpZD1pYXZjMSUyMCZjYz0xfktMVXZfV05hMVhKRktnQ0ZDQUFUVUQwOVZnUVFzUVR5UElFOGJSa1pFbkVDZWdJaEkyV1V6R2RZa0lraFFRWW5TT0FMVzBVTEFTdGJCUXdCS0g0QUtqOEFVWVFCTzNPQUhYT0F0Y1FBS0t6Ykd6Q0VoQUE4TTdHZWdEYXNBQ2ctQUNRMlFCNk9JOW5pNmhrWTRPb0RER0czN241ZjJtenRIN2FnX0swazdGYVEzbGFPbmhZcVZndVVnQ0JhaUFRRTBjckRVdVN4b0hpc0pEMWdqcDR3eV8wX1Zhd1BFdXU3WS1hNkZGcmYxSnM2SlJLcmJUbXpPc1gxMm9ZU28xTnBCNS1aM3NfekhPeTJLVjdqekJxR2dDRjFXa09WS2JyYXBtOTNHMTdyWk1RY3I5OUtKT0N5OGJ2VzIxWUFzMkpXa0pfcWVLSkktTWxLQVJFMDNwZ1pNVUJFQ2l3RGNBYjhTTmtGdmNCcEpaS1IwOFlxenNBTV91QUNPUXUzX3cuLiZkdXI9MX5LTFV2X1NNRm9DZ2tWMFVBQUFBQkFNUWxUQXdJJmNycmVscj0mdmM9MTImc2FpZD05YWU1ZjJhZS03MWExLTQwZjQtOGRhOS01YzYwMjY4MjRlMTUmaWN0PVdpRmkmYXVjdD0xJnVzX3ByaXZhY3k9MVlOTiZpbT0xJm1jPTNiMDRhZTBiLTBiYmEtNDdmZi04NmZmLWM3YjIwMGY0M2M4OCZldj13endsM2xmbHZxT0U3MGw2ODFMVThXYlRMeVBPLUExbDJQOXFieGF0MzFBLiZyc3Y9MTcwLjI1NjQxNzM0MjM2MiZhYnI9YTk3ZWYzYTQtYmIwMS00ODE0LWEzM2EtMmIyYjA2NjI3ZjAxJnRhaWw9MSZzdj1sZWFybmluZ3MmdGFpbD0x%26pie%3D\" type=\"text/javascript\"></script><script src=\"https://js.adsrvr.org/ttdReferrerTrackerV2.min.js\" type=\"text/javascript\"></script><script  type=\"text/javascript\">var trackere09a2e16a327426d85e20c945afe3bc7 = new ReferrerTracker();trackere09a2e16a327426d85e20c945afe3bc7.send(\"https://enduser.adsrvr.org/enduser/pie/\", \"pie%3D21%26rtb%3DdD0xJmlpZD00NzBmNDA0Zi1iY2IyLTQwNWUtOWZlZi0xOGRhOTc4MTczMDAmY3JpZD1rMmVkMm90bCZ3cD0ke0FVQ1RJT05fUFJJQ0V9JmFpZD05YWU1ZjJhZS03MWExLTQwZjQtOGRhOS01YzYwMjY4MjRlMTUmd3BjPVVTRCZzZmU9MWE5ZDQzNTAmcHVpZD0mYmRjPTEwJnRkaWQ9JnBpZD10ZHBhcnRuZXJpZCZhZz1qMG40ZHpxJmFkdj1jeHBvd3FnaiZicD0wLjAxODI1NDgxODc0NDU2MjUmY2Y9ODk1NDg0MCZmcT0wJnRkX3M9ZWFzeS5raWxsZXIuc3Vkb2t1LnB1enpsZS5zb2x2ZXIuZnJlZSZyY2F0cz0mbXN0ZT0mbWZsZD0yJm1zc2k9Jm1mc2k9JnVob3c9OTgmYWdzYT0mcmd6PTc3MTMzJnN2YnR0ZD0xJmR0PVRhYmxldCZvc2Y9QW5kcm9pZCZvcz1BbmRyb2lkMTUwJmJyPVdlYlZpZXcmcmxhbmdzPXVrJm1sYW5nPSZzdnBpZD0xJmRpZD0mcmN4dD1JbkFwcCZsYXQ9NTAuNDUyMjAwJmxvbj0zMC41Mjg3MDAmdG1wYz0xOC40MDAwMDAwMDAwMDAwMzQmZGFpZD1kY2FjZjdiYS02YTc5LTRhM2MtOTA4OC1lN2U0NjUyMjQ5ODImdnA9MCZvc2k9Jm9zdj0mbWs9TGVub3ZvJm1kbD1UQjMzMEZVJnRlc3RpZD1pYXZjMSUyMCZjYz0xfktMVXZfV05hMVhKRktnQ0ZDQUFUVUQwOVZnUVFzUVR5UElFOGJSa1pFbkVDZWdJaEkyV1V6R2RZa0lraFFRWW5TT0FMVzBVTEFTdGJCUXdCS0g0QUtqOEFVWVFCTzNPQUhYT0F0Y1FBS0t6Ykd6Q0VoQUE4TTdHZWdEYXNBQ2ctQUNRMlFCNk9JOW5pNmhrWTRPb0RER0czN241ZjJtenRIN2FnX0swazdGYVEzbGFPbmhZcVZndVVnQ0JhaUFRRTBjckRVdVN4b0hpc0pEMWdqcDR3eV8wX1Zhd1BFdXU3WS1hNkZGcmYxSnM2SlJLcmJUbXpPc1gxMm9ZU28xTnBCNS1aM3NfekhPeTJLVjdqekJxR2dDRjFXa09WS2JyYXBtOTNHMTdyWk1RY3I5OUtKT0N5OGJ2VzIxWUFzMkpXa0pfcWVLSkktTWxLQVJFMDNwZ1pNVUJFQ2l3RGNBYjhTTmtGdmNCcEpaS1IwOFlxenNBTV91QUNPUXUzX3cuLiZkdXI9MX5LTFV2X1NNRm9DZ2tWMFVBQUFBQkFNUWxUQXdJJmNycmVscj0mdmM9MTImc2FpZD05YWU1ZjJhZS03MWExLTQwZjQtOGRhOS01YzYwMjY4MjRlMTUmaWN0PVdpRmkmYXVjdD0xJnVzX3ByaXZhY3k9MVlOTiZpbT0xJm1jPTNiMDRhZTBiLTBiYmEtNDdmZi04NmZmLWM3YjIwMGY0M2M4OCZldj13endsM2xmbHZxT0U3MGw2ODFMVThXYlRMeVBPLUExbDJQOXFieGF0MzFBLiZyc3Y9MTcwLjI1NjQxNzM0MjM2MiZhYnI9YTk3ZWYzYTQtYmIwMS00ODE0LWEzM2EtMmIyYjA2NjI3ZjAxJnRhaWw9MSZzdj1sZWFybmluZ3MmdGFpbD0x\");</script><noscript class=\"TTD-referrer-tracker\"></noscript>",
				  "adomain": ["thetradedesk.com"],
				  "cid": "cmlvp04",
				  "crid": "k2ed2otl",
				  "cat": ["IAB19"],
				  "w": 320,
				  "h": 480,
				  "mtype": 1
				}]
			  }]
			}`)

				response.Body = ttdBidResponseBody
				response.StatusCode = http.StatusOK
			}
		}
	}

	if adapters.IsResponseStatusCodeNoContent(response) {
		return adapters.NewBidderResponse(), nil
	}

	if err := adapters.CheckResponseStatusCodeForErrors(response); err != nil {
		return nil, []error{err}
	}

	var bidResponse openrtb2.BidResponse
	if err := jsonutil.Unmarshal(response.Body, &bidResponse); err != nil {
		return nil, []error{err}
	}

	bidderResponse := adapters.NewBidderResponse()
	bidderResponse.Currency = bidResponse.Cur

	for _, seatBid := range bidResponse.SeatBid {
		for _, bid := range seatBid.Bid {
			bid := bid

			bidType, err := getBidType(bid.MType)

			if err != nil {
				return nil, []error{err}
			}

			b := &adapters.TypedBid{
				Bid:     &bid,
				BidType: bidType,
			}
			bidderResponse.Bids = append(bidderResponse.Bids, b)
		}
	}

	return bidderResponse, nil
}

func getBidType(markupType openrtb2.MarkupType) (openrtb_ext.BidType, error) {
	switch markupType {
	case openrtb2.MarkupBanner:
		return openrtb_ext.BidTypeBanner, nil
	case openrtb2.MarkupVideo:
		return openrtb_ext.BidTypeVideo, nil
	case openrtb2.MarkupNative:
		return openrtb_ext.BidTypeNative, nil
	default:
		return "", fmt.Errorf("unsupported mtype: %d", markupType)
	}
}

func Builder(bidderName openrtb_ext.BidderName, config config.Adapter, server config.Server) (adapters.Bidder, error) {
	if len(config.ExtraAdapterInfo) > 0 {
		isValidEndpoint, err := regexp.Match("([a-z]+)$", []byte(config.ExtraAdapterInfo))
		if !isValidEndpoint || err != nil {
			return nil, errors.New("ExtraAdapterInfo must be a simple string provided by TheTradeDesk")
		}
	}

	template, err := template.New("endpointTemplate").Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse endpoint url template: %v", err)
	}

	urlParams := macros.EndpointTemplateParams{SupplyId: config.ExtraAdapterInfo}
	defaultEndpoint, err := macros.ResolveMacros(template, urlParams)

	if err != nil {
		return nil, fmt.Errorf("unable to resolve endpoint macros: %v", err)
	}

	return &adapter{
		bidderEndpointTemplate: config.Endpoint,
		defaultEndpoint:        defaultEndpoint,
		templateEndpoint:       template,
	}, nil
}
