package addons

import (
	"strings"

	js "github.com/dop251/goja"
)

type (
	Report struct {
		ReportLevel ReportLevel
		Message     string
	}

	ReportLevel string
)

const (
	ReportLevelDebug ReportLevel = "debug"
	ReportLevelInfo  ReportLevel = "info"
	ReportLevelWarn  ReportLevel = "warning"
	ReportLevelError ReportLevel = "error"
)

type Reports struct {
	runtime *js.Runtime `js:"-"` // runtime instance
	stack   []Report    `js:"-"`
}

func NewReports(runtime *js.Runtime) *Reports { return &Reports{runtime: runtime} }

func (r *Reports) Push(args ...js.Value) {
	for _, arg := range args {
		var (
			obj    = arg.ToObject(r.runtime)
			report = Report{
				ReportLevel: ReportLevelDebug,
			}
		)

		if lvl := obj.Get("level"); lvl != nil {
			if value, isString := lvl.Export().(string); isString {
				report.ReportLevel = ReportLevel(strings.ToLower(value))
			}
		}

		if msg := obj.Get("message"); msg != nil {
			if value, isString := msg.Export().(string); isString {
				report.Message = value
			}
		}

		if report.Message != "" {
			r.stack = append(r.stack, report)
		}
	}
}

func (r *Reports) Stack() []Report { return r.stack }

func (r *Reports) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"reports",
		runtime.ToValue(r),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
