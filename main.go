package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/nilemarezz/flexi-chat-v2/handler"
	"github.com/nilemarezz/flexi-chat-v2/service"
)

func init() {
	initTimeZone()
}

func main() {
	godotenv.Load(".env")
	secret := os.Getenv("CHANNEL_SECRET")
	token := os.Getenv("CHANNEL_TOKEN")
	bot, err := linebot.New(
		// os.Getenv("CHANNEL_SECRET"),
		// os.Getenv("CHANNEL_TOKEN"),
		secret,
		token,
	)
	if err != nil {
		log.Fatal(err)
	}

	// create handler

	r := mux.NewRouter()
	spreadSheetService, err := service.NewSpreadSheetService(os.Getenv("SPREADSHEET_ID"))
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	horoszService := service.NewHoroszService()
	expenseService := service.NewExpensesService(spreadSheetService)
	expenseHandler := handler.NewExpensesHandler(bot, expenseService, horoszService)

	// Setup HTTP Server for receiving requests from LINE platform
	r.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					expenseHandler.MessageHandler(event, message.Text)
				case *linebot.ImageMessage:
					fmt.Println(message.OriginalContentURL)
				}
			}
		}
	})

	r.HandleFunc("/img/{user}/{date}/{file}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		user := vars["user"]
		date := vars["date"]
		file := vars["file"]

		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("./img/%v/%v/%v", user, date, file))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)

	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("App start at port " + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic("")
	}
	time.Local = ict
}
