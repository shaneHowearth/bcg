package twilio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	grpc "github.com/shanehowearth/bcg/notify/integration/grpc/proto/v1"
)

// Client -
type Client struct {
	AccountSid, AuthToken, Sender string
}

func NewClient() *Client {

	// SMS Account information
	accountSid, found := os.LookupEnv("TwilioAccountSID")
	if !found {
		log.Fatal("error Twilio specific AccountSID is required to be set in the externalProviders.env file, cannot continue")

	}
	authToken, found := os.LookupEnv("TwilioAuthToken")
	if !found {
		log.Fatal("error Twilio specific AuthToken is required to be set in the externalProviders.env file, cannot continue")

	}
	sender, found := os.LookupEnv("TwilioSenderNumber")
	if !found {
		log.Fatal("error Twilio specific SenderNumber is required to be set in the externalProviders.env file, cannot continue")

	}

	return &Client{
		AccountSid: accountSid,
		AuthToken:  authToken,
		Sender:     sender,
	}
}

// Send -
func (c Client) Send(recipient *grpc.CustomerDetails, message string) error {

	if c.AccountSid == "" || c.AuthToken == "" || c.Sender == "" {
		return fmt.Errorf("error Twilio specific AccountSID, AuthToken, and SenderNumber are required to be set in the externalProviders.env file, cannot continue")

	}
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + c.AccountSid + "/Messages.json"

	if recipient.Phone == "" {
		return fmt.Errorf("no phone number set for %s, cannot continue", recipient.Name)
	}

	// Pack up the data for our message
	msgData := url.Values{}
	msgData.Set("To", recipient.Phone)
	msgData.Set("From", c.Sender)
	msgData.Set("Body", message) // if no message supplied send "" anyway
	msgDataReader := *strings.NewReader(msgData.Encode())


	// Create HTTP request client
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		log.Printf("Error in New Request %v", err)
	}
	req.SetBasicAuth(c.AccountSid, c.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error in Client Do %v", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			log.Print("Message sent successfully")
			log.Print(data["sid"])
		}
	} else {
		log.Print("Message failure")
		log.Print(resp.Status)
	}
	return nil
}
