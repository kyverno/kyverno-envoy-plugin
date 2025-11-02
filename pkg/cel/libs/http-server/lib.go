package httpserver

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
)

type lib struct{}

func Lib() cel.EnvOption {
	// create the cel lib env option
	return cel.Lib(&lib{})
}

func (c *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// register types
		ext.NativeTypes(
			reflect.TypeFor[CheckResponse](),
			ext.ParseStructTags(true),
		),
	}
}

func (*lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}
