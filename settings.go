package leopard

import (
	"errors"
	"os"
	"reflect"
)

type SettingValue interface {
	GetValue() interface{}
	SetValue(interface{}) error
	GetDefault() *interface{}
}
type SettingsObject map[string]SettingValue

type SettingsRegister interface {
	// GetName of settings registry
	GetName() string

	// GetSettings get the settings map
	GetSettings() map[string]SettingValue
}

// TreeValue is a setting value used to make children in the settings tree
type TreeValue struct {
	Map map[string]SettingValue
}

// EnvValue is a setting value using the environment variables
type EnvValue struct {
	path         string
	defaultValue interface{}
}

func (e EnvValue) GetValue() interface{} {
	envVal := os.Getenv(e.path)
	if envVal == "" {
		return e.GetDefault()
	}
	return envVal
}

func (e EnvValue) SetValue(i interface{}) error {
	if reflect.TypeOf(i).Kind() != reflect.String {
		return errors.New("value is not a string so it can not be set in the environment")
	}
	return os.Setenv(e.path, i.(string))
}

func (e EnvValue) GetDefault() *interface{} {
	return &e.defaultValue
}

// EnvSetting get a setting value from the environment
func EnvSetting(path string) SettingValue {
	return EnvValue{
		path:         path,
		defaultValue: nil,
	}
}

// EnvSettingD get a setting value from the environment with default.
func EnvSettingD(path string, defaultValue interface{}) SettingValue {
	return EnvValue{
		path:         path,
		defaultValue: defaultValue,
	}
}

func (a *LeopardApp) registerRegistry(register SettingsRegister) {
	register.GetSettings()
}
