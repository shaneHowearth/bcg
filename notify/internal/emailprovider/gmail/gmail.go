package gmail

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Client -
type Client struct{}

// Send -
func (c Client) Send(recipient, message, sender string) error {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := c.getClient(config)
	srv, err := gmail.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	data := fmt.Sprintf(`From: %s
To: %s
Subject: Confirmation Email

%s
`, sender, recipient, message)
	emailMsg := b64.URLEncoding.EncodeToString([]byte(data))
	sndCall, err := srv.Users.Messages.Send(user, &gmail.Message{Raw: emailMsg}).Do()
	fmt.Printf("Send Call response: %#+vi, %v", sndCall, err)
	return err
}

// Retrieve a token, saves the token, then returns the generated client.
func (c Client) getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens
	tokFile := "token.json"
	tok := &oauth2.Token{}
	f, err := os.Open(tokFile)
	if err != nil {
		log.Panic("No token file present, cannot continue")
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		log.Panicf("Malformed token file, cannot continue, got error: %v", err)
	}
	return config.Client(context.Background(), tok)
}
