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
	return nil, goerr.IOError.Wrap(errors.New("redis err")).Set("user", 123).SetString("internal error")
}

func main() {
	_, err := DoSomething()
	if err != nil {
		log.Printf("log for debug: %s", err.Error())
		log.Printf("output to outside: %s", err.(*goerr.Error).String())
	}

	_, err = DoSomething2()
	if err != nil {
		log.Printf("log for debug: %s, %v", err.Error(), err.(*goerr.Error).GetAll())
		log.Printf("output to outside: %s", err.(*goerr.Error).String())
	}
}
