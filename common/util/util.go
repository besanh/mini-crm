package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
)

func ParseAnyToString(value any) (string, error) {
	ref := reflect.ValueOf(value)
	if ref.Kind() == reflect.String {
		return value.(string), nil
	} else if InArray(ref.Kind(), []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}) {
		return fmt.Sprintf("%d", value), nil
	} else if InArray(ref.Kind(), []reflect.Kind{reflect.Float32, reflect.Float64}) {
		return fmt.Sprintf("%.3f", value), nil
	} else if ref.Kind() == reflect.Bool {
		return fmt.Sprintf("%t", value), nil
	} else if ref.Kind() == reflect.Slice {
		return fmt.Sprintf("%v", value), nil
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MustParseAnyToString(value any) string {
	str, err := ParseAnyToString(value)
	if err != nil {
		return ""
	}
	return str
}

func ParseStringToAny(value string, dest any) error {
	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return err
	}
	return nil
}

func ParseAnyToAny(value any, dest any) (err error) {
	ref := reflect.ValueOf(value)
	var bytes []byte
	if ref.Kind() == reflect.String {
		bytes = []byte(value.(string))
	} else {
		bytes, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func ParseString(value any) string {
	str, ok := value.(string)
	if !ok {
		return str
	}
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Trim(str, "\r\n")
	str = strings.TrimSpace(str)
	return str
}

func InArray[T comparable](element T, slice []T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

var LETTER_RUNES = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

var NUMBER_RUNES = []rune("1234567890")

func GenerateRandomString(n int, letterRunes []rune) string {
	if len(letterRunes) < 1 {
		letterRunes = LETTER_RUNES
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ParseStringToTime(t string, timezone ...string) *time.Time {
	if len(t) == 0 {
		return nil
	}
	c := carbon.Parse(t, timezone...)
	if c.Error != nil {
		return nil
	}
	tPtr := c.StdTime()
	return &tPtr
}

func CheckFromAndToDateValid(from, to time.Time, isAllowZero bool) (isOk bool, err error) {
	if !isAllowZero {
		if from.IsZero() {
			return false, errors.New("is zero")
		}
		if to.IsZero() {
			return false, errors.New("is zero")
		}
	} else if from.After(to) {
		return false, errors.New("from date must be before to date")
	}
	isOk = true
	return
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveEmpty[T string](sliceList []T) []T {
	list := []T{}
	for _, item := range sliceList {
		if item != "" {
			list = append(list, item)
		}
	}
	return list
}

func ParseFloat64(value any) float64 {
	if value == nil {
		return 0
	}
	// convert to string
	str := MustParseAnyToString(value)
	// convert to float
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return result
}

func ParseInt64(value any) int64 {
	if value == nil {
		return 0
	}
	// convert to string
	str := MustParseAnyToString(value)
	// convert to int
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

// func to get end of day
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

// func to parse float64 only with 2 decimal
func ParseFloat64With2Decimal(value float64) float64 {
	return math.Round(value*100) / 100
}

// func to ternary
func Ternary[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
