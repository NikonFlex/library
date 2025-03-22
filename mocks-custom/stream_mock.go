package mocks_custom

import (
	context "context"
	"sync"

	"github.com/project/library/generated/api/library"
	"google.golang.org/grpc/metadata"
)

type MockLibraryGetAuthorBooksServer struct {
	SentBooks []*library.Book
	SendError error
	CancelCtx context.CancelFunc
	mu        sync.Mutex
	ctx       context.Context
}

func NewMockLibraryGetAuthorBooksServer(cancel context.CancelFunc, sendError error) *MockLibraryGetAuthorBooksServer {
	return &MockLibraryGetAuthorBooksServer{
		ctx:       context.Background(),
		SendError: sendError,
		CancelCtx: cancel,
	}
}

// RecvMsg implements library.Library_GetAuthorBooksServer.
func (*MockLibraryGetAuthorBooksServer) RecvMsg(_ any) error {
	panic("unimplemented")
}

// SendHeader implements library.Library_GetAuthorBooksServer.
func (m *MockLibraryGetAuthorBooksServer) SendHeader(metadata.MD) error {
	panic("unimplemented")
}

// SendMsg implements library.Library_GetAuthorBooksServer.
func (*MockLibraryGetAuthorBooksServer) SendMsg(_ any) error {
	panic("unimplemented")
}

// SetHeader implements library.Library_GetAuthorBooksServer.
func (m *MockLibraryGetAuthorBooksServer) SetHeader(metadata.MD) error {
	panic("unimplemented")
}

// SetTrailer implements library.Library_GetAuthorBooksServer.
func (m *MockLibraryGetAuthorBooksServer) SetTrailer(metadata.MD) {
	panic("unimplemented")
}

// Send сохраняет отправленные книги и возвращает заданную ошибку
func (m *MockLibraryGetAuthorBooksServer) Send(book *library.Book) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.SendError == nil {
		m.SentBooks = append(m.SentBooks, book)
	}
	return m.SendError
}

// Context возвращает контекст мока
func (m *MockLibraryGetAuthorBooksServer) Context() context.Context {
	return m.ctx
}

// Reset очищает состояние мока
func (m *MockLibraryGetAuthorBooksServer) Reset() {
	panic("unimplemented")
}

// WithCancel добавляет возможность отмены контекста
func (m *MockLibraryGetAuthorBooksServer) WithCancel() *MockLibraryGetAuthorBooksServer {
	panic("unimplemented")
}

// Cancel вызывает отмену контекста
func (m *MockLibraryGetAuthorBooksServer) Cancel() {
	if m.CancelCtx != nil {
		m.CancelCtx()
	}
}
