package main

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"
)


func main() {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.Username = ""
	server.Password = ""
	server.Authentication = mail.AuthPlain
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 5 * time.Second
	server.SendTimeout = 5 * time.Second
	
	_, err := server.Connect()
	if err != nil {
		log.Println(err)
	}
}