package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Repinoid/diploma56/internal/models"
	"github.com/Repinoid/diploma56/internal/securitate"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var host = "localhost:8080"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()
	models.Sugar = *logger.Sugar()

	if err := initEnvs(); err != nil {
		panic(err)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var err error
	var Interbase securitate.Inter // переменная интерфейса, описание в internal/securitate/Inter.go
	var Accu *securitate.DBstruct
	ctx := context.Background()

	Interbase, err = securitate.ConnectToDB(ctx)

	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return err
	}

	Accu, err = securitate.ConnectToDB(ctx)
	if err != nil {
		fmt.Printf("database connection error  %v", err)
		return err
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", Interbase.RegisterUser).Methods("POST")
	router.HandleFunc("/api/user/login", Interbase.LoginUser).Methods("POST")
	router.HandleFunc("/api/user/balance/withdraw", Interbase.Withdraw).Methods("POST")

	router.HandleFunc("/api/user/orders", Interbase.PutOrder).Methods("POST")
	router.HandleFunc("/api/user/orders", Interbase.GetOrders).Methods("GET")
	router.HandleFunc("/api/user/withdrawals", Interbase.GetWithDrawals).Methods("GET")
	router.HandleFunc("/api/user/balance", Interbase.GetBalance).Methods("GET")

	go Accu.AccuOrders(ctx)

	return http.ListenAndServe(host, router)
}
