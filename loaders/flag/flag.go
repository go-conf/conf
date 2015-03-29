package conf

import (
	"errors"
	"flag"
	"reflect"
	"sync"

	"gopkg.in/conf.v0"
)

var flagSetupOnce sync.Once

type FlagLoader struct {
	FlagSet *flag.FlagSet
	Args    []string
	actions map[reflect.Type]flagHandler
}

func (f *FlagLoader) Load(c *conf.Conf) error {
	if f.FlagSet == nil {
		f.FlagSet = flag.NewFlagSet("", flag.ContinueOnError)
	}
	var fillers []flagFill
	for _, spec := range c.Fields() {
		proc := &flagProc{f, c, spec, spec.GetField(), nil}
		fh, ok := f.actions[spec.Type]
		if !ok {
			return errors.New("Could not find field type " + spec.Type.String())
		}
		fh(proc)
		fillers = append(fillers, proc.fill)
	}
	err := f.FlagSet.Parse(f.Args)
	if err == nil {
		for _, filler := range fillers {
			filler()
		}
	}
	return err
}

func flagFactory() conf.Loader {
	f := &FlagLoader{}
	f.actions = setupFlagActions()
	return f
}

func init() {
	conf.RegisterLoader("flag", flagFactory)
}

type flagHandler func(*flagProc)

func setupFlagActions() map[reflect.Type]flagHandler {
	actions := make(map[reflect.Type]flagHandler)
	setup := func(i interface{}, s flagHandler) {
		actions[reflect.TypeOf(i)] = s
	}
	setup(int(0), procInt)
	setup(int8(0), procInt)
	setup(int16(0), procInt)
	setup(int32(0), procInt)
	setup(int64(0), procInt)
	setup(uint(0), procInt)
	setup(uint8(0), procUint)
	setup(uint16(0), procUint)
	setup(uint32(0), procUint)
	setup(uint64(0), procUint)
	setup("", procString)
	return actions
}

type flagFill func()

type flagProc struct {
	*FlagLoader
	conf  *conf.Conf
	spec  conf.FieldSpec
	field reflect.Value
	fill  flagFill
}

func procInt(proc *flagProc) {
	ptr := proc.FlagSet.Int64(proc.spec.ConfName, proc.field.Int(), "")
	proc.fill = func() {
		proc.field.SetInt(*ptr)
	}
}

func procUint(proc *flagProc) {
	ptr := proc.FlagSet.Uint64(proc.spec.ConfName, proc.field.Uint(), "")
	proc.fill = func() {
		proc.field.SetUint(*ptr)
	}
}

func procString(proc *flagProc) {
	ptr := proc.FlagSet.String(proc.spec.ConfName, proc.spec.Default.(string), "")
	proc.fill = func() {
		proc.field.SetString(*ptr)
	}
}
