package eror

type Error struct {
	Msg string
	Tip string
}

func (e *Error) Error() string {
	return e.Msg
}
