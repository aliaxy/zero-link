// Package configvalidator validates service configuration structs at startup.
package configvalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
)

var validate = validator.New()

// MustValidate validates cfg against its struct tags and calls logx.Severef on failure.
func MustValidate(cfg any) {
	if err := validate.Struct(cfg); err != nil {
		logx.Severef("invalid config: %v", err)
	}
}
