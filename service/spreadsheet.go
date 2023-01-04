package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type spreadSheetService struct {
	service *sheets.Service
	sheetID string
}

func NewSpreadSheetService(sheetID string) (*spreadSheetService, error) {
	// content, err := ioutil.ReadFile("../credential/key.txt")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(content))

	conf := &jwt.Config{
		Email:        os.Getenv("EMAIL"),
		PrivateKey:   []byte(os.Getenv("PRIVATE_KEY")),
		PrivateKeyID: os.Getenv("PRIVATE_ID"),
		TokenURL:     "https://oauth2.googleapis.com/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	client := conf.Client(context.Background())

	// Create a service object for Google sheets
	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &spreadSheetService{service: srv, sheetID: sheetID}, nil
}

func (s spreadSheetService) GetByRange(ranges string) ([][]interface{}, error) {
	resp, err := s.service.Spreadsheets.Values.Get(s.sheetID, ranges).Do()
	if err != nil || resp.HTTPStatusCode != 200 {
		fmt.Println("error 48", err)
		return nil, err.(*googleapi.Error)
	}

	if len(resp.Values) == 0 {
		return nil, errors.New("no record found")
	} else {
		return resp.Values, nil
	}
}

func (s spreadSheetService) InsertRow(sheetName string, records [][]interface{}) error {
	// records := [][]interface{}{{"a1", "b1", "c1"}} // This is a sample value.

	valueInputOption := "USER_ENTERED"
	insertDataOption := "INSERT_ROWS"
	rb := &sheets.ValueRange{
		Values: records,
	}
	resp, err := s.service.Spreadsheets.Values.Append(s.sheetID, sheetName, rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(context.Background()).Do()
	if err != nil || resp.HTTPStatusCode != 200 {
		log.Println(err)
		return err.(*googleapi.Error)
	}
	return nil
}
