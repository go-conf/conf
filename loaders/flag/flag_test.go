package conf

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/conf.v0"
	"testing"
)

type flagExample struct {
	Name     string
	hidden   int
	Renamed1 string `conf:"blah"`
}

func TestFlag(t *testing.T) {
	assert := assert.New(t)
	args := []string{"--Name", "Bob", "--blah", "Hello"}
	c := conf.New(&FlagLoader{Args: args, actions: setupFlagActions()})
	dest := &flagExample{Name: "Hello", Renamed1: "Bye"}
	assert.NoError(c.Parse(dest))
	assert.Equal("Bob", dest.Name)
	assert.Equal("Hello", dest.Renamed1)
}
