package main

import (
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func main() {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.Username = "test"
	server.Password = "test"
	server.Authentication = mail.AuthPlain
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
	}
	log.Print(smtpClient)
}