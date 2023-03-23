package ab

type ABResult struct {
	Exp    map[string]string
	Config map[string]string
	layers map[string]bool
}

func NewABResult() *ABResult {
	return &ABResult{
		Exp:    map[string]string{},
		Config: map[string]string{},
		layers: map[string]bool{},
	}
}

func (this *ABResult) MergeVersion(version *Version) {
	if version == nil {
		return
	}
	if version.Exp == nil || this.layers[version.Exp.Layer.Name] {
		return
	}
	this.layers[version.Exp.Layer.Name] = true
	this.Exp[version.Exp.Name] = version.Name
	for k, v := range version.Config {
		_, exists := this.Config[k]
		if !exists {
			this.Config[k] = v
		}
	}
}

func (this *ABResult) MergeVersions(versions []*Version) {
	for _, version := range versions {
		this.MergeVersion(version)
	}
}

func (this *ABResult) ContainsLayer(layerName string) bool {
	return this.layers[layerName]
}
