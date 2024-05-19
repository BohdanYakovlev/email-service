package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	smtpHost = "smtp-mail.outlook.com"
	smtpPort = 587
	smtpUser = "rate2341@outlook.com"
	smtpPass = "3h9M9Ps7tA"
)

const rateUrl = "http://api-service:8001/api/rate"
const emailsUrl = "http://api-service:8001/api/emails"

func getMessage(rate string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("Subject", "USD to UAH rate")
	m.SetBody("text/plain", fmt.Sprintf("1 USD = %s UAH", rate))
	return m
}

func getRate() string {
	var client http.Client
	resp, err := client.Get(rateUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return string(bodyBytes)
	}
	return ""
}

func getEmails() []string {
	var client http.Client
	resp, err := client.Get(emailsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		res := new([]string)
		err = json.Unmarshal(bodyBytes, res)
		if err != nil {
			log.Fatal(err)
		}
		return *res
	}
	return []string{}
}

func main() {

	flag := true

	for flag {
		hour, _, _ := time.Now().Clock()
		if hour < 1 {
			rate := getRate()
			message := getMessage(rate)
			d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

			emails := getEmails()

			for _, recipient := range emails {
				message.SetHeader("To", recipient)
				if err := d.DialAndSend(message); err != nil {
					log.Printf("Failed to send email to %s: %v", recipient, err)
				} else {
					log.Printf("An email was successfully sent to %s", recipient)
				}
			}
			time.Sleep(time.Hour)
		}
	}
}
