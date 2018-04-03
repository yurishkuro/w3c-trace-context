package params

import (
	"strconv"

	crossdock "github.com/crossdock/crossdock-go"
)

const (
	Actor1    = "actor1"
	Actor2    = "actor2"
	Sampled   = "sampled"
	BitLength = "bit_length"

	// Server is used as an override for the location of the actor
	Server = "server"
)

func GetBool(t crossdock.T, name string) bool {
	fatals := crossdock.Fatals(t)

	val := t.Param(name)
	fatals.NotEmpty(val, "param %s must not be empty", name)

	b, err := strconv.ParseBool(val)
	fatals.NoError(err, "param %s must be true or false", name)

	return b
}
