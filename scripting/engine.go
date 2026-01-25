package scripting

import (
	"time"

	"github.com/dop251/goja"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

type OutputResult struct {
	Value     map[string]any
	Allow     bool
	LogOutput []string
}

type Engine interface {
	Run(code string, ctx map[string]any) (OutputResult, error)
}

type JsEngine struct{}

func (e *JsEngine) Run(code string, ctx map[string]any) (OutputResult, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	time.AfterFunc(200*time.Millisecond, func() {
		vm.Interrupt("halt")
	})

	var outputRes OutputResult

	vm.Set("ctx", ctx)
	vm.Set("log", func(msg string) {
		outputRes.LogOutput = append(outputRes.LogOutput, msg)
	})

	_, err := vm.RunString(code)

	if err != nil {
		return outputRes, err
	}

	return outputRes, nil
}

type LuaEngine struct{}

func (e *LuaEngine) Run(code string, ctx map[string]any) (OutputResult, error) {
	var outputRes OutputResult

	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("ctx", luar.New(L, ctx))
	L.SetGlobal("log", luar.New(L, func(msg string) {
		outputRes.LogOutput = append(outputRes.LogOutput, msg)
	}))

	L.DoString(code)

	return outputRes, nil
}
