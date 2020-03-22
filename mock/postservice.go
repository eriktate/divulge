package mock

import (
	"context"

	"github.com/eriktate/divulge"
	"github.com/google/uuid"
)

type PostService struct {
	SavePostFn    func(ctx context.Context, post divulge.Post) (uuid.UUID, error)
	SavePostCount int

	PublishPostFn    func(ctx context.Context, id uuid.UUID) error
	PublishPostCount int

	RedactPostFn    func(ctx context.Context, id uuid.UUID) error
	RedactPostCount int

	FetchPostFn    func(ctx context.Context, id uuid.UUID) (divulge.Post, error)
	FetchPostCount int

	ListPostsByAccountFn    func(ctx context.Context, accountID uuid.UUID) ([]divulge.Post, error)
	ListPostsByAccountCount int

	RemovePostFn    func(ctx context.Context, id uuid.UUID) error
	RemovePostCount int

	Error error
}

func (m *PostService) SavePost(ctx context.Context, post divulge.Post) (uuid.UUID, error) {
	m.SavePostCount++

	if m.SavePostFn != nil {
		return m.SavePostFn(ctx, post)
	}

	return post.ID, m.Error
}

func (m *PostService) PublishPost(ctx context.Context, id uuid.UUID) error {
	m.PublishPostCount++

	if m.PublishPostFn != nil {
		return m.PublishPostFn(ctx, id)
	}

	return m.Error
}

func (m *PostService) RedactPost(ctx context.Context, id uuid.UUID) error {
	m.RedactPostCount++

	if m.RedactPostFn != nil {
		return m.RedactPostFn(ctx, id)
	}

	return m.Error
}

func (m *PostService) FetchPost(ctx context.Context, id uuid.UUID) (divulge.Post, error) {
	m.FetchPostCount++

	if m.FetchPostFn != nil {
		return m.FetchPostFn(ctx, id)
	}

	return divulge.Post{}, m.Error
}

func (m *PostService) ListPostsByAccount(ctx context.Context, accountID uuid.UUID) ([]divulge.Post, error) {
	m.ListPostsByAccountCount++

	if m.ListPostsByAccountFn != nil {
		return m.ListPostsByAccountFn(ctx, accountID)
	}

	return nil, m.Error
}

func (m *PostService) RemovePost(ctx context.Context, id uuid.UUID) error {
	m.RemovePostCount++

	if m.RemovePostFn != nil {
		return m.RemovePostFn(ctx, id)
	}

	return m.Error
}
