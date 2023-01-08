package main

import (
	"cycir/internal/channeldata"
)

// SendEmail sends an email
func (app *application) SendEmail(mailMessage channeldata.MailData) {
	// if no sender specified, use defaults
	if mailMessage.FromAddress == "" {
		mailMessage.FromAddress = app.PreferenceMap["smtp_from_email"]
		mailMessage.FromName = app.PreferenceMap["smtp_from_name"]
	}

	job := channeldata.MailJob{MailMessage: mailMessage}
	app.MailQueue <- job
}
