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
	ConnectorMsg map[string]interface{} `md:"connectorMsg"`
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
		"connectorMsg": o.ConnectorMsg,
	}
}

// FromMap converts from map to whatever type .
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	// Converts to string
	o.ConnectorMsg, err = coerce.ToObject(values["connectorMsg"])

	if err != nil {
		return err
	}
	return nil
}
