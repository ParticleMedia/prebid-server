package vungle

type VungleBidExt struct {
	Detail VungleBidExtDetail `json:"vungle"`
}

type VungleBidExtDetail struct {
	PlacementReferenceId string `json:"placement_reference_id"`
}
