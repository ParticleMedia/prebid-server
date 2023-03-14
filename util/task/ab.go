package task

import (
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/ParticleMedia/ab-go-sdk"
	"github.com/prebid/prebid-server/util/timeutil"
)

const ab_url = "http://ab-api.ha.nb.com:8220"

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
	ab.Init(&ab.ABConfig{
		App:          "msp-prebid-server",
		Url:          ab_url,
		Layers:       []string{"ad_android", "ad_ios"},
		EnableCohort: false, // if you don't use cohort feature, please set to false
	})
	return &AB{
		httpClient:  httpClient,
		lastUpdated: atomic.Value{},
		time:        &timeutil.RealTime{},
	}
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (abUtl *AB) GetBucket(uid string, os string, cv string) (error, []string) {
	result := []string{}
	os = strings.ToLower(os)
	cv = strings.Replace(cv, ".", "", -1)
	if len(cv) < 6 {
		cv = cv + strings.Repeat("0", 6-len(cv))
	}
	uidInt64, _ := strconv.ParseUint(uid, 10, 64)
	response := ab.AB(
		ab.NewABContext(
			uid,
		).WithUserid(
			uint32(uidInt64),
		).WithConditions(map[string]interface{}{
			"platform": os,
			"cv":       cv,
		}),
	)

	for k, v := range response.Config {
		result = append(result, strings.Replace(k+"_"+v, ".", "_", -1))
	}

	// all bucket is only used for metrics now, it can also be set in the future
	result = append(result, "all")
	return nil, result
}
