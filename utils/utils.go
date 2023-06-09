package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

func IsNumber(str string) bool {
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	return re.MatchString(str)
}

func NumberRegexString(str string) string {
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	var incomingNumber string = re.FindString(str)
	return incomingNumber
}

func ParseSpanishNumberStrToNumber(str string) float64 {
	var incomingNumber string = NumberRegexString(str)
	var numberAmerican string = strings.Replace(incomingNumber, ",", ".", -1)
	priceNum, err := strconv.ParseFloat(numberAmerican, 64)
	if err != nil {
		log.Println("Error parsing " + str)
	}

	return priceNum
}
