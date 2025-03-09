package securitate

import (
	"context"
	"net/http"
	"time"

	"github.com/Repinoid/diploma56/internal/models"
)

func (dataBase *DBstruct) GetBalanceAndWithdrawn(ctx context.Context, UserID int64) (current, withdr float64, err error) {
	db := dataBase.DB

	order := "SELECT (SELECT SUM(orders.accrual) FROM orders where orders.usercode=$1), " +
		"(SELECT COALESCE(SUM(withdrawn.amount),0) FROM withdrawn  where withdrawn.usercode=$1) ;"

	row := db.QueryRow(ctx, order, UserID)
	err = row.Scan(&current, &withdr)
	return

}

func (dataBase *DBstruct) AddToWithdrawn(ctx context.Context, UserID, orderNum int64, sum float64) (err error) {
	db := dataBase.DB

	ordr := "INSERT INTO withdrawn(userCode, orderNumber, amount) VALUES ($1, $2, $3) ;"
	_, err = db.Exec(ctx, ordr, UserID, orderNum, sum)

	return

}

func (dataBase *DBstruct) OrdersList(ctx context.Context, UserID int64) (orda []OrdStruct, status int, err error) {

	db := dataBase.DB
	order := "select ordernumber as number, orderstatus as status, accrual, uploaded_at from orders where usercode=$1 order by uploaded_at ;"
	rows, err := db.Query(ctx, order, UserID) //
	if err != nil {
		status = http.StatusInternalServerError //500 — внутренняя ошибка сервера.
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}
	ord := OrdStruct{}
	//	orda := []models.OrdStruct{}
	var errScan error
	for rows.Next() {
		var tm time.Time
		errScan = rows.Scan(&ord.Number, &ord.Status, &ord.Accrual, &tm)
		ord.UploadedAt = tm.Format(time.RFC3339)
		if errScan != nil {
			break
		}
		orda = append(orda, ord)
	}
	rows.Close()

	if err = rows.Err(); err != nil || errScan != nil { // Err returns any error that occurred while reading. Err must only be called after the Rows is closed
		status = http.StatusInternalServerError // //500 — внутренняя ошибка сервера.
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}

	status = http.StatusOK
	return
}

func (dataBase *DBstruct) WithdrawalsList(ctx context.Context, UserID int64) (orda []WithStruct, status int, err error) {

	db := dataBase.DB
	order := "select ordernumber as number, amount as sum, processed_at from withdrawn where usercode=$1 order by processed_at ;"

	rows, err := db.Query(ctx, order, UserID) //
	if err != nil {
		status = http.StatusInternalServerError //500 — внутренняя ошибка сервера.
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}

	ord := WithStruct{}
	//	orda := []WithStruct{}
	var errScan error
	for rows.Next() {
		var tm time.Time
		errScan = rows.Scan(&ord.Order, &ord.Sum, &tm)
		ord.ProcessedAt = tm.Format(time.RFC3339)
		if errScan != nil {
			break
		}
		orda = append(orda, ord)
	}
	rows.Close()
	if err = rows.Err(); err != nil || errScan != nil { // Err returns any error that occurred while reading. Err must only be called after the Rows is closed
		status = http.StatusInternalServerError //500 — внутренняя ошибка сервера.
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}
	status = http.StatusOK
	return
}
