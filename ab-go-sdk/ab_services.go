package ab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
)

type ABServices struct {
	Cfg           *ABConfig
	status        string
	cron          *cron.Cron
	services      []ABService
	cohortService *CohortService
}

func NewABServices(cfg *ABConfig) *ABServices {
	s := &ABServices{
		Cfg:    cfg,
		status: "success",
	}
	if cfg.EnableCohort {
		s.cohortService = NewCohortService(cfg.Url)
	}
	if err := s.reload(); err != nil {
		panic(err)
	}
	s.cron = cron.New()
	// nolint
	s.cron.AddFunc("* * * * *", s.Reload)
	s.cron.Start()

	return s
}

func (s *ABServices) Reload() {
	start := time.Now()
	if err := s.reload(); err != nil {
		glog.Errorf("Reload AB config err: %s", err)
		s.status = "failed"
		return
	}
	s.status = "success"
	glog.Info("Reload AB config success, cost: ", time.Since(start))
}

type Layers []*Layer

func (s Layers) Len() int {
	return len(s)
}
func (s Layers) Less(i, j int) bool {
	return s[i].Zone < s[j].Zone
}
func (s Layers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *ABServices) reload() error {
	resp, err := http.Get(s.Cfg.Url + "/ab/newsbreak/layers/" + strings.Join(s.Cfg.Layers, ",") + "?" + url.Values{
		"app":            {s.Cfg.App},
		"last_status":    {s.status},
		"client_version": {"go-1.0.0"},
		"cohort":         {fmt.Sprint(s.Cfg.EnableCohort)},
	}.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var layers []*Layer
	if err = json.Unmarshal(body, &layers); err != nil {
		return err
	}
	if s.cohortService != nil {
		s.cohortService.Switching()
	}
	sort.Sort(Layers(layers))
	for _, layer := range layers {
		if err = layer.Init(s.cohortService); err != nil {
			return err
		}
	}
	s.services = []ABService{NewUserService(layers), NewBucketService(layers)}
	if s.cohortService != nil {
		s.cohortService.Switched()
	}

	return nil
}

func (s *ABServices) AB(ctx *ABContext) *ABResult {
	result := NewABResult()
	for _, service := range s.services {
		service.AB(ctx, result)
	}
	return result
}
