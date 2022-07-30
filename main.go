package main

import (
	"btc_service/src/btc"
	"btc_service/src/model"
	"btc_service/src/persistance"
	"btc_service/src/sender"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

const (
	filepath = "C:/Users/gangt/OneDrive/Documents/GitHub/btc_service"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getRate(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(btc.GetRate(goDotEnvVariable("BTC_API_NAME")))
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	var email string
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fdb := persistance.New(filepath, goDotEnvVariable("DATABASE_FILE_NAME"))
	fdb.Save(model.Email(email), goDotEnvVariable("DATABASE_FILE_NAME"))
}

func sendEmails(w http.ResponseWriter, r *http.Request) {
	fdb := persistance.New(filepath, goDotEnvVariable("DATABASE_FILE_NAME"))
	sender := sender.New(goDotEnvVariable("EMAIL_FOR_SENDING"), goDotEnvVariable("SMTP_NAME"), goDotEnvVariable("EMAIL_PASSWORD"), goDotEnvVariable("MAIL_TEXT"), goDotEnvVariable("MAIL_SUBJECT"))
	sender.SendRate(fdb, goDotEnvVariable("BTC_API_NAME"))
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/rate", getRate).Methods("GET")
	router.HandleFunc("/subscribe", subscribe).Methods("POST")
	router.HandleFunc("/sendEmails", sendEmails).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
