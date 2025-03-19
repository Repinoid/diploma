package securitate

import (
	"context"
	"net/http"

	pgx "github.com/jackc/pgx/v5"
)

type DBstruct struct {
	DB     *pgx.Conn
	UserID int64
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
	CloseBase(ctx context.Context) error
}

type Handlera interface {
	RegisterUser(rwr http.ResponseWriter, req *http.Request)
	LoginUser(rwr http.ResponseWriter, req *http.Request)
	Withdraw(rwr http.ResponseWriter, req *http.Request)
	PutOrder(rwr http.ResponseWriter, req *http.Request)
	GetOrders(rwr http.ResponseWriter, req *http.Request)
	GetWithDrawals(rwr http.ResponseWriter, req *http.Request)
	GetBalance(rwr http.ResponseWriter, req *http.Request)
}

type Inter interface {
	Handlera
	TableCreations
	AddUser(ctx context.Context, userName, password, tokenString string) error
	CheckUserPassword(ctx context.Context, userName, password string) error
	IfUserExists(ctx context.Context, userName string) error
	ChangePassword(ctx context.Context, userName string, password string) error
	UpdateToken(ctx context.Context, userName string, tokenString string) error
	GetToken(ctx context.Context, userName string) (string, error)
	UpLoadOrderByID(ctx context.Context, userID int64, orderNumber int64, orderStatus string) error
	GetIDByOrder(ctx context.Context, orderNum int64) (int64, error)
	LoginByToken(rwr http.ResponseWriter, req *http.Request) (int64, error)

	OrdersList(ctx context.Context, UserID int64) (orda []OrdStruct, err error)
	WithdrawalsList(ctx context.Context, UserID int64) (orda []WithStruct, status int, err error)
	GetBalanceAndWithdrawn(ctx context.Context, UserID int64) (current, withdr float64, err error)
	TryWithdraw(ctx context.Context, UserID, orderNum int64, howmuch float64) (notEnough bool, err error)

	AccuOrders(ctx context.Context) (err error)
}
