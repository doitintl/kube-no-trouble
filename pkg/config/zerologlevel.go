package config

import (
	"github.com/rs/zerolog"
)

type ZeroLogLevel int8

func (l ZeroLogLevel) String() string {
	return zerolog.Level(l).String()
}

func (l *ZeroLogLevel) Set(s string) error {
	level, err := zerolog.ParseLevel(s)
	if err != nil {
		return err
	}
	*l = ZeroLogLevel(level)
	return nil
}

func (l ZeroLogLevel) Type() string {
	return "string"
}
