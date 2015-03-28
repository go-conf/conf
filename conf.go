package conf // import "gopkg.in/conf.v0"

import (
	"reflect"
	"strings"
)

type Conf struct {
	loaders []Loader
	fields  []FieldSpec
	Dest    interface{}
	destVal reflect.Value
}

/* Create a New Conf.

loaders specifies the configuration loaders that we want to use.
*/
func New(loaders ...interface{}) *Conf {
	c := &Conf{}
	for _, l := range loaders {
		switch loaderSpec := l.(type) {
		case string:
			factory, ok := registry[loaderSpec]
			// TODO: decide error condition if !ok
			if ok {
				c.loaders = append(c.loaders, factory(c))
			}
		case Loader:
			c.loaders = append(c.loaders, loaderSpec)
		}
	}
	return c
}

// TODO refactor
func (c *Conf) breakdown(s interface{}) {
	v := reflect.ValueOf(s)
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		// Don't know why you're doing this to me, structs only.
		panic("Value must be a struct.")
	}
	c.destVal = v
	t := v.Type()
	c.fields = make([]FieldSpec, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// TODO deal with anonymous structs and nested structs
		if field.PkgPath != "" {
			// skip anonymous fields
			continue
		}
		confName, _ := getConfName(field)
		spec := FieldSpec{
			conf:     c,
			Index:    i,
			RealName: field.Name,
			ConfName: confName,
			Type:     field.Type,
			Default:  v.Field(i).Interface(),
		}
		c.fields = append(c.fields, spec)
	}
}

func (c *Conf) Parse(dest interface{}) (err error) {
	c.Dest = dest
	c.breakdown(dest)
	for _, loader := range c.loaders {
		err = loader.Load(c)
		if err != nil {
			break
		}
	}
	return
}

func (c *Conf) Fields() []FieldSpec {
	return c.fields
}

func (c *Conf) FieldByName(name string) *FieldSpec {
	for _, f := range c.fields {
		if f.ConfName == name {
			return &f
		}
	}
	return nil
}

func getConfName(field reflect.StructField) (confName string, flagTag []string) {
	tags := strings.Split(field.Tag.Get("conf"), ",")
	confName = tags[0]
	if len(tags) > 1 {
		flagTag = tags[1:]
	}
	if confName == "" {
		confName = field.Name
	}
	return
}

type FieldSpec struct {
	conf     *Conf
	Index    int
	RealName string
	ConfName string
	Type     reflect.Type
	Default  interface{}
}

func (f FieldSpec) Get() interface{} {
	//v := reflect.Indirect(reflect.Value(c.Dest))
	return f.GetField().Interface()
}

func (f FieldSpec) GetField() reflect.Value {
	return f.conf.destVal.Field(f.Index)
}
