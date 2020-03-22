package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eriktate/divulge"
	"github.com/eriktate/divulge/mock"
	"github.com/eriktate/divulge/service"
	"github.com/google/uuid"
)

func Test_SavePost(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	mockFS := &mock.FileStore{}
	mockPS := &mock.PostService{}
	postService := service.NewPostService(mockPS, mockFS)

	post := divulge.Post{
		ID: uuid.New(),
	}

	// RUN
	id, err := postService.SavePost(ctx, post)

	// ASSERT
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if id != post.ID {
		t.Fatalf("unexpected post ID: %s", id.String())
	}
}

func Test_SavePost_FSError(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	mockFS := &mock.FileStore{Error: errors.New("forced")}
	mockPS := &mock.PostService{}
	postService := service.NewPostService(mockPS, mockFS)

	post := divulge.Post{
		ID: uuid.New(),
	}

	// RUN
	_, err := postService.SavePost(ctx, post)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_SavePost_PSError(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	mockFS := &mock.FileStore{}
	mockPS := &mock.PostService{Error: errors.New("forced")}
	postService := service.NewPostService(mockPS, mockFS)

	post := divulge.Post{
		ID: uuid.New(),
	}

	// RUN
	_, err := postService.SavePost(ctx, post)

	// ASSERT
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_FetchPost(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	testContent := "this is some test content"
	mockFS := &mock.FileStore{
		ReadFn: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(testContent), nil
		},
	}
	mockPS := &mock.PostService{}
	postService := service.NewPostService(mockPS, mockFS)

	id := uuid.New()

	// RUN
	post, err := postService.FetchPost(ctx, id)

	// ASSERT
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if post.Content != testContent {
		t.Fatalf("expected post.Content to be: %s", testContent)
	}
}

func Test_FetchPost_FSError(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	mockFS := &mock.FileStore{Error: errors.New("forced")}
	mockPS := &mock.PostService{}
	postService := service.NewPostService(mockPS, mockFS)

	id := uuid.New()

	// RUN
	_, err := postService.FetchPost(ctx, id)

	// ASSERT
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_FetchPost_PSError(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	mockFS := &mock.FileStore{}
	mockPS := &mock.PostService{Error: errors.New("forced")}
	postService := service.NewPostService(mockPS, mockFS)

	id := uuid.New()

	// RUN
	_, err := postService.FetchPost(ctx, id)

	// ASSERT
	if err == nil {
		t.Fatal("expected error")
	}
}
