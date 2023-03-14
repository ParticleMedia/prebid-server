package ab

type Version struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
	Users  []string          `json:"users"`
	Exp    *Exp              `json:"-"`
}
