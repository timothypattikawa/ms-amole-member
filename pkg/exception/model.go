package exception

type InternalServerError struct {
	message string
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{message: message}
}

func (e *InternalServerError) Error() string {
	return e.message
}

type BadReqeustError struct {
	message string
}

func NewBadReqeustError(message string) *BadReqeustError {
	return &BadReqeustError{message: message}
}

func (e *BadReqeustError) Error() string {
	return e.message
}

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{message: message}
}

func (e *NotFoundError) Error() string {
	return e.message
}

type Unauthorized struct {
	message string
}

func NewUnauthorized(message string) *Unauthorized {
	return &Unauthorized{message: message}
}

func (e *Unauthorized) Error() string {
	return e.message
}
