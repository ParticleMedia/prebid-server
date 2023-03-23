package ab

type LayerService interface {
	ForEachLayer(func(*Layer) bool) bool
}
