package sl

import (
	"fmt"

	"golang.org/x/exp/slog"
)

func ErrWithOP(err error, op string) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(fmt.Sprintf(err.Error(), " OP: "+op)),
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
