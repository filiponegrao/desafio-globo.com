package controllers

import (
	"log"
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
)

type EmailConfiguration struct {
	IP       string
	Port     string
	Mail     string
	Password string
	Server   string
	MailPort string
	Site     string
}

var mainConf EmailConfiguration

func ConfigEmailEngine(engine EmailConfiguration) {

	mainConf = EmailConfiguration{
		engine.IP,
		engine.Port,
		engine.Mail,
		engine.Password,
		engine.Server,
		engine.MailPort,
		engine.Site,
	}
}

func (conf EmailConfiguration) sendEmail(targetEmail string, text string, subject string) error {

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		conf.Mail,
		conf.Password,
		conf.Server,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	address := conf.Server + ":" + conf.MailPort

	message := email.NewMessage(subject, text)
	sender := mail.Address{}
	sender.Address = conf.Mail
	sender.Name = "Suporte MyBookmarks"
	message.From = sender
	message.To = []string{targetEmail}

	// err := message.Attach("logo1.png")

	err := email.Send(address, auth, message)
	if err != nil {
		log.Println(err)
	}

	return err
}

func EmailChangedPassword(targetEmail string, path string) error {

	link := mainConf.IP + ":" + mainConf.Port + path

	message := "MyBookmarks informa:\n\n"
	message += "Voce solicitou uma troca de senha!\n\n"
	message += "Para efetuar a troca acesse o link: "
	message += link
	message += "\n\nCaso o link nâo funcione entre em contanto em "
	message += mainConf.Site
	message += "\n\nAtenciosamente,\n"
	message += "Equipe MyBookmarks"

	subject := "[MyBookmarks]: Alteração de senha!"

	return mainConf.sendEmail(targetEmail, message, subject)
}
