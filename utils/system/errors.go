package system

import "errors"

var (
	ErrNoFileManager        = errors.New("no file manager found")
	ErrNoTerminal           = errors.New("no terminal found")
	ErrUnsupportedOS        = errors.New("unsupported operating system")
	ErrEditorCommandMissing = errors.New("no supported editor command found")
	ErrUnsupportedEditor    = errors.New("unsupported editor target")
	ErrCustomEditorCommand  = errors.New("invalid custom editor command")
)
