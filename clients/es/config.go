// Package es contain wrapper and helper for elasticsearch client
package es

type Config struct {
	Addresses []string `json:"addresses" yaml:"addresses" env:"ADDRESSES" envSeparator:"," validate:"notEmpty"`
	Username  string   `json:"username" yaml:"username" env:"USERNAME"`
	Password  string   `json:"password" yaml:"password" env:"PASSWORD"`
}
