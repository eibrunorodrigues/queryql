package queryql

import (
	"errors"
	"fmt"
	"github.com/eibrunorodrigues/parameter-handler/constants"
	"github.com/eibrunorodrigues/parameter-handler/structs"
	"github.com/eibrunorodrigues/parameter-handler/types"
	"github.com/eibrunorodrigues/parameter-handler/utils"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	specialFilterStringDecorator = "__"
	andParameter                 = "and"
	inParameter                  = "in"
	orParameter                  = "or"
	notParameter                 = "nin"
	pageParameter                = "page"
	limitParameter               = "limit"
)

type Handle struct {
	Result                []structs.Filter

	avoidTypingParameters []string
	fieldsToRename        map[string]string
}

//AvoidTypingParameters adds and keeps to a list of fields that shouldn't change type
func (h *Handle) AvoidTypingParameters(parameters []string) {
	h.avoidTypingParameters = parameters
}

//RenameField changes the name of the field if arrives
func (h *Handle) RenameField(original string, newKey string) {
	h.fieldsToRename[original] = newKey
}

//AddList appends to result the url.Values correctly formatted
func (h *Handle) AddList(request url.Values, isSpecialFilter bool) error {
	for key, value := range request {
		if len(value) == 1 {
			err := h.AppendToResult(key, value[0], isSpecialFilter, structs.KeyValuePair{})
			if err != nil {
				return err
			}
		} else {
			operations := h.getListOfOperationsFromRepeatedParams(value)
			group := structs.KeyValuePair{}
			for key, value := range operations {

				if key != orParameter { // OrGroup's will be generated later
					group.Key = strconv.Itoa(len(h.Result))
					group.Value = key
				}

				for _, subParameter := range value {
					err := h.AppendToResult(key, subParameter, isSpecialFilter, group)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

//AppendToResult is a pure and simple way to append something to result. The value still is going to be changed
//to the correct typed though.
func (h *Handle) AppendToResult(key string, value string, isSpecialFilter bool, group structs.KeyValuePair) error {
	if err := h.validateParameter(key, value, isSpecialFilter); err != nil {
		return err
	}

	if group.Value == orParameter || isAMatch(value, constants.ContainsOrOperation) {
		h.appendOrOperation(pageParameter, group, value)
		return nil
	}

	expression, err := h.getValueAndOperationSeparated(value)

	if err != nil {
		return err
	}

	if key == "" {
		return errors.New(fmt.Sprintf("empty key passed: %s : %v", key, expression.Value))
	}

	err = h.appendItem(expression.Key, fmt.Sprintf("%v", expression.Value), key, group, isSpecialFilter)

	if err != nil {
		return err
	}

	return nil
}

func (h *Handle) appendItem(operation string, originalValue string, parameterKey string, group structs.KeyValuePair, isSpecialFilter bool) error {
	for key := range h.fieldsToRename {
		if key == parameterKey {
			parameterKey = key
		}
	}

	var typedValue interface{}
	var err error

	if sort.SearchStrings(h.avoidTypingParameters, parameterKey) > 0 {
		typedValue = originalValue
	} else {
		if typedValue, err = h.typeValue(originalValue); err != nil {
			return err
		}
	}

	fieldIsSpecialFilter := strings.HasPrefix(parameterKey, specialFilterStringDecorator)

	if h.isThisADuplicatedInsert(parameterKey, typedValue, operation, group) {
		return nil
	}

	if fieldIsSpecialFilter {
		parameterKey = strings.Replace(parameterKey, specialFilterStringDecorator, "", -1)
	}

	if !isSpecialFilter {
		isSpecialFilter = fieldIsSpecialFilter
	}

	h.Result = append(h.Result, structs.Filter{
		Field:           parameterKey,
		Operation:       operation,
		Value:           typedValue,
		IsSpecialFilter: isSpecialFilter,
		Group:           group,
	})

	return nil
}

func (h *Handle) appendOrOperation(parameterKey string, group structs.KeyValuePair, originalValue string) {
	groupRegex := regexp.MustCompile(constants.IsBetween("\\[", "\\]")).FindStringSubmatch(originalValue)

	groupNumber := ""

	if len(groupRegex) > 0 {
		groupNumber = groupRegex[1]
	}

	if groupNumber != "" {
		if group.Value != nil {
			group = structs.KeyValuePair{Key: "_" + groupNumber, Value: group.Value}
		} else {
			group = structs.KeyValuePair{Key: "_" + groupNumber, Value: orParameter}
		}
	}

	if group.Value == nil {
		group = h.getGroupOrMakeANewOne(orParameter)
	}
	expression, _ := h.getValueAndOperationSeparated(originalValue)
	_ = h.appendItem(expression.Key, fmt.Sprintf("%v", expression.Value), parameterKey, group, false)
}

func (h *Handle) getGroupOrMakeANewOne(group string) structs.KeyValuePair {
	for _, item := range h.Result {
		if !strings.HasPrefix(item.Group.Key, "_") && item.Group.Value == orParameter {
			return item.Group
		}
	}
	return structs.KeyValuePair{Key: strconv.Itoa(len(h.Result)), Value: group}
}

func (h *Handle) getListOfOperationsFromRepeatedParams(listOfValues []string) map[string][]string {
	groupOfOperations := make(map[string][]string)

	addToGroup := func(key string, value string) {
		groupOfOperations[key] = append(groupOfOperations[key], value)
	}

	for _, value := range listOfValues {
		if isAMatch(value, constants.ContainsOperatorsOrFilters) {
			if isAMatch(value, constants.ContainsOrOperation) {
				addToGroup(orParameter, value)
			} else if isAMatch(value, constants.StartsWithDeny) {
				addToGroup(notParameter, value)
			} else {
				addToGroup(andParameter, value)
			}
		} else {
			addToGroup(inParameter, value)
		}
	}

	return groupOfOperations
}

func (h *Handle) typeValue(parameterValue string) (interface{}, error) {
	if parameterValue == "" {
		return "", nil
	}

	if strings.ToLower(parameterValue) == "nil" {
		return nil, nil
	}

	if strings.ToLower(parameterValue) == "true" || strings.ToLower(parameterValue) == "false" {
		result, err := utils.Convert(parameterValue, types.Bool)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	if isAMatch(parameterValue, constants.IsBetweenDoubleQuotes) {
		result := regexp.MustCompile(constants.IsBetweenDoubleQuotes).FindStringSubmatch(parameterValue)[1]
		return result, nil
	}

	if isAMatch(parameterValue, constants.IsAObjectId) {
		result, err := utils.Convert(parameterValue, types.ObjectID)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	if isAMatch(parameterValue, constants.IsAnInteger) {
		result, err := utils.Convert(parameterValue, types.Int)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	if isAMatch(parameterValue, constants.IsADate) {
		result, err := utils.Convert(parameterValue, types.Time)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	if isAMatch(parameterValue, constants.IsAFloat) {
		result, err := utils.Convert(parameterValue, types.Float64)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	return parameterValue, nil
}

func (h *Handle) isThisADuplicatedInsert(key string, value interface{}, operator string, group structs.KeyValuePair) bool {
	for _, item := range h.Result {
		if item.Field == key && fmt.Sprintf("%v", item.Value) == fmt.Sprintf("%v", value) && item.Operation == operator && item.Group.Value == group.Value {
			return true
		}
	}
	return false
}

func (h *Handle) getValueAndOperationSeparated(expression string) (structs.KeyValuePair, error) {
	if isAMatch(expression, constants.ContainsOperatorsOrFilters) {
		expressionOperator := regexp.MustCompile(constants.Operators).FindString(expression)
		expressionValueGroup := regexp.MustCompile(constants.ValueWithoutOperators).FindStringSubmatch(expression)
		return structs.KeyValuePair{Key: expressionOperator, Value: expressionValueGroup[1]}, nil
	}

	return structs.KeyValuePair{Key: "", Value: expression}, nil
}

func (h *Handle) validateParameter(key string, value string, isSpecialFilter bool) error {
	if key == pageParameter || key == limitParameter && isSpecialFilter {
		if _, err := strconv.Atoi(value); err != nil {
			return errors.New(key + " must be an integer")
		}
	}
	return nil
}

func isAMatch(value string, rule string) bool {
	match, err := regexp.MatchString(rule, value)
	if err != nil || !match {
		return false
	}
	return true
}
