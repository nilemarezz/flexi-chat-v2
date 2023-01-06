package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

type visionService struct {
	service *vision.ImageAnnotatorClient
}

func NewVisionService() (*visionService, error) {
	conf := &jwt.Config{
		Email:        os.Getenv("EMAIL"),
		PrivateKey:   []byte(os.Getenv("PRIVATE_KEY")),
		PrivateKeyID: os.Getenv("PRIVATE_ID"),
		TokenURL:     "https://oauth2.googleapis.com/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/cloud-vision",
		},
	}

	b, err := json.Marshal(conf)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	// fmt.Println(string(b))

	// client := conf.Client(context.Background())
	// // conn, err := (context.Background(), option.WithHTTPClient(client))
	// grpc.
	// if err != nil {
	// 	return nil, err
	// }

	v, err := vision.NewImageAnnotatorClient(context.Background(), option.WithCredentialsJSON(b))

	if err != nil {
		return nil, err
	}
	return &visionService{v}, nil
}

func (s visionService) Run(image string) (string, error) {
	ctx := context.Background()

	client := vision.NewImageFromURI(image)

	res := ""

	annotations, err := s.service.DetectTexts(ctx, client, nil, 10)
	if err != nil {
		return "", err
	}

	if len(annotations) == 0 {
		fmt.Println("No text found.")
	} else {
		fmt.Println("Text:")
		for _, annotation := range annotations {
			res += annotation.Description
		}
	}

	return res, nil
}
