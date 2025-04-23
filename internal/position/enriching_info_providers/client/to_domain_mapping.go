package client

import (
	"dp/internal"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func mapEnrichingInfoToDomain(resp *investapi.InstrumentResponse) internal.PositionEnrichingInfo {
	return internal.PositionEnrichingInfo{
		Figi:     internal.Figi(resp.GetInstrument().GetFigi()),
		Name:     resp.GetInstrument().GetName(),
		LogoPath: resp.GetInstrument().GetBrand().GetLogoName(),
	}
}
