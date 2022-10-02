package header

import "go-pdf/pdfio"

type Version uint8

const (
	VersionInvalid Version = iota
	Version10
	Version11
	Version12
	Version13
	Version14
	Version15
	Version16
	Version17
	Version20
)

func ReadVersion(r pdfio.Reader) (Version, error) {
	v, err := r.ReadString(8)
	if err != nil {
		return VersionInvalid, err
	}

	switch v {
	case "%PDF-1.0":
		return Version10, nil
	case "%PDF-1.1":
		return Version10, nil
	case "%PDF-1.2":
		return Version10, nil
	case "%PDF-1.3":
		return Version10, nil
	case "%PDF-1.4":
		return Version10, nil
	case "%PDF-1.5":
		return Version10, nil
	case "%PDF-1.6":
		return Version10, nil
	case "%PDF-1.7":
		return Version10, nil
	case "%PDF-2.0":
		return Version10, nil
	default:
		return VersionInvalid, &ErrInvalidVersion{v}
	}
}

type ErrInvalidVersion struct {
	Version string
}

func (err *ErrInvalidVersion) Error() string {
	return "invalid version specifier: \"" + err.Version + "\""
}
