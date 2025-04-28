package client

import (
	"dp/internal/client"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func mapPositionInfo(resp *investapi.InstrumentResponse) client.Instrument {
	return client.Instrument{
		Figi:     resp.GetInstrument().GetFigi(),
		Name:     resp.GetInstrument().GetName(),
		LogoPath: resp.GetInstrument().GetBrand().GetLogoName(),
	}
}
