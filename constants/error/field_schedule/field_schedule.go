package error

import "errors"

var (
	ErrFieldShceduleExist    = errors.New("field schedule already exist")
	ErrFieldScheduleNotFound = errors.New("field schedule not found")
)

var FieldScheduleErrors = []error{
	ErrFieldShceduleExist,
	ErrFieldScheduleNotFound,
}
