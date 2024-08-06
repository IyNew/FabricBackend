package utils

import (
	"strconv"
	"time"
)

func ConvertToRFC3339(unixtime string) string {
	unixtimeInt, _ := strconv.ParseInt(unixtime, 10, 64)
	t := time.Unix(unixtimeInt, 0)
	return t.Format(time.RFC3339)
}

func ConvertToUnixTime(datetime string) string {
	t, _ := time.Parse(time.RFC3339, datetime)
	return strconv.FormatInt(t.Unix(), 10)
}
