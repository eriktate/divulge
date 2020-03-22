package mock

import "context"

type FileStore struct {
	WriteFn    func(ctx context.Context, key string, data []byte) error
	WriteCount int

	ReadFn    func(ctx context.Context, key string) ([]byte, error)
	ReadCount int

	Error error
}

func (m *FileStore) Write(ctx context.Context, key string, data []byte) error {
	m.WriteCount++

	if m.WriteFn != nil {
		return m.WriteFn(ctx, key, data)
	}

	return m.Error
}

func (m *FileStore) Read(ctx context.Context, key string) ([]byte, error) {
	m.ReadCount++

	if m.ReadFn != nil {
		return m.ReadFn(ctx, key)
	}

	return nil, m.Error
}
