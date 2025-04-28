package internal

import (
	"dp/pkg/nullable"
	"time"
)

type Operation struct {
	Type        OperationType
	Description string
	Time        time.Time
	Payment     nullable.Nullable[Money]
}

type Deal struct {
	Instrument Instrument
	Price      Money
	Quantity   Quantity
}

type DealOperation struct {
	Operation Operation
	Deal      nullable.Nullable[Deal]
}
