package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/util/secretutil"
	"github.com/prebid/prebid-server/util/timeutil"
)

type DealInfo struct {
	DealId          string
	DealName        string
	CustomTargeting map[string][]string
}

const (
	pubmaticTargetingApiUrl = "https://api.pubmatic.com/v1/inventory/targeting"
	pubmaticAccessKey       = "PUBMATIC_ACCESS_TOKEN"
)

type PubmaticResponse struct {
	Items []PubmaticDealItem `json:"items"`
}

type PubmaticStatus struct {
	Name string `json:"name"`
}

type PubmaticTargeting struct {
	Id                int                    `json:"id"`
	KeyValueTargeting [][]PubmaticKeyValPair `json:"keyValueTargeting"`
}

type PubmaticKeyValPair struct {
	Key   string          `json:"k"`
	Value []PubmaticValue `json:"v"`
}

type PubmaticValue struct {
	Val string `json:"eV"`
}

type PubmaticDealItem struct {
	DealId        int               `json:"id"`
	DealDisplayId string            `json:"dealId"`
	DealName      string            `json:"name"`
	Status        PubmaticStatus    `json:"status"`
	Targeting     PubmaticTargeting `json:"targeting"`
}

// DealFetcher holds the deal dictionary
type DealFetcher struct {
	httpClient    httpClient
	lastUpdated   atomic.Value
	time          timeutil.Time
	bidderInfo    config.BidderInfos
	pubmaticToken string
	dealLock      sync.RWMutex
	// bidderName -> deal_id -> dealInfo
	dealMap map[string]map[string]DealInfo
}

func NewDealFetcher(
	httpClient httpClient,
	bidderInfo config.BidderInfos,
) *DealFetcher {
	return &DealFetcher{
		httpClient:    httpClient,
		lastUpdated:   atomic.Value{},
		time:          &timeutil.RealTime{},
		bidderInfo:    bidderInfo,
		pubmaticToken: secretutil.GetSecret(pubmaticAccessKey),
		dealMap:       map[string]map[string]DealInfo{},
		dealLock:      sync.RWMutex{},
	}
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (df *DealFetcher) Run() error {
	for bidderName, info := range df.bidderInfo {
		if info.IsEnabled() && info.DealEndpoint != "" {
			err := df.fetch(bidderName, info)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func (df *DealFetcher) fetch(bidderName string, info config.BidderInfo) error {
	switch bidderName {
	case "pubmatic":
		return df.fetchPubmatic(bidderName, info)
	default:
		fmt.Println("Unsupported deal bidder:", bidderName)
	}

	return nil
}

func (df *DealFetcher) generatePubmaticGetRequestWithHeaders(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", df.pubmaticToken))

	// TODO: move these account information to account level in the future
	q := request.URL.Query()
	q.Add("view", "summary")
	q.Add("filters", "loggedInOwnerId eq 162568")
	q.Add("filters", "loggedInOwnerTypeId eq 1")
	request.URL.RawQuery = q.Encode()
	return request, nil
}

func (df *DealFetcher) GetDealInfo(bidderName string, dealId string) *DealInfo {
	df.dealLock.RLock()
	defer df.dealLock.RUnlock()
	dealInfos, ok := df.dealMap[bidderName]
	if ok {
		dealInfo, ok := dealInfos[dealId]
		if ok {
			return &dealInfo
		}
	}
	return nil
}

func (df *DealFetcher) fetchPubmatic(bidderName string, info config.BidderInfo) error {
	request, err := df.generatePubmaticGetRequestWithHeaders(info.DealEndpoint)
	if err != nil {
		return err
	}

	response, err := df.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		message := fmt.Sprintf("The deal api request to %s failed with status code %d", bidderName, response.StatusCode)
		return &errortypes.BadServerResponse{Message: message}
	}

	defer response.Body.Close()

	bytesJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	dealResponse := &PubmaticResponse{}

	err = json.Unmarshal(bytesJSON, dealResponse)
	if err != nil {
		return err
	}

	dealInfos := map[string]DealInfo{}
	for _, deal := range dealResponse.Items {
		if deal.Status.Name == "Active" {
			targetingRequest, err := df.generatePubmaticGetRequestWithHeaders(fmt.Sprintf("%s/%d", pubmaticTargetingApiUrl, deal.Targeting.Id))
			if err == nil {
				response, err = df.httpClient.Do(targetingRequest)
				if err == nil && response.StatusCode == 200 {
					bytesJSON, err := ioutil.ReadAll(response.Body)
					if err == nil {
						targeting := &PubmaticTargeting{}
						err = json.Unmarshal(bytesJSON, targeting)
						if err == nil {
							dealInfo := DealInfo{
								DealId:          deal.DealDisplayId,
								DealName:        deal.DealName,
								CustomTargeting: map[string][]string{},
							}
							for _, kv := range targeting.KeyValueTargeting {
								for _, kvPair := range kv {
									for _, val := range kvPair.Value {
										dealInfo.CustomTargeting[kvPair.Key] = append(dealInfo.CustomTargeting[kvPair.Key], val.Val)
									}
								}
							}
							if len(dealInfo.CustomTargeting) > 0 {
								dealInfos[deal.DealDisplayId] = dealInfo
							}
						}
					}

				}
			}
		}
	}

	df.dealLock.Lock()
	df.dealMap[bidderName] = dealInfos
	df.dealLock.Unlock()
	fmt.Printf("Updated %s deals with custom targeting, count: %d\n", bidderName, len(dealInfos))
	return nil
}
