package persistance

import (
	"btc_service/src/model"
	"bufio"
	"errors"
	"log"
	"os"
	//"path/filepath"
)

type fileDatabase struct {
	filePath string
	buffer   []model.Email
}

func New(filepath string) *fileDatabase {
	_, err := os.Stat("database.txt")

	if errors.Is(err, os.ErrNotExist) {
		CreateFile()
	}
	return &fileDatabase{
		filePath: filepath,
		buffer:   readFromFile(filepath),
	}
}

func CreateFile() {
	file, err := os.Create("database.txt")

	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func readFromFile(filepath string) []model.Email {
	var emails []model.Email
	f, err := os.Open("database.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		emails = append(emails, model.Email(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return emails
}

func (fdb *fileDatabase) Save() {

}

func (fdb *fileDatabase) Exists(email model.Email) bool {
	exists := false
	for _, element := range fdb.buffer {
		if element == email {
			exists = true
		}
	}
	return exists
}
