package errors

type Error struct {
	Code int
	Err  string
}

func New(typ int, err string) *Error {
	return &Error{Code: typ, Err: err}
}

func (e *Error) Error() string {
	return e.Err
}
