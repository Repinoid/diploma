package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Repinoid/diploma56/internal/rual"
	"github.com/Repinoid/diploma56/internal/securitate"
)

func initEnvs() error {
	enva, exists := os.LookupEnv("RUN_ADDRESS")
	if exists {
		host = enva
		fmt.Printf("LookupEnv(RUN_ADDRESS)   %s \n", enva) // хотел определить откуда тесты берут параметры. оказалось из ENVs
	}
	enva, exists = os.LookupEnv("DATABASE_URI")
	if exists {
		securitate.DBEndPoint = enva
		fmt.Printf("LookupEnv(DATABASE_URI)   %s \n", enva)
	}
	enva, exists = os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if exists {
		rual.Accrualhost = enva
		fmt.Printf("LookupEnv(ACCRUAL_SYSTEM_ADDRESS)   %s \n", enva)
	}

	var hostFlag, dbFlag, acchostFlag string
	flag.StringVar(&hostFlag, "a", host, "Only -a={host:port} flag is allowed here")
	flag.StringVar(&dbFlag, "d", securitate.DBEndPoint, "Only -a={host:port} flag is allowed here")
	flag.StringVar(&acchostFlag, "r", rual.Accrualhost, "Only -a={host:port} flag is allowed here")
	flag.Parse()

	securitate.DBEndPoint = "postgres://postgres:passwordas@localhost:5432/forgo"

	if _, exists := os.LookupEnv("RUN_ADDRESS"); !exists { // если закомментить все IF ниже - у флагов будет приоритет перед переменными окружения
		host = hostFlag
		fmt.Printf("flag host   %s \n", host)
	}
	if _, exists := os.LookupEnv("DATABASE_URI"); !exists {
		securitate.DBEndPoint = dbFlag
		fmt.Printf("dbase   %s \n", dbFlag)
	}
	if _, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); !exists {
		rual.Accrualhost = acchostFlag
		fmt.Printf("accccc   %s \n", acchostFlag)
	}

	return nil
}
