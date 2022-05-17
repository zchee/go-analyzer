# go-analyzer/pkgerrors

This analyzer rewrites the [github.com/pkg/errors](https://github.com/pkg/errors) (that has been deprecated) to the [fmt.Errorf](https://pkg.go.dev/fmt#Errorf) with `%w` verb provided after the go1.13.

Support functions:

- `github.com/pkg/errors.As(err error, target interface{}) bool`  
  - Replace to `errors.As(err error, target interface{}) bool`

- `github.com/pkg/errors.Cause(err error) error`   
  - **NO**

- `github.com/pkg/errors.Errorf(format string, args ...interface{}) error`  
  - Replace to `fmt.Errorf(format string, args ...interface{}) error`

- `github.com/pkg/errors.Is(err, target error) bool`   
  - Replace to `errors.Is(err, target error) bool`

- `github.com/pkg/errors.New(message string) error`   
  - Replace to `errors.New(message string) error`

- `github.com/pkg/errors.Unwrap(err error) error`   
  - Replace to `errors.Unwrap(err error) error`

- `github.com/pkg/errors.WithMessage(err error, message string) error`   
  - **?**

- `github.com/pkg/errors.WithMessagef(err error, format string, args ...interface{}) error`   
  - **?**

- `github.com/pkg/errors.WithStack(err error) error`   
  - **?**

- `github.com/pkg/errors.Wrap(err error, message string) error`   
  - **Rewrite** to `fmt.Errorf("message: %w", err) error`

- `github.com/pkg/errors.Wrapf(err error, format string, args ...interface{}) error`   
  - **Rewrite** to `fmt.Errorf("format: %w", args..., err) error`
