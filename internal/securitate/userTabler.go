package securitate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Repinoid/diploma56/internal/models"
	pgx "github.com/jackc/pgx/v5"
)



func ConnectToDB(ctx context.Context) (*DBstruct, error) {

	dataBase := &DBstruct{}
	baza, err := pgx.Connect(ctx, DBEndPoint)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DB %s err %w", DBEndPoint, err)
	}
	dataBase.DB = baza

	if err := dataBase.UsersTableCreation(ctx); err != nil {
		return nil, fmt.Errorf("UsersTableCreation %w", err)
	}
	if err := dataBase.OrdersTableCreation(ctx); err != nil {
		return nil, fmt.Errorf("OrdersTableCreation %w", err)
	}
	if err := dataBase.TokensTableCreation(ctx); err != nil {
		return nil, fmt.Errorf("TokensTableCreation %w", err)
	}
	if err := dataBase.WithdrawalsTableCreation(ctx); err != nil {
		return nil, fmt.Errorf("WithdrawalsTableCreation %w", err)
	}
	return dataBase, nil
}

func (dataBase *DBstruct) AddUser(ctx context.Context, userName, password, tokenString string) error {
	db := dataBase.DB

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error db.Begin  %[1]w", err)
	}
	defer tx.Rollback(ctx)

	order := "INSERT INTO " + "accounts" + " (login, password) VALUES ($1, crypt($2, gen_salt('md5'))) ;"
	_, err = tx.Exec(ctx, order, userName, password)
	if err != nil {
		return fmt.Errorf("add user error is %w", err)
	}
	order = "INSERT INTO tokens(userCode, token) VALUES ((select usercode from accounts where login = $1), $2) ;"
	_, err = tx.Exec(ctx, order, userName, tokenString)
	if err != nil {
		return fmt.Errorf("add TOKEN %w", err)
	}
	return tx.Commit(ctx)
}

func (dataBase *DBstruct) CheckUserPassword(ctx context.Context, userName, password string) error {
	db := dataBase.DB
	order := "SELECT (password = crypt($2, password)) AS password_match FROM " + "accounts" + " WHERE login= $1 ;"
	row := db.QueryRow(ctx, order, userName, password) // password here - what was entered
	var yes bool
	err := row.Scan(&yes)
	if err != nil {
		return fmt.Errorf("QueryRow, error is %w", err)
	}
	if !yes {
		return fmt.Errorf("password not match")
	}
	return nil
}

// nil - user exists
func (dataBase *DBstruct) IfUserExists(ctx context.Context, userName string) error {
	db := dataBase.DB
	order := "SELECT 7 from " + "accounts" + " WHERE login= $1 ;"
	row := db.QueryRow(ctx, order, userName) // password here - what was entered
	var yes int
	err := row.Scan(&yes)
	if err != nil {
		return fmt.Errorf(" QueryRow, error is %w", err)
	}
	if yes != 7 {
		return fmt.Errorf("user %s does not exist", userName)
	}
	return nil
}

func (dataBase *DBstruct) ChangePassword(ctx context.Context, userName string, password string) error {
	db := dataBase.DB
	order := "UPDATE " + "accounts" + " SET password = crypt($2, gen_salt('md5')) WHERE login= $1 ;"
	_, err := db.Exec(ctx, order, userName, password)
	if err != nil {
		return fmt.Errorf("change password error %w", err)
	}
	return nil
}

func (dataBase *DBstruct) UpdateToken(ctx context.Context, userName string, tokenString string) error {
	db := dataBase.DB
	order := "UPDATE tokens SET token = $2 WHERE userCode = (select usercode from accounts where login = $1) ;"
	_, err := db.Exec(ctx, order, userName, tokenString)
	if err != nil {
		return fmt.Errorf("add TOKEN %w", err)
	}
	return nil
}

func (dataBase *DBstruct) GetToken(ctx context.Context, userName string, tokenString *string) error {
	db := dataBase.DB
	//				получить токен из токен-таблицы  где код пользователя равен коду юзера из юзер-таблицы с именем UserName
	order := "SELECT token from " + "tokens" + " WHERE userCode = (select usercode from accounts where login = $1) ;"
	row := db.QueryRow(ctx, order, userName)
	var str string
	err := row.Scan(&str)
	if err != nil {
		return fmt.Errorf("GT %w", err)
	}
	*tokenString = str
	return nil
}

func (dataBase *DBstruct) UpLoadOrderByID(ctx context.Context, userID int64, orderNumber int64, orderStatus string, accrual float64) error {
	db := dataBase.DB
	if orderStatus == "" || orderStatus == "REGISTERED" {
		orderStatus = "NEW"
	}
	order := "INSERT INTO orders(userCode, orderNumber, orderStatus, accrual) VALUES ($1, $2, $3, $4) ;"
	_, err := db.Exec(ctx, order, userID, orderNumber, orderStatus, accrual)
	if err != nil {
		return fmt.Errorf("add ORDER %w", err)
	}
	return nil
}

func (dataBase *DBstruct) GetIDByOrder(ctx context.Context, orderNum int64, orderID *int64) error {
	db := dataBase.DB
	order := "SELECT usercode from " + "orders" + " WHERE orderNumber =  $1 ;"
	row := db.QueryRow(ctx, order, orderNum)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("GT %w", err)
	}
	*orderID = id
	return nil
}

func (dataBase *DBstruct) AddOrder(ctx context.Context, userName string, orderNumber int64, orderStatus string, accrual float64) error {
	db := dataBase.DB

	order := "INSERT INTO orders(userCode, ordernumber, orderStatus, accrual) VALUES ((select id from accounts where login = $1), $2, $3, $4) ;"

	_, err := db.Exec(ctx, order, userName, orderNumber, orderStatus, accrual)
	if err != nil {
		return fmt.Errorf("add ORDER %w", err)
	}
	return nil
}

func (dataBase *DBstruct) LoginByToken(rwr http.ResponseWriter, req *http.Request) (int64, error) {

	tokenStr := req.Header.Get("Authorization")
	tokenStr, niceP := strings.CutPrefix(tokenStr, "Bearer <") // обрезаем -- Bearer <token>
	tokenStr, niceS := strings.CutSuffix(tokenStr, ">")

	var UserID int64
	if niceP && niceS {
		order := "SELECT usercode from " + "tokens" + " WHERE token =  $1 ;"
		row := dataBase.DB.QueryRow(req.Context(), order, tokenStr)
		err := row.Scan(&UserID)
		if err == nil {
			return UserID, nil
		}
	}
	rwr.WriteHeader(http.StatusUnauthorized)            // 401 — неверная пара логин/пароль;
	fmt.Fprintf(rwr, `{"status":"StatusUnauthorized"}`) // либо токена неверный формат, либо по нему нет юзера в базе
	models.Sugar.Debug("Authorization header\n")
	return 0, errors.New("Unauthorized")
}
