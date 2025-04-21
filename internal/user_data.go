package internal

import "github.com/russianinvestments/invest-api-go-sdk/investgo"

type UserToken = string
type AccountId = string

type AccountIdWithAttachedClientttt struct {
	AccountId AccountId
	Client    *investgo.Client
}
