package main

import (
	"errors"
	"log"

	goerr "github.com/weaming/go-err"
)

func DoSomething() (any, error) {
	return nil, goerr.DBError.Extend("connect fail, user %d", 123)
}

func DoSomething2() (any, error) {
	return nil, goerr.IOError.Wrap(errors.New("redis err")).Set("user", 123).SetMsg("internal error")
}

func main() {
	_, err := DoSomething()
	if err != nil {
		log.Printf("log for debug: %s", err.(*goerr.Error).Debug())
		log.Printf("output to outside: %s", err) // without change exiting code
	}

	_, err = DoSomething2()
	if err != nil {
		err2 := err.(*goerr.Error)
		log.Printf("log for debug: %s, %v", err2.Debug(), err2.GetAll())
		log.Printf("output to outside: %s", err)
	}
}
