package main

import (
	"log"
	"mxcat/internal/smtp"
)

func main() {
	client := smtp.New("127.0.0.1", 25)
	err := client.Start()
	if err != nil {
		log.Fatal(err)
	}
}
