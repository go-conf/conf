package conf

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	//"reflect"
	"testing"

	"gopkg.in/conf.v0"
)

type jsonExample struct {
	JsonFile string
	Name     string
	hidden   int
	Renamed1 string `conf:"blah"`
}

func TestJson(t *testing.T) {
	assert := assert.New(t)
	c := conf.New("json")
	jsonpath, err := tempJson()
	assert.NoError(err)
	dest := &jsonExample{Name: "Hello", Renamed1: "Bye", JsonFile: jsonpath}
	assert.NoError(c.Parse(dest))
	assert.Equal("Bob", dest.Name)
	assert.Equal("Bye", dest.Renamed1)
	os.Remove(jsonpath)
}

func tempJson() (string, error) {
	f, err := ioutil.TempFile("", "jsonexample")
	if err != nil {
		return "", err
	}
	path := f.Name()
	f.Write([]byte(`{"Name":"Bob","Age": 26}`))
	f.Close()
	return path, nil
}
