package itunes

// ErrInvalidFormat is returned when an input file
// does not meet the expected itune library xml format
type ErrInvalidFormat struct {
    s string
}

func (err ErrInvalidFormat) Error() string {
    return err.s
}

// NewInvalidFormatError created a new invalid format error
// that formats to the given string
func NewInvalidFormatError(s string) ErrInvalidFormat {
    return ErrInvalidFormat{s}
}
