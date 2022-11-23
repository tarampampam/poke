package addons

import (
	"crypto/md5" //nolint:gosec
	"crypto/sha256"
	"encoding/hex"

	js "github.com/dop251/goja"
)

type Hashing struct {
	runtime *js.Runtime
}

func NewHashing(runtime *js.Runtime) *Hashing { return &Hashing{runtime: runtime} }

func (h *Hashing) Md5(args ...js.Value) js.Value {
	if len(args) == 0 {
		return js.Undefined()
	}

	var hash = md5.Sum([]byte(args[0].String())) //nolint:gosec

	return h.runtime.ToValue(hex.EncodeToString(hash[:]))
}

func (h *Hashing) Sha256(args ...js.Value) js.Value {
	if len(args) == 0 {
		return js.Undefined()
	}

	var hash = sha256.New()

	hash.Write([]byte(args[0].String()))

	return h.runtime.ToValue(hex.EncodeToString(hash.Sum(nil)))
}

func (h *Hashing) Register(runtime *js.Runtime) error {
	return runtime.GlobalObject().DefineDataProperty(
		"hashing",
		runtime.ToValue(h),
		js.FLAG_FALSE, // writable
		js.FLAG_FALSE, // configurable
		js.FLAG_TRUE,  // enumerable
	)
}
