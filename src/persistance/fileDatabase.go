package persistance

import (
	"btc_service/src/model"
	"bufio"
	"errors"
	"log"
	"os"
)

type FileDatabase struct {
	filePath string
	Buffer   []model.Email
}

func New(filepath string, fileName string) (*FileDatabase, error) {
	_, err := os.Stat(fileName)
	var createFileError error
	if errors.Is(err, os.ErrNotExist) {
		createFile(fileName)
	}
	return &FileDatabase{
		filePath: filepath,
		Buffer:   readFromFile(filepath, fileName),
	}, createFileError
}

func createFile(fileName string) error {
	var createFileError error
	file, err := os.Create(fileName)
	if err != nil {
		createFileError = err
	}
	defer file.Close()
	return createFileError
}

func readFromFile(filepath string, fileName string) []model.Email {
	var emails []model.Email
	f, err := os.Open(fileName)

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

func (fdb *FileDatabase) Save(email model.Email, fileName string) bool {
	exists := fdb.Exists(email)
	if !exists {
		fdb.AddNewEmail(email, fileName)
	}
	fdb.Buffer = readFromFile(fdb.filePath, fileName)
	return exists
}

func (fdb *FileDatabase) Exists(email model.Email) bool {
	exists := false
	for _, element := range fdb.Buffer {
		if element == email {
			exists = true
		}
	}
	return exists
}

func (fdb *FileDatabase) AddNewEmail(email model.Email, fileName string) {
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(string(email) + "\n"); err != nil {
		log.Fatal(err)
	}
}
