package common_mappers

import (
	"dp/internal"
	"dp/internal/client"
)

func MapInstrument(instrument client.Instrument) internal.Instrument {
	return internal.Instrument{
		Figi:     internal.Figi(instrument.Figi),
		Name:     instrument.Name,
		LogoPath: instrument.LogoPath,
	}
}
