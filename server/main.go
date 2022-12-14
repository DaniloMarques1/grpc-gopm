package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/danilomarques1/grpc-gopm/pb"
	"github.com/danilomarques1/grpc-gopm/server/repository"
	"github.com/danilomarques1/grpc-gopm/server/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

const TABLES = `
CREATE TABLE IF NOT EXISTS passwords(
	id varchar(40) primary key,
	key varchar(100) not null unique,
	password varchar(100) not null
);
`

func main() {
	db, err := sql.Open("sqlite3", "db.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(TABLES); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	passwordRepository := repository.NewPasswordRepositoryImpl(db)
	server := grpc.NewServer()
	pb.RegisterPasswordServer(server, service.NewPasswordServer(passwordRepository))

	log.Printf("Starting grpc server\n")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
