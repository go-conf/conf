package env

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/conf.v0"
)

func init() {
	conf.RegisterLoader("env", envFactory)
}

type envLoader struct {
	actions map[reflect.Type]envHandler
}

func (f *envLoader) Load(c *conf.Conf) error {
	for _, spec := range c.Fields() {
		keys := []string{spec.ConfName, strings.ToUpper(spec.ConfName)}
		var envVal string
		for _, k := range keys {
			envVal = os.Getenv(k)
			if envVal != "" {
				break
			}
		}
		proc := &envProc{f, c, spec, spec.GetField(), envVal}
		fh, ok := f.actions[spec.Type]
		if !ok {
			return errors.New("Could not find field type " + spec.Type.String())
		}
		err := fh(proc)
		if err != nil {
			return err
		}
	}
	return nil
}

func envFactory() conf.Loader {
	f := &envLoader{
		actions: setupEnvActions(),
	}
	return f
}

type envHandler func(*envProc) error

func setupEnvActions() map[reflect.Type]envHandler {
	actions := make(map[reflect.Type]envHandler)
	setup := func(i interface{}, s envHandler) {
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

type envProc struct {
	*envLoader
	conf   *conf.Conf
	spec   conf.Field
	field  reflect.Value
	envVal string
}

func procInt(proc *envProc) error {
	i, err := strconv.ParseInt(proc.envVal, 10, 64)
	if err != nil {
		return err
	}
	proc.field.SetInt(i)
	return nil
}

func procUint(proc *envProc) error {
	i, err := strconv.ParseUint(proc.envVal, 10, 64)
	if err != nil {
		return err
	}
	proc.field.SetUint(i)
	return nil
}

func procString(proc *envProc) error {
	// TODO do we need to work on handling intentionally blanking value?
	proc.field.SetString(proc.envVal)
	return nil
}
