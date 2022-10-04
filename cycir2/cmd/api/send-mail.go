package main

import (
	"cycir/internal/channeldata"
	"log"
)

// SendEmail sends an email
func (app *application) SendEmail(mailMessage channeldata.MailData) {
	// if no sender specified, use defaults
	if mailMessage.FromAddress == "" {
		mailMessage.FromAddress = app.PreferenceMap["smtp_from_email"]
		mailMessage.FromName = app.PreferenceMap["smtp_from_name"]
	}
	log.Println(app.MailQueue)

	job := channeldata.MailJob{MailMessage: mailMessage}
	app.MailQueue <- job
}
