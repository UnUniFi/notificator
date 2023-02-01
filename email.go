package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func sendMail(domain string, apiKey string, sender string, subject string, body string, recipient string) {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, apiKey)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
