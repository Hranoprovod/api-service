package main

import (
	"fmt"
	"strconv"
	"time"
)

func getFloat(value string) float64 {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return num
}

func valToStr(value float64) (string, error) {
	return fmt.Sprintf("%0.2f", value), nil
}

func timeToStr(ts time.Time) (string, error) {
	return ts.Format("01.02.2006 15:04"), nil
}
