package invalid

import (
	"errors"
	"fmt"
)

const (
	LevelWarning = "warning"
	LevelError   = "error"
)

type Level string

type Result struct {
	Level   Level
	Message string
	Range   *Range
}

func NewFieldMissingError(key string) error {
	return errors.New(fmt.Sprintf("Key [%s] is Expected here", key))
}

func NewResult(level Level, msg string, r *Range) Result {
	return Result{
		Level:   level,
		Message: msg,
		Range:   r,
	}
}
