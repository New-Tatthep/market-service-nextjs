package custom_config

import (
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type CustomConfiguration struct {
	Test string `mapstructure:"test" json:"test"`
}

func (cfg *CustomConfiguration) String() string {
	return stringutil.Json(*cfg)
}

func (cfg *CustomConfiguration) ToMap() map[string]any {
	return convutil.Obj2Map(*cfg)
}
