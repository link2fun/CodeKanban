package system

import "errors"

var (
	ErrNoFileManager = errors.New("no file manager found")
	ErrNoTerminal    = errors.New("no terminal found")
	ErrUnsupportedOS = errors.New("unsupported operating system")
)
