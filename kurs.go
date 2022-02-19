package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

func main() {

	type ValCurs struct {
		XMLName xml.Name `xml:"ValCurs"`
		Text    string   `xml:",chardata"`
		Date    string   `xml:"Date,attr"`
		Name    string   `xml:"name,attr"`
		Valute  []struct {
			Text     string `xml:",chardata"`
			ID       string `xml:"ID,attr"`
			NumCode  string `xml:"NumCode"`
			CharCode string `xml:"CharCode"`
			Nominal  string `xml:"Nominal"`
			Name     string `xml:"Name"`
			Value    string `xml:"Value"`
			ValFloat float64
		} `xml:"Valute"`
	}
	valCur := make([]ValCurs, 90)
	t := time.Now()
	var maxVal float64
	var datMax string
	var minVal float64
	var datMin string
	var sumVal float64
	var srVal float64

	for i := 0; i <= 89; i++ {
		t1 := t.Format("02/01/2006")
		s := "http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req="
		s1 := s + t1
		t = t.Add(-24 * time.Hour)

		resp, err := http.Get(s1)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		reader := bytes.NewReader(body)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		err = decoder.Decode(&valCur[i])
		for j := 0; j <= 33; j++ {
			valCur[i].Valute[j].Value = strings.Replace(valCur[i].Valute[j].Value, ",", ".", -1)
			valCur[i].Valute[j].ValFloat, err = strconv.ParseFloat(valCur[i].Valute[j].Value, 64)
			if err != nil {
				log.Fatalln(err)
			}

		}
	}
	for j := 0; j <= 33; j++ {
		maxVal = 0

		for i := 0; i <= 89; i++ {
			if maxVal < valCur[i].Valute[j].ValFloat {
				maxVal = valCur[i].Valute[j].ValFloat
				datMax = valCur[i].Date
			}
		}
		fmt.Println("Значение максимального курса валюты:", maxVal, "Название валюты:", valCur[0].Valute[j].Name, "Дата максимального значения:", datMax)

	}
	fmt.Println()
	for j := 0; j <= 33; j++ {
		minVal = 999999999

		for i := 0; i <= 89; i++ {
			if minVal > valCur[i].Valute[j].ValFloat {
				minVal = valCur[i].Valute[j].ValFloat
				datMin = valCur[i].Date
			}
		}
		fmt.Println("Значение минимального курса валюты:", minVal, "Название валюты:", valCur[0].Valute[j].Name, "Дата минимального значения:", datMin)

	}
	fmt.Println()

	for j := 0; j <= 33; j++ {
		sumVal = 0
		for i := 0; i <= 89; i++ {

			sumVal += valCur[i].Valute[j].ValFloat
		}
		srVal = math.Round((sumVal/90)*10000) / 10000
		fmt.Println("Среднее значение курса рубля за весь период:", srVal, "Верно для валюты:", valCur[0].Valute[j].Name)

	}

}
