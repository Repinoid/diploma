package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Repinoid/diploma56/internal/models"
	"github.com/Repinoid/diploma56/internal/securitate"
)

func GetWithDrawals(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	UserID, err := securitate.Interbase.LoginByToken(rwr, req)
	if err != nil {
		return
	}

	orda, status, err := securitate.Interbase.WithdrawalsList(req.Context(), UserID)

	if status == http.StatusInternalServerError {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}

	if len(orda) == 0 {
		rwr.WriteHeader(http.StatusNoContent) // 204 No Content — сервер успешно обработал запрос, но в ответе были переданы только заголовки без тела сообщения
		fmt.Fprintf(rwr, `{"status":"StatusNoContent"}`)
		models.Sugar.Debug("No withdrawals\n")
		return
	}
	rwr.WriteHeader(http.StatusOK)
	json.NewEncoder(rwr).Encode(orda)
}
