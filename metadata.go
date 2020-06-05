package aeroqualaqy1

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings for the package
type Settings struct {
	Host       string `md:"host,required,"`
	Port       string `md:"port,required"`
	Username   string `md:"username,required"`
	Password   string `md:"password,required"`
	Instrument string `md:"instrument,required"`
	Entity     string `md:"entity,required"`
	Mappings   string `md:"mappings,required"`
}

// Input for the package
type Input struct {
}

// Output for the package
type Output struct {
	EspMqttMsg map[string]interface{} `md:"espMqttMsg"`
}

// ToMap converts from structure to a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

// FromMap converts fields in map to type specified in structure
func (i *Input) FromMap(values map[string]interface{}) error {
	return nil
}

// ToMap converts from structure to a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"espMqttMsg": o.EspMqttMsg,
	}
}

// FromMap converts from map to whatever type .
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	// Converts to string
	o.EspMqttMsg, err = coerce.ToObject(values["espMqttMsg"])

	if err != nil {
		return err
	}
	return nil
}
