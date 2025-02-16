package goerr

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var _ ErrorInterface = &Error{}

type ErrorInterface interface {
	// compatible,
	// and for output without change existing code
	error

	// keep more context values
	Set(key string, value any) *Error
	Get(key string) (value any, ok bool)
	GetAll() map[string]any

	// cooperate with errors.Is
	Is(error) bool

	// wrap a new subtype error
	Wrap(err error) *Error
	// create a new variant
	Extend(format string, a ...any) *Error
	// update the message
	SetMsg(format string, a ...any) *Error

	// for log
	Debug() string
}

type Error struct {
	err    error
	errMsg string
	values sync.Map
}

func NewError(err error) (err2 *Error) {
	err2 = &Error{err: err} // Keep the original error

	switch x := err.(type) {
	case interface{ Unwrap() error }:
		err = x.Unwrap()
		if err != nil {
			err2.errMsg = err.Error()
		}
	case interface{ Unwrap() []error }:
		errs := x.Unwrap()
		if len(errs) > 0 {
			// less output, details see err2.err
			err2.errMsg = errs[0].Error()
		}
	default:
		err2.errMsg = err.Error()
	}
	return
}

func NewErrorFromString(msg string) *Error {
	return &Error{errMsg: msg, err: errors.New(msg)}
}

func (e *Error) Set(key string, value any) *Error {
	e.values.Store(key, value)
	return e
}

func (e *Error) Get(key string) (value any, ok bool) {
	value, ok = e.values.Load(key)
	return
}

func (e *Error) GetAll() map[string]any {
	m := make(map[string]any)
	e.values.Range(func(key, value any) bool {
		m[key.(string)] = value
		return true
	})
	return m
}

func (e *Error) Error() string {
	if e.errMsg != "" {
		return e.errMsg
	}
	return e.Error()
}

func (e *Error) Debug() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.errMsg
}

func (e *Error) Is(target error) bool {
	if e2, ok := target.(*Error); ok {
		if strings.HasPrefix(e.errMsg, e2.errMsg) {
			return true
		}
	}

	// use the original error to compare
	return errors.Is(e.err, target)
}

func (e *Error) Wrap(err error) *Error {
	return NewError(errors.Join(err, e))
}

func (e *Error) Extend(format string, a ...any) *Error {
	fmtStr := fmt.Sprintf(format, a...)
	suffix := ": " + fmtStr

	var err error
	if e.err != nil {
		switch x := e.err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
			if err != nil {
				err = errors.New(err.Error() + suffix)
			}
		case interface{ Unwrap() []error }:
			errs := x.Unwrap()
			if len(errs) > 0 {
				new0 := errors.New(errs[0].Error() + suffix)
				errsNew := make([]error, 1, len(errs))
				errsNew[0] = new0
				errs = append(errsNew, errs...)
			} else {
				errs = []error{errors.New(fmtStr)}
			}
			err = errors.Join(errs...)
		}
	} else {
		err = errors.New(e.errMsg + suffix)
	}
	return NewError(err)
}

func (e *Error) SetMsg(format string, a ...any) *Error {
	e.errMsg = fmt.Sprintf(format, a...)
	return e
}
