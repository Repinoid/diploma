package securitate

import (
	"context"
	"fmt"
)

func (dataBase *DBstruct) UsersTableCreation(ctx context.Context) error {

	db := dataBase.DB
	db.Exec(ctx, "CREATE EXTENSION pgcrypto;") // расширение для хэширования паролей

	creatorOrder :=
		"CREATE TABLE IF NOT EXISTS " + "accounts" +
			"(userCode INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY ," +
			"login VARCHAR(100) UNIQUE," +
			"password VARCHAR(200) NOT NULL," +
			"user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"

	_, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("create users table. %w", err)
	}
	return nil
}
func (dataBase *DBstruct) OrdersTableCreation(ctx context.Context) error {
	db := dataBase.DB
	creatorOrder :=
		"CREATE TABLE IF NOT EXISTS " + "orders" +
			"(id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY," +
			"userCode INT NOT NULL," +
			"orderNumber BIGINT NOT NULL UNIQUE," +
			"orderStatus VARCHAR(20)," +
			"accrual FLOAT8," +
			"uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
			"FOREIGN KEY (userCode) REFERENCES " + "accounts" + "(usercode) ON DELETE CASCADE);"

	_, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("create orders table. %w", err)
	}
	return nil
}

func (dataBase *DBstruct) TokensTableCreation(ctx context.Context) error {
	db := dataBase.DB
	creatorOrder :=
		"CREATE TABLE IF NOT EXISTS " + "tokens" +
			"(id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY," +
			"userCode INT NOT NULL UNIQUE," +
			//			"balance FLOAT8 DEFAULT 0," +
			//			"bonus FLOAT8 DEFAULT 0," +
			"token VARCHAR(1000) NOT NULL," +
			"token_valid_until TIMESTAMP," +
			"token_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
			"FOREIGN KEY (userCode) REFERENCES " + "accounts" + "(usercode) ON DELETE CASCADE);"
	_, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("create orders table. %w", err)
	}
	return nil
}
func (dataBase *DBstruct) WithdrawalsTableCreation(ctx context.Context) error {
	db := dataBase.DB
	creatorOrder :=
		"CREATE TABLE IF NOT EXISTS " + "withdrawn" +
			"(id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY," +
			"userCode INT NOT NULL," +
			"orderNumber BIGINT NOT NULL UNIQUE," +
			"amount FLOAT8 DEFAULT 0," +
			"processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
			"FOREIGN KEY (userCode) REFERENCES " + "accounts" + "(usercode) ON DELETE CASCADE);"
	_, err := db.Exec(ctx, creatorOrder)
	if err != nil {
		return fmt.Errorf("create orders table. %w", err)
	}
	return nil
}
