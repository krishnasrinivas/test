package main

import (
	"fmt"

	"github.com/krishnasrinivas/test/pkg/minerrors"
)
import "errors"

func f1() error {
	return minerrors.NewError(errors.New("testing"))
}

func f2() error {
	return f1()
}

func f3() error {
	return f2()
}

func main() {
	err := f3()
	if err != nil {
		fmt.Println(err, err.(*minerrors.Error).Stack())
	}
}
