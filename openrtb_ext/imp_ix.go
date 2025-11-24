package openrtb_ext

// ExtImpIx defines the contract for bidrequest.imp[i].ext.prebid.bidder.ix
type ExtImpIx struct {
	Floor    float64 `json:"floor"`
	SiteId   string  `json:"siteId"`
	Size     []int   `json:"size"`
	Sid      string  `json:"sid"`
	Endpoint string  `json:"ix_endpoint"`
}
