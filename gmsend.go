package gmsend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
)

//https://nathanleclaire.com/blog/2013/12/17/sending-email-from-gmail-using-golang/

// SMTPAuthentication contains gmail identification, server and port
type SMTPAuthentication struct {
	Username    string `json:"username"`    // e.g. "ralphmalph@gmail.com"
	Password    string `json:"password"`    // e.g. "mypassword1234!"
	EmailServer string `json:"emailserver"` // e.g. "smtp.gmail.com"
	Port        int    `json:"port"`        // e.g. 587
}

// Message contains content and subject
type Message struct {
	HideRecipients bool
	From           string
	Subject        string
	Content        string
}

func defaultEmailUser() (*SMTPAuthentication, error) {

	var result SMTPAuthentication

	prioritizedLocations := []string{
		os.Getenv("GMSEND"),
		fmt.Sprintf("%s/.gmsend.json", os.Getenv("HOME")),
		"/opt/etc/gmsend.json"}

	for _, location := range prioritizedLocations {

		if location == "" {
			continue
		}

		bytes, err := ioutil.ReadFile(location)
		if err != nil {
			continue
		}

		err = json.Unmarshal(bytes, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	return nil, fmt.Errorf("no default settings found in [$GMSEND,~/.gmsend.json,/opt/etc/gmsend.json")
}

// Send takes a *SMTPAuthentication (which is usually nil), a message
// and recipients.
func Send(emailUser *SMTPAuthentication, message Message, recipients []string) error {

	if emailUser == nil {
		defaulted, err := defaultEmailUser()
		if err != nil {
			return err
		}
		emailUser = defaulted
	}

	if len(recipients) == 0 {
		return fmt.Errorf("empty recipients list")
	}

	if message.Content == "" {
		return fmt.Errorf("empty message content")
	}

	if message.Subject == "" {
		return fmt.Errorf("no subject supplied")
	}

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer)

	if message.From == "" {
		message.From = emailUser.Username
	}

	FromHeader := fmt.Sprintf("From: %s\n", message.From)
	SubjectHeader := fmt.Sprintf("Subject: %s\n", message.Subject)

	var ToHeader string
	if !message.HideRecipients {
		ToHeader = fmt.Sprintf("To: %s\n", strings.Join(recipients, ","))
	}

	text := fmt.Sprintf("%s%s%s\n%s\n",
		FromHeader,
		ToHeader,
		SubjectHeader,
		message.Content)

	serverPort := fmt.Sprintf("%s:%d",
		emailUser.EmailServer,
		emailUser.Port)

	return smtp.SendMail(serverPort,
		auth,
		emailUser.Username,
		recipients,
		[]byte(text))

}
