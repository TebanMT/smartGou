package domain

type DomainErrorCodes string

var (
	DataValidationCode DomainErrorCodes = "01"
	UnexpectedErrors   DomainErrorCodes = "00"
)

type CustomError struct {
	Err       error
	ErrorCode DomainErrorCodes
}

func (e CustomError) Error() string {
	return e.Err.Error()
}

func (e CustomError) Unwrap() error {
	return e.Err
}

func NewValidationError(err error) CustomError {
	return CustomError{
		Err:       err,
		ErrorCode: DataValidationCode,
	}
}

func NewUnexpectedError(err error) CustomError {
	return CustomError{
		Err:       err,
		ErrorCode: UnexpectedErrors,
	}
}
