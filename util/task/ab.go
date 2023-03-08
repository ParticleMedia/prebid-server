package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/util/timeutil"
)

const ab_url = "http://ab-api.ha.nb.com:8220/ab/newsbreak/user"

type AbResponse struct {
	Abkv map[string]string `json:"exp"`
}

type AB struct {
	httpClient  httpClient
	lastUpdated atomic.Value
	time        timeutil.Time
}

func NewAB(
	httpClient httpClient,
) *AB {
	return &AB{
		httpClient:  httpClient,
		lastUpdated: atomic.Value{},
		time:        &timeutil.RealTime{},
	}
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (ab *AB) GetBucket(uid string, os string, cv string) (error, []string) {
	// all bucket is only used for metrics now, it can also be set in the future
	result := []string{"all"}
	request, err := http.NewRequest("GET", ab_url+"/"+uid, nil)
	if err != nil {
		return err, result
	}
	os = strings.ToLower(os)
	cv = strings.Replace(cv, ".", "", -1)
	if len(cv) < 6 {
		cv = cv + strings.Repeat("0", 6-len(cv))
	}
	q := request.URL.Query()
	q.Add("tag", fmt.Sprintf("{\"platform\":\"%s\", \"cv\":\"%s\"}", os, cv))
	request.URL.RawQuery = q.Encode()

	response, err := ab.httpClient.Do(request)
	if err != nil {
		return err, result
	}

	if response.StatusCode >= 400 {
		message := fmt.Sprintf("The ab api request failed with status code %d", response.StatusCode)
		return &errortypes.BadServerResponse{Message: message}, nil
	}

	defer response.Body.Close()

	bytesJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err, result
	}

	abResponse := &AbResponse{}

	err = json.Unmarshal(bytesJSON, abResponse)
	if err != nil {
		return err, result
	}

	for k, v := range abResponse.Abkv {
		result = append(result, k+"_"+v)
	}

	return nil, result
}
