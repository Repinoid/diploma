package securitate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Repinoid/diploma56/internal/models"
	"github.com/jackc/pgx/v5"

	"github.com/theplant/luhn"
)

// /api/user/balance/withdraw
func (dataBase *DBstruct) Withdraw(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса; не text/plain
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debug("not text/plain \n")
		return
	}

	UserID, err := dataBase.LoginByToken(rwr, req)
	if err != nil {
		return
	}

	telo, err := io.ReadAll(req.Body)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
		fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
		models.Sugar.Debugf("io.ReadAll %+v\n", err)
		return
	}
	defer req.Body.Close()

	wdrStruct := struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}{}
	errM := json.Unmarshal([]byte(telo), &wdrStruct)

	orderNum, err := strconv.ParseInt(wdrStruct.Order, 10, 64)     //
	if err != nil || errM != nil || (!luhn.Valid(int(orderNum))) { // если не распарсилось или не по ЛУНУ
		rwr.WriteHeader(http.StatusUnprocessableEntity) // 422 — неверный формат номера заказа;
		fmt.Fprintf(rwr, `{"status":"StatusUnprocessableEntity"}`)
		models.Sugar.Debugf("422 — неверный формат номера заказа; %d\n", orderNum)
		return
	}
	//	var orderID int64
	_, err = dataBase.GetUserIDByOrder(req.Context(), orderNum)
	if err == nil { // // если заказ уже есть
		rwr.WriteHeader(http.StatusUnprocessableEntity) // 422 — неверный формат номера заказа;
		fmt.Fprintf(rwr, `{"status":"StatusUnprocessableEntity"}`)
		models.Sugar.Debug("422 — неверный формат номера заказа;\n")
	}
	if errors.Is(err, pgx.ErrNoRows) { //если запись не найдена

		// -------------------------------------------------------------------------
		noMany, err := dataBase.TryWithdraw(req.Context(), UserID, orderNum, wdrStruct.Sum)

		if err != nil {
			rwr.WriteHeader(http.StatusUnprocessableEntity) // 422 — неверный формат номера заказа;
			fmt.Fprintf(rwr, `{"status":"StatusUnprocessableEntity"}`)
			models.Sugar.Debugf("422 — невернная сумма на списание; %d\n", wdrStruct.Sum)
			return
		}
		//		if wdrStruct.Sum > current-withdr { // денег на счету
		if noMany { // если бабла было недостаточно
			rwr.WriteHeader(http.StatusPaymentRequired) //402 Payment Required
			fmt.Fprintf(rwr, `{"status":"StatusPaymentRequired"}`)
			models.Sugar.Debug("402 Payment Required\n")
			return
		}

		err = dataBase.UpLoadOrderByID(req.Context(), UserID, orderNum, "WITHDRAWN") // Лучше использовать вообще иной статус. Например, WITHDRAWN.
		if err != nil {
			rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
			fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
			models.Sugar.Debug("500 — внутренняя ошибка сервера.\n")
			return
		}
		//		}
		rwr.WriteHeader(http.StatusOK) //
		fmt.Fprintf(rwr, `{"status":"StatusOK"}`)
		return
	}
	rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
	fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
	models.Sugar.Debug("500 — внутренняя ошибка сервера.\n")

}
