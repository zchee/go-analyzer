package a

import (
	errorspkg "errors"
	"fmt"

	"github.com/pkg/errors" // ok TODO(zchee): support `found use "github.com/pkg/errors" package`
)

var _ = errorspkg.New
var _ fmt.State
var _ errors.Frame

func PkgErrorsCause(err error) error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")  // want `found use location of the deprecated github.com/pkg/errors`
	e3 := errors.Wrap(e2, "middle") // want `found use location of the deprecated github.com/pkg/errors`
	e4 := errors.Wrap(e3, "outer")  // want `found use location of the deprecated github.com/pkg/errors`
	return errors.Cause(e4)         // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsErrorf() error {
	return errors.Errorf("whoops: %s", "foo") // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsErrorf2() (any, error) {
	return nil, errors.Errorf("whoops: %s", "foo") // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsErrorf3() (any, struct{}, error) {
	return nil, struct{}{}, errors.Errorf("whoops: %s", "foo") // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsErrorfWithVerb() error {
	return errors.Errorf("withverb: %s: %v", "foo", errors.New("bar")) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsErrorfWithVerb2() (any, error) {
	return nil, errors.Errorf("withverb: %s: %v", "foo", errors.New("bar")) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsErrorfWithVerb3() (any, struct{}, error) {
	return nil, struct{}{}, errors.Errorf("withverb: %s: %v", "foo", errors.New("bar")) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsErrorfWithVVerb() error {
	return errors.Errorf("withverb: %s: %v", "foo", errors.New("bar")) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsNew() error {
	return errors.New("whoops") // ok
}

func pkgErrorsWithMessage() error {
	cause := errors.New("whoops")
	return errors.WithMessage(cause, "oh noes") // want `found use location of the deprecated github.com/pkg/errors`
}

func pkgErrorsWithMessage2() (any, error) {
	cause := errors.New("whoops")
	return nil, errors.WithMessage(cause, "oh noes") // want `found use location of the deprecated github.com/pkg/errors`
}

func pkgErrorsWithMessage3() (any, struct{}, error) {
	cause := errors.New("whoops")
	return nil, struct{}{}, errors.WithMessage(cause, "oh noes") // want `found use location of the deprecated github.com/pkg/errors`
}

func pkgErrorsWithMessagef() error {
	cause := errors.New("whoops")
	return errors.WithMessagef(cause, "oh noes: %s", "yeah") // want `found use location of the deprecated github.com/pkg/errors`
}

func pkgErrorsWithMessagef2() (any, error) {
	cause := errors.New("whoops")
	return nil, errors.WithMessagef(cause, "oh noes: %s", "yeah") // want `found use location of the deprecated github.com/pkg/errors`
}

func pkgErrorsWithMessagef3() (any, struct{}, error) {
	cause := errors.New("whoops")
	return nil, struct{}{}, errors.WithMessagef(cause, "oh noes: %s", "yeah") // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrap() error {
	cause := errors.New("whoops")
	return errors.Wrap(cause, "oh noes") // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrap2() (any, error) {
	cause := errors.New("whoops")
	return nil, errors.Wrap(cause, "oh noes") // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrap3() (any, struct{}, error) {
	cause := errors.New("whoops")
	return nil, struct{}{}, errors.Wrap(cause, "oh noes") // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrap4() (any, struct{}, error) {
	cause := errors.New("whoops")
	return nil, struct{}{}, errors.Wrap(cause, fmt.Sprintf("oh noes: %s", "yeah")) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrapf() error {
	cause := errors.New("whoops")
	return errors.Wrapf(cause, "oh noes #%d", 1) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrapf2() (any, error) {
	cause := errors.New("whoops")
	return nil, errors.Wrapf(cause, "oh noes #%d", 2) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrapf3() (any, struct{}, error) {
	cause := errors.New("whoops")
	return nil, struct{}{}, errors.Wrapf(cause, "oh noes #%d", 3) // want `found use location of the deprecated github.com/pkg/errors`
}

func PkgErrorsWrapf4() (any, struct{}, error) {
	cause := errors.New("whoops")
	return nil, struct{}{}, errors.Wrapf(cause, fmt.Sprintf("oh noes #%d", 3)) // want `found use location of the deprecated github.com/pkg/errors`
}
