package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Repinoid/diploma56/internal/models"
	"github.com/Repinoid/diploma56/internal/securitate"
)

func GetBalance(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	UserID, err := securitate.Interbase.LoginByToken(rwr, req)
	if err != nil {
		return
	}

	current, withdr, err := securitate.Interbase.GetBalanceAndWithdrawn(req.Context(), UserID)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) // //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("row.Scan %+v\n", err)
		return
	}

	bs := models.BalanceStruct{Current: current - withdr, Withdrawn: withdr} // текущий счёт - сумма бонусов минус сумма списаний

	rwr.WriteHeader(http.StatusOK)
	json.NewEncoder(rwr).Encode(bs)
}
