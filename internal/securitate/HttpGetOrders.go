package securitate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Repinoid/diploma56/internal/models"
)

// api/user/orders
func (dataBase *DBstruct) GetOrders(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	UserID, err := dataBase.LoginByToken(rwr, req)
	if err != nil {
		return
	}

	orda, err := dataBase.OrdersList(req.Context(), UserID)

	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}

	if len(orda) == 0 {
		rwr.WriteHeader(http.StatusNoContent) // 204 No Content — сервер успешно обработал запрос, но в ответе были переданы только заголовки без тела сообщения
		fmt.Fprintf(rwr, `{"status":"StatusNoContent"}`)
		models.Sugar.Debug("No ORDERS\n")
		return
	}
	rwr.WriteHeader(http.StatusOK)
	models.Sugar.Debugf("orda[0].Status  \"%+v\"\n", orda[0].Status)
	json.NewEncoder(rwr).Encode(orda)
}
