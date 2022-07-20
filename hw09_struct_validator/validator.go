package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrVarType           = errors.New("incorrect var type")
	ErrValueType         = errors.New("incorrect value type for validate")
	ErrNotFoundValidator = errors.New("validator not found")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	str := strings.Builder{}

	for _, oneV := range v {
		str.WriteString(oneV.Err.Error())
	}

	return str.String()
}

type validationErrWrapper struct {
	errs []error
}

func (vew *validationErrWrapper) Error() string {
	return fmt.Sprintf("%v", vew.errs)
}

const (
	validatorNameTag = "validate"
)

func Validate(v interface{}) error {
	var validationErrs ValidationErrors

	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Struct {
		return ErrVarType
	}

	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		errValidate := validateFiled(vv.Field(i), vt.Field(i))

		var vew *validationErrWrapper
		if errValidate != nil {
			if !errors.As(errValidate, &vew) {
				return errValidate
			}

			validationErrs = addErrors(vew, validationErrs, vt.Field(i).Name)
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func addErrors(wrapper *validationErrWrapper, validationErrors ValidationErrors, fieldName string) ValidationErrors {
	for _, err := range wrapper.errs {
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
	}
	return validationErrors
}

func validaByRules(rules []Rule) error {
	var errWrapper validationErrWrapper

	for _, r := range rules {
		err := r.ValidateRule()
		if err != nil {
			errWrapper.errs = append(errWrapper.errs, err)
		}
	}

	return &errWrapper
}

func validateFiled(fv reflect.Value, ft reflect.StructField) error {
	rulesString := ft.Tag.Get(validatorNameTag)
	if rulesString == "" {
		return nil
	}

	var rules []Rule
	var err error

	switch fv.Kind() {
	case reflect.Int:
		rules, err = getIntRules(rulesString, int(fv.Int()))
	case reflect.String:
		rules, err = getStringRules(rulesString, fv.String())
	case reflect.Slice:
		rules, err = getSliceRules(rulesString, fv)
	default:
		return ErrNotFoundValidator
	}

	if err != nil {
		return err
	}

	return validaByRules(rules)
}

func getSliceRules(rs string, fv reflect.Value) ([]Rule, error) {
	var rules []Rule

	switch slice := fv.Interface().(type) {
	case []string:
		for _, str := range slice {
			rulesIteration, err := getStringRules(rs, str)
			if err != nil {
				return rules, err
			}

			rules = append(rules, rulesIteration...)
		}
	case []int:
		for _, number := range slice {
			rulesIteration, err := getIntRules(rs, number)
			if err != nil {
				return rules, err
			}

			rules = append(rules, rulesIteration...)
		}
	default:
		return rules, ErrNotFoundValidator
	}

	return rules, nil
}

func getStringRules(rs string, value string) ([]Rule, error) {
	var rules []Rule

	sliceRuleString := strings.Split(rs, "|")

	if len(sliceRuleString) == 0 {
		return rules, ErrVarType
	}

	for _, osr := range sliceRuleString {
		sosr := strings.Split(osr, ":")
		switch sosr[0] {
		case "len":
			val, err := strconv.Atoi(sosr[1])

			if err != nil {
				return rules, err
			}

			rules = append(rules, &StringLenRule{
				len:        val,
				fieldValue: value,
			})
		case "regexp":
			re, err := regexp.Compile(sosr[1])
			if err != nil {
				return rules, err
			}
			rules = append(rules, &StringRegexpRule{
				regexp:     re,
				fieldValue: value,
			})
		case "in":
			rules = append(rules, &StringInRule{
				in:         strings.Split(sosr[1], ","),
				fieldValue: value,
			})
		default:
			return rules, ErrValueType
		}
	}

	return rules, nil
}

func getIntRules(rs string, value int) ([]Rule, error) {
	rules := make([]Rule, 0)

	sliceRuleString := strings.Split(rs, "|")

	if len(sliceRuleString) == 0 {
		return rules, ErrVarType
	}

	for _, osr := range sliceRuleString {
		sosr := strings.Split(osr, ":")
		switch sosr[0] {
		case "min":
			val, err := strconv.Atoi(sosr[1])
			if err != nil {
				return rules, err
			}
			rules = append(rules, &IntMinRule{
				min:        val,
				fieldValue: value,
			})
		case "max":
			val, err := strconv.Atoi(sosr[1])
			if err != nil {
				return rules, err
			}
			rules = append(rules, &IntMaxRule{
				max:        val,
				fieldValue: value,
			})
		case "in":
			in, err := covertStringToIntSlice(sosr[1])
			if err != nil {
				return rules, err
			}
			rules = append(rules, &IntInRule{
				in:         in,
				fieldValue: value,
			})
		default:
			return rules, ErrValueType
		}
	}

	return rules, nil
}

func covertStringToIntSlice(str string) ([]int, error) {
	strSlice := strings.Split(str, ",")
	intSlice := make([]int, 0)

	for _, val := range strSlice {
		intVar, err := strconv.Atoi(val)
		if err != nil {
			return intSlice, err
		}
		intSlice = append(intSlice, intVar)
	}

	return intSlice, nil
}
