package conf

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type example1 struct {
	Name     string
	hidden   int
	Renamed1 string `conf:"blah"`
}

func TestExample1(t *testing.T) {
	assert := assert.New(t)
	c := New()
	c.breakdown(&example1{Name: "Hello", Renamed1: "Bye"})
	assert.Equal(2, len(c.fields))
	assert.Equal("Name", c.fields[0].RealName)
	assert.Equal("Name", c.fields[0].ConfName)
	assert.Equal(reflect.TypeOf(""), c.fields[0].Type)
	assert.Equal("Hello", c.fields[0].Default)
	assert.Equal(2, c.fields[1].Index)
	assert.Equal("blah", c.fields[1].ConfName)
}
