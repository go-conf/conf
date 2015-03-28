package conf

import (
	"encoding/json"
	"errors"
	"os"

	"gopkg.in/conf.v0"
)

type JsonLoader struct {
	JsonFile string
}

func (l *JsonLoader) Load(c *conf.Conf) error {
	if l.JsonFile == "" {
		jsonField := c.FieldByName("JsonFile")
		if jsonField == nil {
			return errors.New("Could not find JSON configuration")
		}
		p := jsonField.Get()
		if p != nil {
			jf := p.(string)
			l.JsonFile = jf
		}
	}
	fp, err := os.Open(l.JsonFile)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(fp)
	err = dec.Decode(c.Dest)
	return err
}

func jsonFactory(*conf.Conf) conf.Loader {
	return &JsonLoader{}
}

func init() {
	conf.RegisterLoader("json", jsonFactory)
}
