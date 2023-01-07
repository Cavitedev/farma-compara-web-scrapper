package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

func ParseSpanishNumberStrToNumber(str string) float64 {
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	var incomingNumber string = re.FindString(str)
	var numberAmerican string = strings.Replace(incomingNumber, ",", ".", -1)
	priceNum, err := strconv.ParseFloat(numberAmerican, 64)
	if err != nil {
		log.Println("Error parsing " + str)
	}

	return priceNum
}
