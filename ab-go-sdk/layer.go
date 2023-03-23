package ab

import "strconv"

const (
	bucketSize = 1000
)

type Layer struct {
	Name          string `json:"name"`
	Platform      string `json:"platform"`
	Exps          []*Exp `json:"exps"`
	ShuffleTs     string `json:"shuffle_ts"`
	shufflePrefix string
	BucketsMap    map[string]string `json:"buckets"`
	Zone          string            `json:"zone"`
	Buckets       []*Version        `json:"-"`
}

func (this *Layer) Init(cohortService *CohortService) error {
	this.Buckets = make([]*Version, bucketSize)
	this.shufflePrefix = this.Name + this.ShuffleTs
	versionMap := map[string]*Version{}
	for _, exp := range this.Exps {
		if exp == nil {
			continue
		}
		exp.Layer = this
		if len(exp.Versions) == 0 {
			continue
		}
		if cohortService != nil {
			var err error
			exp.Cohort, err = cohortService.GetCohort(exp.Cohort)
			if err != nil {
				return err
			}
		}
		for _, version := range exp.Versions {
			version.Exp = exp
			versionMap[exp.Name+"-"+version.Name] = version
		}
	}
	for bucket, versionName := range this.BucketsMap {
		version := versionMap[versionName]
		n, err := strconv.Atoi(bucket)
		if err != nil {
			continue
		}
		if n >= 0 && n < bucketSize && version != nil {
			this.Buckets[n] = version
		}
	}
	return nil
}

func (this *Layer) ForEachLayer(action func(*Layer) bool) bool {
	return action(this)
}

func (this *Layer) GetBucket(n uint32) *Version {
	return this.Buckets[n%bucketSize]
}
