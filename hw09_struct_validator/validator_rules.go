package hw09structvalidator

import (
	"errors"
	"regexp"
)

var (
	ErrIntMin    = errors.New("поле не должно быть меньше указанного значения")
	ErrIntMax    = errors.New("поле не должно быть больше указанного значения")
	ErrIntIn     = errors.New("число должна входить в множество чисел")
	ErrStrRegexp = errors.New("строка должна соответствовать шаблону")
	ErrStrLen    = errors.New("длина строки не соответствует указанному значению")
	ErrStrIn     = errors.New("строка должна входить в множество строк")
)

type Rule interface {
	ValidateRule() error
}

type StringLenRule struct {
	fieldValue string
	len        int
}

func (r *StringLenRule) ValidateRule() error {
	if len(r.fieldValue) == r.len {
		return nil
	}

	return ErrStrLen
}

type StringRegexpRule struct {
	fieldValue string
	regexp     *regexp.Regexp
}

func (r *StringRegexpRule) ValidateRule() error {
	if r.regexp.MatchString(r.fieldValue) {
		return nil
	}

	return ErrStrRegexp
}

type StringInRule struct {
	fieldValue string
	in         []string
}

func (r *StringInRule) ValidateRule() error {
	err := ErrStrIn

	for _, str := range r.in {
		if str == r.fieldValue {
			err = nil

			break
		}
	}

	return err
}

type IntInRule struct {
	fieldValue int
	in         []int
}

func (r *IntInRule) ValidateRule() error {
	err := ErrIntIn

	for _, number := range r.in {
		if number == r.fieldValue {
			err = nil
		}
	}

	return err
}

type IntMinRule struct {
	fieldValue int
	min        int
}

func (r *IntMinRule) ValidateRule() error {
	if r.fieldValue < r.min {
		return ErrIntMin
	}

	return nil
}

type IntMaxRule struct {
	fieldValue int
	max        int
}

func (r *IntMaxRule) ValidateRule() error {
	if r.fieldValue > r.max {
		return ErrIntMax
	}

	return nil
}
