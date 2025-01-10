package config

import "github.com/go-playground/validator/v10"

func NewValidator(conf *Config) *validator.Validate {
	return validator.New()
}
