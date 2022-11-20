package interpreter

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed common.js
var script string

func Run() error {
	var runtime, err = NewRuntime(context.TODO())
	if err != nil {
		return err
	}

	v, err := runtime.RunString(script)
	if err != nil {
		return err
	}

	fmt.Println(">> RESULT:", v)

	for _, r := range runtime.Reports() {
		fmt.Printf(">> REPORT (level %s): %s\n", r.ReportLevel, r.Message)
	}

	return nil
}
