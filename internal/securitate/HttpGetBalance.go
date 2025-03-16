package securitate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Repinoid/diploma56/internal/models"
)

// api/user/balance
func (dataBase *DBstruct) GetBalance(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	UserID, err := dataBase.LoginByToken(rwr, req)
	if err != nil {
		return
	}

	current, withdr, err := dataBase.GetBalanceAndWithdrawn(req.Context(), UserID)
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
