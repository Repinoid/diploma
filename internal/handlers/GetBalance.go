package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Repinoid/diploma56/internal/models"
	"github.com/Repinoid/diploma56/internal/securitate"
)

type BalanceStruct struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func GetBalance(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	UserID, err := securitate.DataBase.LoginByToken(rwr, req)
	if err != nil {
		return
	}

	order := "SELECT (SELECT SUM(orders.accrual) FROM orders where orders.usercode=$1), " +
		"(SELECT COALESCE(SUM(withdrawn.amount),0) FROM withdrawn  where withdrawn.usercode=$1) ;"

	row := securitate.DataBase.DB.QueryRow(context.Background(), order, UserID)
	var current, withdr float64
	err = row.Scan(&current, &withdr)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) // //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("row.Scan %+v\n", err)
		return
	}

	bs := BalanceStruct{Current: current - withdr, Withdrawn: withdr} // текущий счёт - сумма бонусов минус сумма списаний

	rwr.WriteHeader(http.StatusOK)
	json.NewEncoder(rwr).Encode(bs)
}
