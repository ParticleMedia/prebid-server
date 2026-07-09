package openrtb_ext

func mspBidderNames() []BidderName {
	return []BidderName{
		BidderMspGoogle,
		BidderMspNova,
		BidderMspNovaAlpha,
		BidderMspNovaBeta,
		BidderMspNovaGamma,
		BidderMspFbAlpha,
		BidderMspFbBeta,
		BidderMspFbGamma,
		BidderMspFbDelta,
		BidderMspMoloco,
		BidderMspMolocoNative,
		BidderMspBidease,
		BidderMspTtdVideo,
		BidderMspVungleAlpha,
		BidderMspVungleBeta,
		BidderMspVungleGamma,
		BidderMspVungleDelta,
	}
}

func MspAllBidderNames() []BidderName {
	core := CoreBidderNames()
	return append(core, mspBidderNames()...)
}
