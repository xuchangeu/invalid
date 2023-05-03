package invalid

import (
	"errors"
	"fmt"
)

const (
	KeyMissing        ResultType = "keyMissing"
	TypeMismatch                 = "typeMismatch"
	StrLengthMismatch            = "strLengthMismatch"
	RegxMismatch                 = "regxMismatch"
)

type ResultType string

type Result struct {
	Type  ResultType
	Error error
	Range *Range
}

func NewKeyMissingError(key string) error {
	return errors.New(fmt.Sprintf("key [%s] is expected here", key))
}

func NewTypeMismatchError(key, ty string) error {
	return errors.New(fmt.Sprintf("type for [%s] must be [%s]", key, ty))
}

func NewStrLengthError1(key string, len int) error {
	return errors.New(fmt.Sprintf("length of value in [%s] must < %d", key, len))
}

func NewStrLengthError2(key string, len int) error {
	return errors.New(fmt.Sprintf("length of value in [%s] must > %d", key, len))
}

func NewRegxError(key, regx string) error {
	return errors.New(fmt.Sprintf("value for [%s] must match regexp : %s", key, regx))
}

func NewKeyNameError(key, regx string) error {
	return errors.New(fmt.Sprintf("key name for [%s] must match regexp ： %s", key, regx))
}

func NewResult(t ResultType, err error, r *Range) Result {
	return Result{
		Type:  t,
		Error: err,
		Range: r,
	}
}
