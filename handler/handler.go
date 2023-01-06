package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/nilemarezz/flexi-chat-v2/service"
)

type ExpensesHandler struct {
	linebot        *linebot.Client
	expenseService service.ExpensesService
	horoszService  *service.HoroszService
}

func NewExpensesHandler(linebot *linebot.Client, service service.ExpensesService, horosz *service.HoroszService) ExpensesHandler {
	return ExpensesHandler{linebot: linebot, expenseService: service, horoszService: horosz}
}

func (e ExpensesHandler) MessageHandler(event *linebot.Event, message string) {
	fmt.Println("message in : ", message)
	s := strings.Split(message, " ")
	switch s[0] {
	case "สรุป", "summary":
		e.summary(message, event)
	case "รายการ", "list":
		e.list(message, event)
	case "hi", "สวัสดี", "คำสั่ง", "Hi":
		e.hello(event)
	case "ดูดวง":
		e.horosz(event)
	default:
		e.save(message, event)
	}

}

func (e ExpensesHandler) save(message string, event *linebot.Event) {
	res, err := e.expenseService.Save(message, event.Source.UserID)

	if err != nil {
		e.sendMessageWithExample(event.ReplyToken, err.Error(), "คำสั่ง -> [รับ,จ่าย]/จำนวน/รายการ\nตย. -> จ่าย/200/ค่ารถ")
		return

	}

	e.sendMessage(event.ReplyToken, res.Message)
}

func (e ExpensesHandler) list(message string, event *linebot.Event) {

	res, err := e.expenseService.List(message, event.Source.UserID)

	if err != nil {
		e.sendMessage(event.ReplyToken, err.Error())
		// fmt.Println(err)
		return
	}
	// res.Message
	lineMessages := []linebot.SendingMessage{}
	for _, v := range res.Message {
		lineMessages = append(lineMessages, linebot.NewImageMessage(v, v))
	}
	// e.sendMessage(event.ReplyToken, res.)
	if _, err := e.linebot.ReplyMessage(event.ReplyToken, lineMessages...).Do(); err != nil {
		log.Print(err)
	}
}

func (e ExpensesHandler) summary(message string, event *linebot.Event) {
	fmt.Println("summary , " + event.Source.UserID)
	res, err := e.expenseService.Summary(message, event.Source.UserID)

	if err != nil {
		e.sendMessageWithExample(event.ReplyToken, err.Error(), "คำสั่ง -> สรุป เดือน/ปี\nตย. -> สรุป 2/2023")
		return
	}
	e.sendMessage(event.ReplyToken, res.Message)
}

func (e ExpensesHandler) horosz(event *linebot.Event) {
	fmt.Println("horosz , " + event.Source.UserID)
	res, err := e.horoszService.Run()

	if err != nil {
		e.sendMessageWithExample(event.ReplyToken, err.Error(), "เกิดข้อผิดพลาด")
		return
	}

	fmt.Println(res)

	lineMessages := []linebot.SendingMessage{}

	for i := 0; i < 4; i++ {
		lineMessages = append(lineMessages, linebot.NewTextMessage(res[i]))
	}

	lineMessages = append(lineMessages, linebot.NewImageMessage(res[4], res[4]))

	if _, err := e.linebot.ReplyMessage(event.ReplyToken, lineMessages...).Do(); err != nil {
		log.Print(err)
	}
}

func (e ExpensesHandler) hello(event *linebot.Event) {
	s1 := "รายการคำสั่ง"
	s2 := "เพิ่มรายการ รับ/จ่าย \n- [รับ/จ่าย]/จำนวน/ประเภท \n- จ่าย/100/ขนม"
	s3 := "สรุป \n- สรุป \n- สรุป 2/2023"
	s4 := "รายการ \n- รายการ \n- รายการ 2/2023"

	s := []string{s1, s2, s3, s4}

	lineMessages := []linebot.SendingMessage{}
	for _, v := range s {
		lineMessages = append(lineMessages, linebot.NewTextMessage(v))
	}
	if _, err := e.linebot.ReplyMessage(event.ReplyToken, lineMessages...).Do(); err != nil {
		log.Print(err)
	}
}

func (e ExpensesHandler) sendMessage(replyToken string, message string) {
	fmt.Println("message response" + message)
	if _, err := e.linebot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		log.Print(err)
	}
}

func (e ExpensesHandler) sendMessageWithExample(replyToken string, message string, example string) {
	if _, err := e.linebot.ReplyMessage(replyToken,
		linebot.NewTextMessage(message),
		linebot.NewTextMessage(example)).Do(); err != nil {
		log.Print(err)
	}
}
