package google

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prebid/openrtb/v19/openrtb2"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
)

const (
	QUERY_INFO_KEY = "query_info"
)

type GoogleAdapter struct {
	endpoint string
}

func getQueryInfo(imp *openrtb2.Imp) string {
	var bidderExt openrtb_ext.ExtImpBidderGoogle
	if err := json.Unmarshal(imp.Ext, &bidderExt); err != nil {
		if queryInfo, ok := bidderExt.Context.DataMap[QUERY_INFO_KEY]; ok && len(queryInfo) > 0 {
			return queryInfo[0]
		}
	}
	return ""
}

func (a *GoogleAdapter) MakeRequests(request *openrtb2.BidRequest, _ *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	headers := http.Header{}
	headers.Add("Accept-Encoding", "gzip,deflate")
	impressions := request.Imp

	result := make([]*adapters.RequestData, 0, len(impressions))
	errs := make([]error, 0, len(impressions))

	if request.Device != nil && request.Device.UA != "" {
		headers.Add("User-Agent", request.Device.UA)
	} else {
		return result, errs
	}

	if request.Device.IP != "" {
		headers.Add("X-Forwarded-For", request.Device.IP)
	}

	for k, v := range headers {
		for _, val := range v {
			fmt.Printf("%s:%s\n", k, val)
		}
	}

	for _, impression := range impressions {
		reqUrl, err := url.Parse(a.endpoint)
		query := reqUrl.Query()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		queryInfo := getQueryInfo(&impression)
		if queryInfo == "" {
			errs = append(errs, &errortypes.BadInput{
				Message: "Missing query info for Google adapter",
			})
			continue
		}

		query.Set("gsig", queryInfo)
		if impression.Banner == nil && impression.Video == nil && impression.Native == nil {
			errs = append(errs, &errortypes.BadInput{
				Message: "Google only supports banner, video or native ads",
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
			query.Set("sz", fmt.Sprintf("%dx%d", *impression.Banner.W, *impression.Banner.H))
		}
		if len(impression.Ext) == 0 {
			errs = append(errs, errors.New("impression extensions required"))
			continue
		}
		query.Set("pubf", fmt.Sprintf("%d", int(impression.BidFloor*1e6)))
		var bidderExt adapters.ExtImpBidder
		err = json.Unmarshal(impression.Ext, &bidderExt)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(bidderExt.Bidder) == 0 {
			errs = append(errs, errors.New("bidder required"))
			continue
		}
		var impressionExt openrtb_ext.ExtImpGoogle
		err = json.Unmarshal(bidderExt.Bidder, &impressionExt)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if impressionExt.IU == "" {
			errs = append(errs, errors.New("google ad unit required"))
			continue
		}

		query.Set("iu", impressionExt.IU)
		query.Set("c", fmt.Sprintf("%d", time.Now().UnixMilli()))

		result = append(result, &adapters.RequestData{
			Method:  "GET",
			Uri:     fmt.Sprintf("%s?%s", reqUrl.String(), getParamString(query)),
			Headers: headers,
		})
		fmt.Println("url:", fmt.Sprintf("%s?%s", reqUrl.String(), getParamString(query)))
	}

	request.Imp = impressions

	if len(result) == 0 {
		return nil, errs
	}
	return result, errs
}

func getParamString(values url.Values) string {
	var buf strings.Builder
	for key, valueList := range values {
		for _, value := range valueList {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(key)
			buf.WriteByte('=')
			buf.WriteString(value)
		}
	}
	return buf.String()
}

func (a *GoogleAdapter) MakeBids(bidReq *openrtb2.BidRequest, adapterReq *adapters.RequestData, adapterResp *adapters.ResponseData) (*adapters.BidderResponse, []error) {
	fmt.Println("status code:", adapterResp.StatusCode)
	fmt.Println("body:", string(adapterResp.Body))
	if adapterResp.StatusCode != http.StatusOK {
		return nil, nil
	}
	if string(adapterResp.Body) == "" {
		return nil, nil
	}

	reader, err := gzip.NewReader(bytes.NewReader(adapterResp.Body))
	if err != nil {
		return nil, nil
	}
	defer reader.Close()

	if err != nil {
		return nil, nil
	}

	body, err := io.ReadAll(reader)

	if err != nil {
		return nil, nil
	}

	reqUrl, err := url.Parse(adapterReq.Uri)
	if err != nil {
		return nil, nil
	}
	size := reqUrl.Query().Get("sz")
	var width int64
	var height int64
	_, err = fmt.Sscanf(size, "%dx%d", &width, &height)
	if err != nil {
		return nil, nil
	}
	adUnit := reqUrl.Query().Get("iu")
	bidderResp := adapters.NewBidderResponseWithBidsCapacity(len(bidReq.Imp))
	var errors []error

	ext := openrtb_ext.GoogleBidExt{
		Detail: openrtb_ext.GoogleBidExtDetail{
			AdUnit: adUnit,
		},
	}
	extBytes, err := json.Marshal(&ext)
	if err != nil {
		fmt.Println("error when marshal ext:", err.Error())
		return nil, nil
	}
	bid := openrtb2.Bid{
		ID:      fmt.Sprintf("%d", time.Now().UnixMilli()),
		ImpID:   bidReq.ID,
		Price:   bidReq.Imp[0].BidFloor,
		AdM:     string(body),
		ADomain: []string{"unknown"},
		CrID:    "unknown",
		W:       width,
		H:       height,
		Ext:     extBytes,
	}
	bidderResp.Bids = append(bidderResp.Bids, &adapters.TypedBid{
		Bid:     &bid,
		BidType: openrtb_ext.BidTypeBanner,
	})
	return bidderResp, errors
}

func Builder(bidderName openrtb_ext.BidderName, config config.Adapter, server config.Server) (adapters.Bidder, error) {
	bidder := &GoogleAdapter{
		endpoint: config.Endpoint,
	}
	return bidder, nil
}
