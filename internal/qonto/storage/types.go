package storage

type (

	// Storage defines interface to be satisfied by concrete storage implementation
	Storage interface {
		// Wait runs provided wait function until it returns no error
		Wait(f WaiterFunc) error

		// Close closes underlying storage connection if supported by concrete implementation
		Close() error
	}
)
