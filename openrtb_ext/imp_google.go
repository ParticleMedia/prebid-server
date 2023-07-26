package openrtb_ext

type ExtImpGoogle struct {
	IU string `json:"iu"`
}

type GoogleContext struct {
	DataMap map[string][]string `json:"data"`
}

type ExtImpBidderGoogle struct {
	Context GoogleContext `json:"context"`
}

type GoogleBidExt struct {
	Detail GoogleBidExtDetail `json:"google"`
}

type GoogleBidExtDetail struct {
	AdUnit string `json:"ad_unit_id"`
}
