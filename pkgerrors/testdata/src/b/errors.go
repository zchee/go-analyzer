package a

import (
	"fmt"

	"github.com/pkg/errors" // want `found use "github.com/pkg/errors" package`
)

func PkgErrorsNew() error {
	return errors.New("whoops") // ok
}

func PkgErrorsErrorf() error {
	foo := fmt.Sprintf("foo")
	return errors.Errorf("whoops: %s", foo) // want `found use location of the deprecated github.com/pkg/errors`
}
