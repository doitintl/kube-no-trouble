package test

import "github.com/rs/zerolog"

func Setup() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}
