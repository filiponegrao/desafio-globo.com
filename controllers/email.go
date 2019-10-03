package controllers

import (
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
)

type EmailConfiguration struct {
	Mail     string
	Password string
	Server   string
	Port     string
	Site     string
}

var mainConf EmailConfiguration

func ConfigEmailEngine(engine EmailConfiguration) {

	mainConf = EmailConfiguration{
		engine.Mail,
		engine.Password,
		engine.Server,
		engine.Port,
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

	address := conf.Server + ":" + conf.Port

	message := email.NewMessage(subject, text)
	sender := mail.Address{}
	sender.Address = conf.Mail
	sender.Name = "Suporte Convivva"
	message.From = sender
	message.To = []string{targetEmail}

	// err := message.Attach("logo1.png")
	// if err != nil {
	// 	log.Println(err)
	// }

	err := email.Send(address, auth, message)

	return err
}

// Teste
func EmailChangedPassword(targetEmail string, link string) error {

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
