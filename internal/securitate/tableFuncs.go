package securitate

import (
	"context"
	"time"

	"github.com/Repinoid/diploma56/internal/models"
	"github.com/Repinoid/diploma56/internal/rual"
)

func (dataBase *DBstruct) TryWithdraw(ctx context.Context, UserID, orderNum int64, howmuch float64) (notEnough bool, err error) {
	db := dataBase.DB

	tx, err := db.Begin(ctx)
	if err != nil {
		models.Sugar.Debugf("error db.Begin  %[1]w", err)
	}
	ordr := "INSERT INTO withdrawn(userCode, orderNumber, amount) VALUES ($1, $2, $3) ;" // добавить в withdrawn сумму списания
	_, err = tx.Exec(ctx, ordr, UserID, orderNum, howmuch)
	if err != nil {
		models.Sugar.Debugf("tx.Exec %+v\n", err)
		return false, err
	}
	order := "SELECT (SELECT SUM(orders.accrual) FROM orders where orders.usercode=$1) - " + // получить разницу суммы кешбеков и списаний
		"(SELECT COALESCE(SUM(withdrawn.amount),0) FROM withdrawn  where withdrawn.usercode=$1) ;"
	var val float64
	row := tx.QueryRow(ctx, order, UserID)
	err = row.Scan(&val)
	if err != nil {
		models.Sugar.Debugf("row.Scan %+v\n", err)
		return false, err
	}
	if val < 0 { // если бабла недостаточно откатываем транзакцию
		tx.Rollback(ctx)
		return true, err // true - no $$$
	}
	err = tx.Commit(ctx)
	if err != nil {
		models.Sugar.Debugf(" tx.Commit %+v\n", err)
		return false, err
	}
	return false, err
}

func (dataBase *DBstruct) GetBalanceAndWithdrawn(ctx context.Context, UserID int64) (current, withdr float64, err error) {
	db := dataBase.DB

	order := "SELECT (SELECT SUM(orders.accrual) FROM orders where orders.usercode=$1), " +
		"(SELECT COALESCE(SUM(withdrawn.amount),0) FROM withdrawn  where withdrawn.usercode=$1) ;"

	row := db.QueryRow(ctx, order, UserID)
	err = row.Scan(&current, &withdr)
	return

}

func (dataBase *DBstruct) OrdersList(ctx context.Context, UserID int64) (orda []OrdStruct, err error) {

	db := dataBase.DB
	order := "select ordernumber as number, orderstatus as status, accrual, uploaded_at from orders where usercode=$1 order by uploaded_at ;"
	rows, err := db.Query(ctx, order, UserID) //
	if err != nil {
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}
	ord := OrdStruct{}
	for rows.Next() {
		var tm time.Time
		err = rows.Scan(&ord.Number, &ord.Status, &ord.Accrual, &tm)
		ord.UploadedAt = tm.Format(time.RFC3339)
		if err != nil {
			return
		}
		orda = append(orda, ord)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil { // Err returns any error that occurred while reading. Err must only be called after the Rows is closed
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}
	return
}

func (dataBase *DBstruct) WithdrawalsList(ctx context.Context, UserID int64) (orda []WithStruct, err error) {

	db := dataBase.DB
	order := "select ordernumber as number, amount as sum, processed_at from withdrawn where usercode=$1 order by processed_at ;"

	rows, err := db.Query(ctx, order, UserID) //
	if err != nil {
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}

	ord := WithStruct{}
	for rows.Next() {
		var tm time.Time
		err = rows.Scan(&ord.Order, &ord.Sum, &tm)
		ord.ProcessedAt = tm.Format(time.RFC3339)
		if err != nil {
			return
		}
		orda = append(orda, ord)
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil { // Err returns any error that occurred while reading. Err must only be called after the Rows is closed
		models.Sugar.Debugf("db.Query %+v\n", err)
		return
	}
	return
}

func (dataBase *DBstruct) AccuOrders(ctx context.Context) (err error) {

	db := dataBase.DB
	order := "select ordernumber, orderstatus, accrual from orders WHERE orderstatus != 'INVALID' AND orderstatus != 'PROCESSED' AND orderstatus != 'WITHDRAWN' ;"

	for {

		ord := OrdStruct{}
		orda := []OrdStruct{}
		rows, err := db.Query(ctx, order) //
		if err != nil {
			models.Sugar.Debugf("db.Query %+v\n", err)
			return err
		}
		for rows.Next() {
			err = rows.Scan(&ord.Number, &ord.Status, &ord.Accrual)
			if err != nil {
				return err
			}
			orda = append(orda, ord)
		}
		defer rows.Close()
		err = rows.Err()
		if err != nil { // Err returns any error that occurred while reading. Err must only be called after the Rows is closed
			models.Sugar.Debugf("db.Query %+v\n", err)
			return err
		}

		for _, ord := range orda {
			accuOrderStat, status, err := rual.GetFromAccrual(ord.Number)
			if err != nil {
				models.Sugar.Debugf("GetFromAccrual stat %d %v", status, err)
				continue
			}
			updateOrder := "UPDATE orders SET orderStatus = $2, accrual = $3 WHERE orderNumber  = $1 ;"
			_, err = db.Exec(ctx, updateOrder, ord.Number, accuOrderStat.Status, accuOrderStat.Accrual)
			if err != nil {
				models.Sugar.Debugf("error exex %v", err)
				return err
			}

		}

		time.Sleep(time.Second) // секунда задержки ... а сколько надо ставить ? или надо запускать после/перед каждым http обращением к таблице заказов ? триггеря через канал, например
	}
}
