package utils

import (
	"fmt"
	"github.com/yaa110/go-persian-calendar"
	"strconv"
	"strings"
	"time"
)

func GetDateStr(t time.Time) string {
	p := ptime.New(t)
	return fmt.Sprintf("%04d/%02d/%02d", p.Year(), p.Month(), p.Day())
}

func ParsePersianDate(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, "/")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("فرمت اشتباه تاریخ: %s", dateStr)
	}

	year, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	day, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return time.Time{}, fmt.Errorf("مقدار تاریخ نامعتبر است")
	}

	p := ptime.Date(year, ptime.Month(time.Month(month)), day, 0, 0, 0, 0, time.Local)
	return p.Time(), nil
}
