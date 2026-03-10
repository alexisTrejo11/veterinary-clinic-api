package shared

import (
	"errors"
	"fmt"
)

var ErrEntityNotFound = errors.New("entity not supported")

type IntegerID interface {
	Value() uint
	Equals(number uint) bool
	String() string
	IsZero() bool
}

type BaseID struct {
	Value uint
}

func (id BaseID) Equals(number uint) bool {
	return id.Value == number
}

func (id BaseID) String() string {
	return fmt.Sprintf("%d", id.Value)
}

func (id BaseID) IsZero() bool {
	return id.Value == 0
}

func (id BaseID) Int32() int32 {
	return int32(id.Value)
}
