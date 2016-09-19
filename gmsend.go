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

type EmailUser struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	EmailServer string `json:"emailserver"`
	Port        int    `json:"port"`
}

type Message struct {
	HideRecipients bool
	From           string
	Subject        string
	Content        string
}

func defaultEmailUser() (*EmailUser, error) {

	var result EmailUser

	prioritizedLocations := []string{
		os.Getenv("GMAIL_USER_AUTHENTICATION"),
		fmt.Sprintf("%s/.gmail_user_authentication.json", os.Getenv("HOME")),
		"/opt/etc/gmail_user_authentication.json"}

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

	return nil, fmt.Errorf("no default settings available")
}

func Send(emailUser *EmailUser, message Message, recipients []string) error {

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
		ToHeader = "To: " + strings.Join(recipients, ",") + "\n"
	}

	text := fmt.Sprintf("%s%s%s\n%s\n",
		FromHeader,
		ToHeader,
		SubjectHeader,
		message.Content)

	fmt.Printf("%s\n", text)

	serverPort := fmt.Sprintf("%s:%d",
		emailUser.EmailServer,
		emailUser.Port)

	return smtp.SendMail(serverPort,
		auth,
		emailUser.Username,
		recipients,
		[]byte(text))

}