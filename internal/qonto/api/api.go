package api

import (
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
)

func NewAPI(transferManager core.TransferManager) *qontoAPI {
	return &qontoAPI{
		manager: transferManager,
	}
}
