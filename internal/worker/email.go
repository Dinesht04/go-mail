package worker

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type Email struct {
	Content   string
	Recipient string
}

func sendEmail(email Email) bool {
	fmt.Println("sending mail?")
	auth := smtp.PlainAuth("", "dineshtyagi567@gmail.com", os.Getenv("smtp_pass"), "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{"dineshtyagi2022@vitbhopal.ac.in"}
	msg := []byte("To: dineshtyagi2022@vitbhopal.ac.in\r\n" +
		"Subject: From my server!\r\n" +
		"\r\n" +
		"This is a scheduled task.\r\n")
	err := smtp.SendMail("smtp.gmail.com:25", auth, "dineshtyagi567@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func GenerateOtp() {}

func VerifyOtp() {}

func Sendmessage() {}

func Unsubscribe() {}
