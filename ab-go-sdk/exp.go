package ab

type Exp struct {
	Name       string     `json:"name"`
	Versions   []*Version `json:"versions"`
	Conditions string     `json:"conditions"`
	Layer      *Layer     `json:"-"`
	Cohort     *Cohort    `json:"cohort"`
}
