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

/*
Create a New Conf.

loaders specifies the configuration loaders that we want to use.
Loaders can be either string references to loaders registered with
conf, or it can be an object which implements Loader:

    conf.New("json", "flag", MyLoader{})
*/
func New(loaders ...interface{}) *Conf {
	c := &Conf{}
	for _, l := range loaders {
		switch loaderSpec := l.(type) {
		case string:
			factory, ok := registry[loaderSpec]
			// TODO: decide error condition if !ok
			if ok {
				c.loaders = append(c.loaders, factory())
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

/*
Actually load configuration.

dest must be a pointer to a struct to which you want configuration to be loaded into.

Each exported struct field will have configuration loaded into it, with the
expected configuration name being the same as the struct field name, unless
struct tags are used with the key conf:

	// Field is included in conf with the default name
	Field int

	// Field is loaded with flag/var/field name "my_name"
	Field int `conf:"my_name"`

	// Field gets name "name" and config option "noflag" and JSON name "FullName"
	Name string `conf:"name,noflag" json:"FullName"`

How each Loader behaves is up to the loader, but typically the loaders are set
up so that they only overwrite specified known values and the ordering allows
the precedence.
*/
func (c *Conf) Load(dest interface{}) (err error) {
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

// Get all Fields that this conf knows about.
func (c *Conf) Fields() []FieldSpec {
	return c.fields
}

// FieldByName gets the field with ConfName matching 'name'
// If field not found, returns nil.
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

// FieldSpec is metadata inferred from reflecting into the struct that is given to conf.
type FieldSpec struct {
	conf     *Conf
	Index    int          // The index within the struct
	RealName string       // The name of the field on the struct.
	ConfName string       // The name of the field used for conf, like building flags
	Type     reflect.Type // The type of the struct field
	Default  interface{}  // The default value provided on this struct field
}

// Get the current value of the field referenced by this FieldSpec
func (f FieldSpec) Get() interface{} {
	//v := reflect.Indirect(reflect.Value(c.Dest))
	return f.GetField().Interface()
}

// Get the reflect.Value pointing to the struct field this FieldSpec describes.
func (f FieldSpec) GetField() reflect.Value {
	return f.conf.destVal.Field(f.Index)
}
