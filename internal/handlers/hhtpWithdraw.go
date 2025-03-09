package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Repinoid/diploma56/internal/models"
	"github.com/Repinoid/diploma56/internal/securitate"

	"github.com/theplant/luhn"
)

func Withdraw(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		rwr.WriteHeader(http.StatusBadRequest) //400 — неверный формат запроса; не text/plain
		fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
		models.Sugar.Debug("not text/plain \n")
		return
	}

	UserID, err := securitate.Interbase.LoginByToken(rwr, req)
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
	var orderID int64
	err = securitate.Interbase.GetIDByOrder(req.Context(), orderNum, &orderID)
	if err != nil { // если такого номера заказа нет в базе вносим его

		current, withdr, err := securitate.Interbase.GetBalanceAndWithdrawn(req.Context(), UserID)

		//	db := securitate.Interbase.DB
		// ordr := "SELECT (SELECT SUM(orders.accrual) FROM orders where orders.usercode=$1)- " +
		// 	"(SELECT COALESCE(SUM(withdrawn.amount),0) FROM withdrawn where withdrawn.usercode=$1) ;"
		// row := db.QueryRow(req.Context(), ordr, UserID) //
		// var accs float64                                // денег на счету
		// err := row.Scan(&accs)
		if err != nil {
			rwr.WriteHeader(http.StatusUnprocessableEntity) // 422 — неверный формат номера заказа;
			fmt.Fprintf(rwr, `{"status":"StatusUnprocessableEntity"}`)
			models.Sugar.Debugf("422 — невернная сумма на списание; %d\n", wdrStruct.Sum)
			return
		}
		if wdrStruct.Sum > current-withdr { // денег на счету
			rwr.WriteHeader(http.StatusPaymentRequired) //402 Payment Required
			fmt.Fprintf(rwr, `{"status":"StatusPaymentRequired"}`)
			models.Sugar.Debug("402 Payment Required\n")
			return
		}
		// -------------------------------------------------------------------------
		err = securitate.Interbase.AddToWithdrawn(req.Context(), UserID, orderNum, wdrStruct.Sum)
		//	ordr := "INSERT INTO withdrawn(userCode, orderNumber, amount) VALUES ($1, $2, $3) ;"
		//	_, err = db.Exec(req.Context(), ordr, UserID, orderNum, wdrStruct.Sum)
		if err != nil {
			rwr.WriteHeader(http.StatusInternalServerError) //500 — внутренняя ошибка сервера.
			fmt.Fprintf(rwr, `{"status":"StatusInternalServerError"}`)
			models.Sugar.Debug("error insert 2 withdrawn.\n")
			return
		}

		err = securitate.Interbase.UpLoadOrderByID(req.Context(), UserID, orderNum, "REGISTERED", 0)
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
	rwr.WriteHeader(http.StatusUnprocessableEntity) // 422 — неверный формат номера заказа;
	fmt.Fprintf(rwr, `{"status":"StatusUnprocessableEntity"}`)
	models.Sugar.Debug("422 — неверный формат номера заказа;\n")

}
