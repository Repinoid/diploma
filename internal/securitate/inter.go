package securitate

import (
	"context"
	"net/http"

	pgx "github.com/jackc/pgx/v5"
)

type DBstruct struct {
	DB *pgx.Conn
}

var DBEndPoint string

type OrdStruct struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

type WithStruct struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type TableCreations interface {
	UsersTableCreation(ctx context.Context) error
	OrdersTableCreation(ctx context.Context) error
	TokensTableCreation(ctx context.Context) error
	WithdrawalsTableCreation(ctx context.Context) error
}

type Inter interface {
	TableCreations
	AddUser(ctx context.Context, userName, password, tokenString string) error
	CheckUserPassword(ctx context.Context, userName, password string) error
	IfUserExists(ctx context.Context, userName string) error
	ChangePassword(ctx context.Context, userName string, password string) error
	UpdateToken(ctx context.Context, userName string, tokenString string) error
	GetToken(ctx context.Context, userName string, tokenString *string) error
	UpLoadOrderByID(ctx context.Context, userID int64, orderNumber int64, orderStatus string, accrual float64) error
	GetIDByOrder(ctx context.Context, orderNum int64, orderID *int64) error
	AddOrder(ctx context.Context, userName string, orderNumber int64, orderStatus string, accrual float64) error
	LoginByToken(rwr http.ResponseWriter, req *http.Request) (int64, error)

	OrdersList(ctx context.Context, UserID int64) (orda []OrdStruct, status int, err error)
	WithdrawalsList(ctx context.Context, UserID int64) (orda []WithStruct, status int, err error)
	GetBalanceAndWithdrawn(ctx context.Context, UserID int64) (current, withdr float64, err error)
	AddToWithdrawn(ctx context.Context, UserID, orderNum int64, sum float64) (err error)
}



var Interbase *DBstruct
//var Interbase Inter