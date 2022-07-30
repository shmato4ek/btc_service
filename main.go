package main

import (
	"btc_service/src/btc"
	"btc_service/src/model"
	"btc_service/src/persistance"
	"btc_service/src/sender"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type config struct {
	emailPassword    string
	emailForSending  string
	smtpName         string
	mailText         string
	mailSubject      string
	databaseFileName string
	btcApiName       string
	filePath         string
}

func newConfig() config {
	var newConfig config
	newConfig.btcApiName = goDotEnvVariable("BTC_API_NAME")
	newConfig.emailPassword = goDotEnvVariable("EMAIL_PASSWORD")
	newConfig.databaseFileName = goDotEnvVariable("DATABASE_FILE_NAME")
	newConfig.emailForSending = goDotEnvVariable("EMAIL_FOR_SENDING")
	newConfig.smtpName = goDotEnvVariable("SMTP_NAME")
	newConfig.mailText = goDotEnvVariable("MAIL_TEXT")
	newConfig.mailSubject = goDotEnvVariable("MAIL_SUBJECT")
	newConfig.filePath = goDotEnvVariable("FILEPATH")
	return newConfig
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

/*func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page")
}*/

func getRate(w http.ResponseWriter, r *http.Request) {
	config := newConfig()
	rate, err := btc.GetRate(config.btcApiName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(rate)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	config := newConfig()
	var email string
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fdb, err := persistance.New(config.filePath, config.databaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exists := fdb.Save(model.Email(email), config.databaseFileName)
	if exists {
		http.Error(w, "This email has already subscribed", http.StatusConflict)
	}
}

func sendEmails(w http.ResponseWriter, r *http.Request) {
	config := newConfig()
	fdb, err := persistance.New(config.filePath, config.databaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sender := sender.New(config.emailForSending, config.smtpName, config.emailPassword, config.mailText, config.mailSubject)
	sendErr, rateErr := sender.SendRate(fdb, config.btcApiName)
	if sendErr != nil {
		http.Error(w, sendErr.Error(), http.StatusBadRequest)
	} else if rateErr != nil {
		http.Error(w, rateErr.Error(), http.StatusBadRequest)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/", homeLink)
	router.HandleFunc("/rate", getRate).Methods("GET")
	router.HandleFunc("/subscribe", subscribe).Methods("POST")
	router.HandleFunc("/sendEmails", sendEmails).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
