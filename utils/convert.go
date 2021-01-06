package utils

import (
	"errors"
	"fmt"
	"github.com/eibrunorodrigues/gitql/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

func Convert(value string, newType types.Types) (interface{}, error) {
	switch newType {
	case types.String:
		return fmt.Sprintf("%v", value), nil
	case types.Int:
		newValue, err := strconv.Atoi(value)
		if err != nil {
			return "", err
		}
		return newValue, nil
	case types.Float64:
		newValue, err := strconv.ParseFloat(strings.Replace(value, ",", ".", -1), 64)
		if err != nil {
			return "", err
		}
		return newValue, nil
	case types.Float32:
		newValue, err := strconv.ParseFloat(strings.Replace(value, ",", ".", -1), 32)
		if err != nil {
			return "", err
		}
		return newValue, nil
	case types.Bool:
		newValue, err := strconv.ParseBool(value)
		if err != nil {
			return "", err
		}
		return newValue, nil
	case types.ObjectID:
		newValue, err := primitive.ObjectIDFromHex(value)
		if err != nil {
			return "", err
		}
		return newValue, nil
	case types.Time:
		newValue, err := ParseTime(value)
		if err != nil {
			return "", err
		}
		return newValue, nil
	default:
		return "", errors.New("could not convert value")
	}
}
