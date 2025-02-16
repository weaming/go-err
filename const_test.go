package goerr

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func TestError(t *testing.T) {
	keyErr := NewError(errors.New("key error"))
	log.Printf("FromError: err %s, fmt %s", keyErr.Debug(), keyErr.Error())

	log.Printf("InternalError: err %s, fmt %s", InternalError.Debug(), InternalError.Error())
	log.Printf("IOError: err %s, fmt %s", IOError.Debug(), IOError.Error())
	log.Printf("DBError: err %s, fmt %s", DBError.Debug(), DBError.Error())

	equal(t, "Database error", fmt.Sprintf("%s", DBError), "format %s failed")
	equal(t, "Database error", fmt.Sprintf("%v", DBError), "format %v failed")
	equal(t, "Database error", fmt.Sprintf("%+v", DBError), "format %+v failed")

	assert(t, IOError.Is(InternalError), "IOError.Is(InternalError) should be true")
	assert(t, DBError.Is(IOError), "DBError.Is(IOError) should be true")
	assert(t, DBError.Is(InternalError), "DBError.Is(InternalError) should be true")

	assert(t, errors.Is(IOError, InternalError), "IOError should be an InternalError")
	assert(t, errors.Is(DBError, IOError), "DBError should be an IOError")
	assert(t, errors.Is(DBError, InternalError), "DBError should be an InternalError")

	var customDbErr = DBError.Extend("dial tcp fail")
	log.Printf("customDbErr: err %v, fmt %v", customDbErr.Debug(), customDbErr.Error())
	assert(t, errors.Is(customDbErr, InternalError), "customDbErr should be an InternalError")
	assert(t, errors.Is(customDbErr, IOError), "customDbErr should be an IOError")

	var customDbErr2 = customDbErr.SetMsg("internal error")
	assert(t, customDbErr2.Error() == "internal error", "customDbErr2 should be an InternalError at the view of outside")
	log.Printf("customDbErr2: err %v, fmt %v", customDbErr2.Debug(), customDbErr2.Error())
	assert(t, errors.Is(customDbErr2, InternalError), "customDbErr2 should be an InternalError")
	assert(t, errors.Is(customDbErr2, IOError), "customDbErr2 should be an IOError")
}

func equal(t *testing.T, expect, actual any, msg string) {
	if expect != actual {
		t.Fatalf("%s: expect '%v', got '%v'", msg, expect, actual)
	}
}

func assert(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Fatalf("Assertion failed: %s", msg)
	}
}

/*
=== RUN   TestError
2024/07/10 14:31:28 InternalError: err Internal error, fmt Internal error

2024/07/10 14:31:28 IOError: err I/O error
Internal error, fmt I/O error

2024/07/10 14:31:28 DBError: err Database error
I/O error
Internal error, fmt Database error

2024/07/10 14:31:28 customDbErr: err Database error: dial tcp fail
Database error
I/O error
Internal error, fmt Database error: dial tcp fail

2024/07/10 14:31:28 customDbErr2: err Database error: dial tcp fail
Database error
I/O error
Internal error, fmt internal error
*/
