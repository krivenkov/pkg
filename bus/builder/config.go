package builder

import (
	"github.com/krivenkov/pkg/bus/franz"
)

type Transport string

const (
	TransportFranz Transport = "franz"
)

type Config struct {
	Transport Transport    `json:"transport" yaml:"transport" env:"TRANSPORT" validate:"notEmpty"`
	Franz     franz.Config `json:"franz" yaml:"franz" envPrefix:"FRANZ_"`
}
