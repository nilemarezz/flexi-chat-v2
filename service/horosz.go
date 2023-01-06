package service

import (
	"fmt"
	"strings"

	"github.com/irfansofyana/go-aztro-api-wrapper/aztro"
	"golang.org/x/image/colornames"
)

type HoroszService struct {
}

func NewHoroszService() *HoroszService {
	return &HoroszService{}
}

func (HoroszService) Run() ([]string, error) {
	aztroClient, err := aztro.NewAztroClient()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	aztroParam := aztro.NewAztroRequestParam(aztro.Sign(aztro.Aries), aztro.WithDay(aztro.Today))
	todayHoroscope, aztroErr := aztroClient.GetHoroscope(aztroParam)
	if aztroErr != nil {
		fmt.Println(aztroErr)
		return nil, err
	}
	fmt.Println(todayHoroscope)

	color := convertWebToHex(strings.ReplaceAll(todayHoroscope.Color, " ", ""))
	color = strings.ReplaceAll(color, "#", "")

	return []string{
		todayHoroscope.Description,
		"Lucky Number : " + todayHoroscope.LuckyNumber,
		"Lucky Time : " + todayHoroscope.LuckyTime,
		"Lucky Color : " + todayHoroscope.Color,
		"https://singlecolorimage.com/get/" + color + "/400x100.png",
	}, nil

}

func convertWebToHex(webcolorname string) (hexcolor string) {
	c, ok := colornames.Map[strings.ToLower(webcolorname)]
	if !ok {
		// Unknown name
		return ""
	}
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
