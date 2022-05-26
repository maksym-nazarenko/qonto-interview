package api

import (
	"context"

	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
)

type mockManager struct {
	err error
}

func newMockManager() *mockManager {
	return &mockManager{}
}

func (mm *mockManager) ProcessTransfers(ctx context.Context, request *core.Request) error {
	return mm.err
}

func (mm *mockManager) WithError(err error) *mockManager {
	mm.err = err
	return mm
}
