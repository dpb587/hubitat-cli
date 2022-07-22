package cmdflags

type ErrorCode struct {
	Err  error
	Code int
}

var _ error = ErrorCode{}

func (err ErrorCode) Error() string {
	return err.Err.Error()
}
