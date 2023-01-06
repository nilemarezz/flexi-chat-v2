package service

import (
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/karmdip-mi/go-fitz"
	"github.com/nilemarezz/flexi-chat-v2/dto"
	"google.golang.org/api/googleapi"
)
type ExpensesService struct {
	spreadsheet *spreadSheetService
}

func NewExpensesService(spreadsheet *spreadSheetService) ExpensesService {
	return ExpensesService{spreadsheet: spreadsheet}
}

func (service ExpensesService) Save(message string, user string) (*dto.SaveResponse, error) {
	s := strings.Split(message, "/")

	// check must have 3 length
	if len(s) < 3 {
		fmt.Println("command less than 3")
		return nil, errors.New("คำสั่งไม่ครบคับ 😭")
	}

	// check index 0 must be รับ or จ่าย
	if s[0] != "รับ" && s[0] != "จ่าย" {
		fmt.Println("no รับ || จ่าย")
		return nil, errors.New("ค่าแรกเป็น รับ,จ่าย คับ 😭")
	}

	// check index 1 must be รับ or จ่าย
	if _, err := strconv.ParseFloat(s[1], 64); err != nil {
		fmt.Println("no numuric")
		return nil, errors.New("ค่าสองเป็นตัวเลขงับ 😭")
	}

	date := time.Now().Local()

	// get sheet name by date now (1_2023)
	sheetName := fmt.Sprintf("%d_%d", date.Month(), date.Year())

	// create record to insert
	d := date.Format("2006-01-02") + " - " + date.Format("15:04:05")
	records := [][]interface{}{{
		d,
		user,
		message,
		s[0],
		s[1],
		s[2],
	}}

	fmt.Println(records)
	// insert row at sheet
	err := service.spreadsheet.InsertRow(sheetName, records)
	if err != nil {
		fmt.Println("insert error")
		return nil, errors.New("ระบบมีปัญหา ขอเช็กแปปป 😭")
	}
	// response line chat
	res := fmt.Sprintf("เรียบร้อยงับ 😊\n%s - %s", date.Format("15:04:05"), message)
	return &dto.SaveResponse{Message: res}, nil
}

func (service ExpensesService) Summary(message, user string) (*dto.SummaryResponse, error) {
	var sheetName string
	date := time.Now()

	s := strings.Split(message, " ")

	// สรุป
	if len(s) == 1 {
		sheetName = fmt.Sprintf("%d_%d", date.Month(), date.Year())
	}

	// สรุป 1/2022
	if len(s) >= 2 {
		inputDate := s[1]
		ids := strings.Split(string(inputDate), "/")
		if len(ids) != 2 {
			return nil, errors.New("รูปแบบ เดือน/ปี ไม่ถูกงับ 😭")
		}
		sheetName = fmt.Sprintf("%v_%v", ids[0], ids[1])
	}

	totalIncome := 0.0
	totalExpense := 0.0

	ranges := sheetName + "!A2:F"

	fmt.Println("sheetName" + ranges)

	records, err := service.spreadsheet.GetByRange(ranges)
	if err != nil {
		if err.Error() == "no record found" {
			return nil, errors.New("ยังไม่มีข้อมูลในระบบงับ 😭")
		}
		if err.(*googleapi.Error).Code == 400 {
			return nil, errors.New("รูปแบบ เดือน/ปี ไม่ถูกงับ 😭")
		}
	}

	for _, row := range records {
		if row[1] == user {

			amount, _ := strconv.ParseFloat(fmt.Sprint(row[4]), 64)

			if row[3] == "รับ" {
				totalIncome += amount
			} else if row[3] == "จ่าย" {
				totalExpense += amount
			}
		}
	}

	res := fmt.Sprintf(" สรุปเดือน %v %v \n--------------------\nรายรับ   : %v \nรายจ่าย : %v \nยอดสุทธิ  : %v",
		date.Month().String()[:3], date.Year(),
		totalIncome,
		totalExpense,
		totalIncome-totalExpense)

	return &dto.SummaryResponse{Message: res}, nil
}

func (service ExpensesService) List(message, user string) (*dto.ListResponse, error) {
	var sheetName string
	date := time.Now()

	s := strings.Split(message, " ")

	// สรุป
	if len(s) == 1 {
		sheetName = fmt.Sprintf("%d_%d", date.Month(), date.Year())
	}

	// สรุป 1/2022
	if len(s) >= 2 {
		inputDate := s[1]
		ids := strings.Split(string(inputDate), "/")
		sheetName = fmt.Sprintf("%v_%v", ids[0], ids[1])
	}

	ranges := sheetName + "!A2:F"
	records, _ := service.spreadsheet.GetByRange(ranges)

	contents := [][]string{}

	for _, v := range records {
		if v[1] == user {
			data := []string{fmt.Sprint(v[0]), fmt.Sprint(v[3]), fmt.Sprint(v[4]), fmt.Sprint(v[5])}
			contents = append(contents, data)
		}
	}
	newpath := filepath.Join(".", "temp")
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		return nil, errors.New("เกิดข้อผิดพลาดงับ 😭")
	}
	filename := fmt.Sprintf("%v/%v.pdf", "temp", user)
	err = createTable(contents, filename)

	if err != nil {
		return nil, errors.New("เกิดข้อผิดพลาดงับ 😭")
	}

	err = createImage(sheetName)
	if err != nil {
		return nil, errors.New("เกิดข้อผิดพลาดงับ 😭")
	}
	var files []string

	root := "./img/" + user + "/" + sheetName
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" {
			fmt.Println(path)

			files = append(files, os.Getenv("ADDRESS")+"/"+path)
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("เกิดข้อผิดพลาดงับ 😭")
	}
	fmt.Println(files)
	return &dto.ListResponse{Message: files}, nil
}

func createTable(contents [][]string, fileName string) error {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)
	m.AddUTF8Font("CustomArial", consts.Normal, "assets/arial.ttf")
	m.AddUTF8Font("CustomArial", consts.Italic, "assets/arial.ttf")
	m.AddUTF8Font("CustomArial", consts.Bold, "assets/arial.ttf")
	m.AddUTF8Font("CustomArial", consts.BoldItalic, "assets/arial.ttf")

	buildList(m, contents)
	err := m.OutputFileAndClose(fileName)
	if err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
		return err
	}

	fmt.Println("PDF saved successfully")
	return nil
}

func buildList(m pdf.Maroto, contents [][]string) {
	date := time.Now()
	tableHeadings := []string{"Date", "Type", "Amount", "Note"}
	lightPurpleColor := getLightPurpleColor()

	m.SetBackgroundColor(getTealColor())
	m.SetDefaultFontFamily("CustomArial")

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("รายการ รายรับ - รายจ่าย %v %v", date.Month().String()[:3], date.Year()), props.Text{
				Top:    2,
				Size:   16,
				Color:  color.NewWhite(),
				Family: "CustomArial",
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.TableList(tableHeadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      14,
			GridSizes: []uint{4, 1, 2, 5},
		},
		ContentProp: props.TableListContent{
			Size:      12,
			GridSizes: []uint{4, 1, 2, 5},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
}

func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}

func getTealColor() color.Color {
	return color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	}
}

func createImage(sheetName string) error {
	var files []string

	root := "temp/"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".pdf" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, file := range files {
		doc, err := fitz.New(file)
		if err != nil {
			return err
		}
		folder := strings.TrimSuffix(path.Base(file), filepath.Ext(path.Base(file))) + "/" + sheetName
		fmt.Println(folder)
		// Extract pages as images
		for n := 0; n < doc.NumPage(); n++ {
			img, err := doc.Image(n)
			if err != nil {
				return err
			}
			err = os.MkdirAll("./img/"+folder, 0755)
			if err != nil {
				return err
			}

			f, err := os.Create(filepath.Join("./img/"+folder+"/", fmt.Sprintf("%v.jpg", n)))
			if err != nil {
				return err
			}

			err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			if err != nil {
				return err
			}

			f.Close()

		}
	}
	return nil
}
