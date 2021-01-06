package utils

import (
	"strings"
	"time"
)

func ParseTime(value string) (time.Time, error) {
	if value == "00000000" {
		return time.Time{}, nil
	}

	if len(value) == 8 {
		newValue, err := time.Parse("20060102", value)
		if err != nil {
			newValue, err := time.Parse("02012006", value)
			if err != nil {
				return time.Time{}, err
			}
			return newValue, nil
		}
		return newValue, nil
	}

	if !strings.Contains(value, "T") {
		newValue, err := time.Parse("2006-01-02 15:04:05", value)
		if err != nil {
			return time.Time{}, nil
		}
		return newValue, nil
	}

	newValue, err := time.Parse("2006-01-02T15:04", value)
	if err != nil {
		newValue, err := time.Parse("2006-01-02T15:04:05", value)
		if err != nil {
			return time.Time{}, err
		}
		return newValue, nil
	}
	return newValue, nil
}
