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
		BidderMspTtdVideo,
		BidderMspLiftoff,
		BidderMspLiftoffNative,
	}
}

func MspAllBidderNames() []BidderName {
	core := CoreBidderNames()
	return append(core, mspBidderNames()...)
}
