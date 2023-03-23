package ab

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/RoaringBitmap/roaring"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
)

type CohortService struct {
	Current map[string]*Cohort
	NG      map[string]*Cohort
	url     string
}

func NewCohortService(url string) *CohortService {
	return &CohortService{
		Current: map[string]*Cohort{},
		NG:      map[string]*Cohort{},
		url:     url,
	}
}

func (u *CohortService) Switching() {
	u.NG = map[string]*Cohort{}
}

func (u *CohortService) Switched() {
	u.Current = u.NG
	u.NG = map[string]*Cohort{}
}

func (u *CohortService) GetCohort(input *Cohort) (tag *Cohort, err error) {
	if input == nil {
		return nil, nil
	}
	location := input.Location
	update := input.UpdateTime
	if tag = u.NG[location]; tag != nil && tag.UpdateTime >= update {
		return tag, nil
	}
	if tag = u.Current[location]; tag != nil && tag.UpdateTime >= update {
		u.NG[location] = tag
		return tag, nil
	}
	tag = &Cohort{
		Location:   location,
		Users:      roaring.New(),
		UpdateTime: update,
	}
	if err := u.setupAwsEnv(); err != nil {
		return nil, err
	}
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	var parsedUrl *url.URL
	if parsedUrl, err = url.Parse(location); err != nil {
		return nil, err
	}
	s3Object, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(parsedUrl.Host),
		Key:    aws.String(parsedUrl.Path),
	})
	if err != nil {
		return nil, fmt.Errorf("Get Cohort: %s, err: %s", location, err)
	}
	defer s3Object.Body.Close()
	scanner := bufio.NewScanner(s3Object.Body)
	for scanner.Scan() {
		userID, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		tag.Users.AddInt(userID)
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	u.NG[location] = tag
	glog.Info("Reload Cohort success, location: " + location + ", update_time: " + update)
	return tag, nil
}

func (this *CohortService) setupAwsEnv() error {
	resp, err := http.Get(this.url + "/sys/aws_envs")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	envs := map[string]string{}
	if err = json.Unmarshal(body, &envs); err != nil {
		return err
	}
	for k, v := range envs {
		os.Setenv(k, v)
	}
	return nil
}
