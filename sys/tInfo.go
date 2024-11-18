package sys

import (
	"strconv"
	"time"
)

func IsBefore(dateStr string) bool {
	currTime := time.Now().Unix()
	expDate, _ := time.Parse(time.DateOnly, dateStr)
	expDateStamp := expDate.Unix()

	if currTime < expDateStamp {
		return true
	} else {
		return false
	}
}

// for reading Prophet results
type Date struct {
	Year  int
	Month int
}

func (this *Date) YYYYMM() int {
	return this.Year*100 + this.Month
}

func (this *Date) MMYYYY() (res string) {
	if this.Month >= 10 {
		res = strconv.Itoa(this.Month) + "/" + strconv.Itoa(this.Year)
	} else {
		res = "0" + strconv.Itoa(this.Month) + "/" + strconv.Itoa(this.Year)
	}

	return res
}

func (this *Date) CalendarMth(i int) (res int) {
	num := (this.Month + i) % 12
	if num == 0 {
		res = 12
	} else {
		res = num
	}
	return
}

func (this *Date) CalendarYr(i int) (res int) {
	if this.Month == 12 {
		res = this.Year + (i+11)/12
	} else {
		res = this.Year + (this.Month+i-1)/12
	}
	return
}

func (this *Date) CalendarDate(i int) Date {
	res := Date{this.CalendarYr(i), this.CalendarMth(i)}
	return res
}

func (this *Date) Months(Year, Month int) int {
	res := (Year-this.Year)*12 + Month - this.Month
	return res
}
