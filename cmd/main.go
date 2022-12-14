package main

import (
	"log"

	"github.com/danilomarques1/grpc-gopm/cmd/cli"
	"github.com/danilomarques1/grpc-gopm/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewPasswordClient(conn)
	cli := cli.NewCli(client)
	cli.Shell()
}
