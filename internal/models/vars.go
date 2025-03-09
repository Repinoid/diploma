package models

import (
	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger

type BalanceStruct struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
