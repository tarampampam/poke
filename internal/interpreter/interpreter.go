package interpreter

import (
	_ "embed"
	"fmt"
)

//go:embed common.js
var script string

func Run() error {
	var runtime, err = NewRuntime()
	if err != nil {
		return err
	}

	v, err := runtime.RunString(script)
	if err != nil {
		return err
	}

	fmt.Println("RESULT:", v)

	return nil
}
