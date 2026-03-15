package sl

import (
	"fmt"

	slog "golang.org/x/exp/slog"
)

func ErrWithOP(err error, op string) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error() + fmt.Sprintf(" OP: %s", op)),
	}
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
