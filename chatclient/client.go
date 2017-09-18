package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("Unknown%d", rand.Intn(1000))

	fmt.Println("Starting chat...")
	fmt.Println("Please, tell me your name?")
	// In case of the scan failed for any reason, name will keep its previous value.
	fmt.Scan(&name)

	fmt.Println("Starting connection with the server...")
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("Could not connect to the server!", err)
	}
	fmt.Printf("Connection estabilished. You're on, %s!\n", name)
	name += ": "
	defer conn.Close()

	// Set up reader
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	//Set up Writer and block the execution
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := name + scanner.Text() + "\n"
		fmt.Fprint(conn, msg)
	}
}
