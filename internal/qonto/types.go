package qonto

type (
	Logger interface {
		Info(msg string, args ...interface{})
		Error(msg string, args ...interface{})

		// SubLogger creates new sublogger with provided instance name
		SubLogger(name string) Logger
	}
)
