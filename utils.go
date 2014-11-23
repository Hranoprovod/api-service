package apiservice

import (
	"fmt"
	"strconv"
	"time"
)

const (
	PRECISION = 1000
)

func getFloat(value string) float32 {
	num, err := strconv.ParseFloat(value, 32)
	if err != nil {
		panic(err)
	}
	return float32(num)
}

func floatToInt(value float32) int {
	return int(value * PRECISION)
}

func valToStr(value int) (string, error) {
	return fmt.Sprintf("%0.2f", float32(value)/float32(PRECISION)), nil
}

func timeToStr(ts time.Time) (string, error) {
	return ts.Format("01.02.2006 15.04"), nil
}
