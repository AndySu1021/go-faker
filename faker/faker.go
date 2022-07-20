package faker

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandSlice(slice []interface{}) interface{} {
	return slice[rand.Intn(len(slice))]
}

func RandStr(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandDateTime(start, end string) string {
	layout := "2006-01-02 15:04:05"
	startTime, _ := time.Parse(layout, start)
	startStamp := startTime.UTC().Unix()
	endTime, _ := time.Parse(layout, end)
	endStamp := endTime.UTC().Unix()
	return time.Unix(RandInt64(startStamp, endStamp), 0).Format(layout)
}

func RandDate(start, end string) string {
	layout := "2006-01-02"
	startTime, _ := time.Parse(layout, start)
	startStamp := startTime.UTC().Unix()
	endTime, _ := time.Parse(layout, end)
	endStamp := endTime.UTC().Unix()
	return time.Unix(RandInt64(startStamp, endStamp), 0).Format(layout)
}
