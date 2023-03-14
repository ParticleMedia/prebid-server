package ab

type ZoneLayerService struct {
	Left  *ZoneLayerService
	Selfs []LayerService
}

func NewZoneLayerService() *ZoneLayerService {
	return &ZoneLayerService{}
}

func (this *ZoneLayerService) ForEachLayer(action func(layer *Layer) bool) bool {
	hit := false
	for _, self := range this.Selfs {
		hit = self.ForEachLayer(action) || hit
	}
	return hit || this.Left != nil && this.Left.ForEachLayer(action)
}
