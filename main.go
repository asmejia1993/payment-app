package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/asmejia1993/payment-app/api"
	db "github.com/asmejia1993/payment-app/db/sqlc"
	"github.com/asmejia1993/payment-app/db/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Printf("connection: %s", config.DBSource)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB:", err)
	}
	if err := conn.Ping(); err != nil {
		fmt.Printf("error connecting to db: %v", err)
		log.Fatal("failed to ping db")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("failed to create server", err)
	}

	if err := server.Run(config.ServerAddress); err != nil {
		log.Fatal("failed to start server")
	}
}
