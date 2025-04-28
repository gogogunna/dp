package internal

const (
	divider           = 10000000
	multiplier        = 100
	percentMultiplier = 100
)

type UnitsWithNano int

func MapUnitsWithNano(
	value interface {
		GetUnits() int64
		GetNano() int32
	},
) UnitsWithNano {
	units := value.GetUnits()
	nanos := int64(value.GetNano() / divider)
	return UnitsWithNano(units*multiplier + nanos)
}

func CalculatePercent(diff UnitsWithNano, all UnitsWithNano) Percent {
	return Percent(float64(diff) / float64(all) * multiplier * percentMultiplier)
}

func CalculateDiff(percent UnitsWithNano, all UnitsWithNano) Money {
	return Money(float64(percent) / multiplier / percentMultiplier * float64(all))
}
