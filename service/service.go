package service

import (
	"context"
	"fmt"

	"github.com/eriktate/divulge"
	"github.com/google/uuid"
)

// A PostService implements the divulge.PostService interface.
type PostService struct {
	ps divulge.PostService
	fs divulge.FileStore
}

// NewPostService returns a new PostService.
func NewPostService(ps divulge.PostService, fs divulge.FileStore) PostService {
	return PostService{
		ps: ps,
		fs: fs,
	}
}

// SavePost saves the post content in a file store and then passes off to another PostService
// to persist the metdata.
func (s PostService) SavePost(ctx context.Context, post divulge.Post) (uuid.UUID, error) {
	if err := s.fs.Write(ctx, post.ContentPath, []byte(post.Content)); err != nil {
		return post.ID, fmt.Errorf("failed to write post content: %w", err)
	}

	id, err := s.ps.SavePost(ctx, post)
	if err != nil {
		return id, err
	}

	return id, nil
}

// PublishPost passes off to another PostService to publish a Post.
func (s PostService) PublishPost(ctx context.Context, id uuid.UUID) error {
	return s.ps.PublishPost(ctx, id)
}

// RedactPost passes off to another PostService to redact a Post.
func (s PostService) RedactPost(ctx context.Context, id uuid.UUID) error {
	return s.ps.RedactPost(ctx, id)
}

// FetchPost fetches the post content from a FileStore and then combines it with metadata
// fetched from another PostService.
func (s PostService) FetchPost(ctx context.Context, id uuid.UUID) (divulge.Post, error) {
	post, err := s.ps.FetchPost(ctx, id)
	if err != nil {
		return post, err
	}

	content, err := s.fs.Read(ctx, post.ContentPath)
	if err != nil {
		return post, fmt.Errorf("failed to fetch post content: %w", err)
	}

	post.Content = string(content)
	return post, err
}

func (s PostService) ListPosts(ctx context.Context, req divulge.ListPostsReq) ([]divulge.Post, error) {
	return s.ps.ListPosts(ctx, req)
}

// RemovePost passes off to another PostService to remove a Post.
func (s PostService) RemovePost(ctx context.Context, id uuid.UUID) error {
	return s.ps.RemovePost(ctx, id)
}
