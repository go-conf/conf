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
	assert.Equal(2, c.fields[1].path[0])
	assert.Equal("blah", c.fields[1].ConfName)
}

type embed1 struct {
	example1
	Foo string
}

func TestEmbed(t *testing.T) {
	assert := assert.New(t)
	c := New()
	c.Load(&embed1{example1{Name: "Hello", Renamed1: "Bye"}, "foo"})
	assert.Equal(3, len(c.fields))
	assert.Equal([]int{1}, c.fields[0].path)
	assert.Equal("foo", c.fields[0].Default)
	assert.Equal("Foo", c.fields[0].ConfName)
	assert.Equal([]int{0, 0}, c.fields[1].path)
	assert.Equal("Hello", c.fields[1].Default)
	assert.Equal([]int{0, 2}, c.fields[2].path)
	assert.Equal("Bye", c.fields[2].Default)

}
