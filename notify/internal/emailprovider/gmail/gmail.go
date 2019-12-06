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
	srv, err := gmail.New(client) // TODO use NewService()

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
// func (c Client) getClient(config *oauth2.Config) *http.Client {
// 	// The file token.json stores the user's access and refresh tokens, and is
// 	// created automatically when the authorization flow completes for the first
// 	// time.
// 	tokFile := "token.json"
// 	tok, err := c.tokenFromFile(tokFile)
// 	if err != nil {
// 		tok = c.getTokenFromWeb(config)
// 		c.saveToken(tokFile, tok)
// 	}
// 	return config.Client(context.Background(), tok)
// }

// Retrieves a token from a local file.
// func (c Client) tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(tok)
// 	return tok, err
// }

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

// Saves a token to a file path.
// func (c Client) saveToken(path string, token *oauth2.Token) {
// 	fmt.Printf("Saving credential file to: %s\n", path)
// 	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		log.Fatalf("Unable to cache oauth token: %v", err)
// 	}
// 	defer f.Close()
// 	err = json.NewEncoder(f).Encode(token)
// 	if err != nil {
// 		log.Fatalf("Unable to encode oauth token: %v", err)
// 	}
// }

// Request a token from the web, then returns the retrieved token.
// func (c Client) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
// 	fmt.Printf("Go to the following link in your browser then type the "+
// 		"authorization code: \n%v\n", authURL)

// 	var authCode string
// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		log.Fatalf("Unable to read authorization code: %v", err)
// 	}

// 	tok, err := config.Exchange(context.TODO(), authCode)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve token from web: %v", err)
// 	}
// 	return tok
// }
