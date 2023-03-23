package ab

type ABService interface {
	AB(*ABContext, *ABResult)
}

var ab *ABServices

func Init(cfg *ABConfig) {
	if ab == nil {
		ab = NewABServices(cfg)
	}
}

func AB(ctx *ABContext) *ABResult {
	return ab.AB(ctx)
}
