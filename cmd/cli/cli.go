package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/danilomarques1/grpc-gopm/pb"
)

const (
	SAVE   = "save"
	FIND   = "find"
	KEYS   = "keys"
	DELETE = "delete"
	UPDATE = "update"
	CLEAR  = "clear"
)

type Cli struct {
	passwordClient pb.PasswordClient
	scanner        *bufio.Scanner
}

func NewCli(client pb.PasswordClient) *Cli {
	scanner := bufio.NewScanner(os.Stdin)
	return &Cli{passwordClient: client, scanner: scanner}
}

func (c *Cli) Shell() {
	for {
		fmt.Print(">> ")
		var input string
		if c.scanner.Scan() {
			input = c.scanner.Text()
		}
		cmd, arg, err := c.parseInput(input)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}

		switch cmd {
		case SAVE:
			args := strings.Split(arg, " ")
			if len(args) != 2 {
				fmt.Println("Wrong number of arguments")
				continue
			}
			key, pwd := args[0], args[1]
			req := &pb.CreatePasswordRequest{
				Key:      key,
				Password: pwd,
			}
			response, err := c.passwordClient.SavePassword(context.Background(), req)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}

			if !response.GetOK() {
				fmt.Println("Could not save the password")
				continue
			}
			fmt.Println("Password saved successfully")
		case KEYS:
			response, err := c.passwordClient.FindAllKeys(context.Background(), &pb.Empty{})
			if err != nil {
				fmt.Printf("ERROR: %v", err)
				continue
			}
			for _, key := range response.GetKeys() {
				fmt.Printf("- %v\n", key)
			}
		case FIND:
			args := strings.Split(arg, " ")
			if len(args) != 1 {
				fmt.Println("Wrong number of arguments")
				continue
			}
			key := args[0]
			req := &pb.FindPasswordRequest{Key: key}
			response, err := c.passwordClient.FindPassword(context.Background(), req)
			if err != nil {
				fmt.Printf("ERROR: %v", err)
				continue
			}

			fmt.Println(response.GetPassword())
		case DELETE:
			args := strings.Split(arg, " ")
			if len(args) != 1 {
				fmt.Println("Wrong number of arguments")
				continue
			}
			key := args[0]
			req := &pb.DeletePasswordRequest{Key: key}
			response, err := c.passwordClient.DeletePassword(context.Background(), req)
			if err != nil {
				fmt.Printf("ERROR: %v", err)
				continue
			}
			if !response.GetOK() {
				fmt.Println("Could not remove the password")
				continue
			}
			fmt.Println("Password removed successfully")
		case UPDATE:
			args := strings.Split(arg, " ")
			if len(args) != 2 {
				fmt.Println("Wrong number of arguments")
				continue
			}
			key, pwd := args[0], args[1]
			req := &pb.UpdatePasswordRequest{Key: key, Password: pwd}
			response, err := c.passwordClient.UpdatePassword(context.Background(), req)
			if err != nil {
				fmt.Printf("ERROR: %v", err)
				continue
			}
			if !response.GetOK() {
				fmt.Println("Could not update the password")
				continue
			}
			fmt.Println("Password updated successfully")
		case CLEAR:
			operatingSystem := runtime.GOOS
			switch operatingSystem {
			case "windows":
				exec.Command("cls").Run()
			default:
				clearCommand := exec.Command("clear")
				out, err := clearCommand.Output()
				if err != nil {
					continue
				}
				fmt.Print(string(out))
			}

		default:
			fmt.Println("Command not found")
		}
	}
}

func (c *Cli) parseInput(input string) (string, string, error) {
	inputs := strings.Split(input, " ")
	if len(inputs) == 0 {
		return "", "", errors.New("Invalid command")
	}

	var cmd, arg string
	cmd = inputs[0]
	if len(inputs) > 1 {
		arg = strings.Join(inputs[1:], " ")
	}
	return cmd, arg, nil
}
