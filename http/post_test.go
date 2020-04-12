package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/eriktate/divulge"
	"github.com/eriktate/divulge/http"
	"github.com/eriktate/divulge/mock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func structToString(s interface{}) string {
	data, _ := json.Marshal(s)
	return string(data)
}

func Test_GetPost(t *testing.T) {
	cases := []struct {
		name            string
		postService     divulge.PostService
		postID          string
		expectedCode    int
		expectedContent string
	}{
		{
			name:            "Happy path",
			postService:     &mock.PostService{},
			postID:          uuid.New().String(),
			expectedCode:    stdhttp.StatusOK,
			expectedContent: structToString(divulge.Post{}),
		},
		{
			name:            "Invalid postID",
			postService:     &mock.PostService{},
			postID:          "abc123",
			expectedCode:    stdhttp.StatusBadRequest,
			expectedContent: "id is improperly formatted",
		},
		{
			name:            "PostService error",
			postService:     &mock.PostService{Error: errors.New("forced")},
			postID:          uuid.New().String(),
			expectedCode:    stdhttp.StatusInternalServerError,
			expectedContent: "failed to fetch post",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// SETUP
			rec := httptest.NewRecorder()
			req, err := stdhttp.NewRequest(stdhttp.MethodGet, fmt.Sprintf("http://localhost/post/%s", c.postID), nil)
			if err != nil {
				t.Fatalf("unexpected error while building request: %s", err)
			}

			// RUN
			http.BuildRouter(c.postService, logrus.New()).ServeHTTP(rec, req)

			// ASSERT
			if rec.Code != c.expectedCode {
				t.Fatalf("expected status code %d, but got %d", c.expectedCode, rec.Code)
			}

			content := rec.Body.String()
			if content != c.expectedContent {
				t.Fatalf("expected body '%s', but got '%s'", c.expectedContent, content)
			}
		})
	}
}

func Test_GetPosts(t *testing.T) {
	accountID := uuid.New()
	cases := []struct {
		name            string
		postService     divulge.PostService
		query           string
		expectedCode    int
		expectedContent string
	}{
		{
			name:            "Happy path",
			postService:     &mock.PostService{},
			query:           fmt.Sprintf("?accountId=%s", accountID.String()),
			expectedCode:    stdhttp.StatusOK,
			expectedContent: "null",
		},
		{
			name:            "Invalid accountID",
			postService:     &mock.PostService{},
			query:           fmt.Sprintf("?accountId=abc123"),
			expectedCode:    stdhttp.StatusBadRequest,
			expectedContent: "accountId is improperly formatted",
		},
		{
			name:            "PostService error",
			postService:     &mock.PostService{Error: errors.New("forced")},
			query:           fmt.Sprintf("?accountId=%s", accountID.String()),
			expectedCode:    stdhttp.StatusInternalServerError,
			expectedContent: "failed to list posts",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// SETUP
			rec := httptest.NewRecorder()
			req, err := stdhttp.NewRequest(stdhttp.MethodGet, fmt.Sprintf("http://localhost/post%s", c.query), nil)
			if err != nil {
				t.Fatalf("unexpected error while building request: %s", err)
			}

			// RUN
			http.BuildRouter(c.postService, logrus.New()).ServeHTTP(rec, req)

			// ASSERT
			if rec.Code != c.expectedCode {
				t.Fatalf("expected status code %d, but got %d", c.expectedCode, rec.Code)
			}

			content := rec.Body.String()
			if content != c.expectedContent {
				t.Fatalf("expected body '%s', but got '%s'", c.expectedContent, content)
			}
		})
	}
}

func Test_PostPost(t *testing.T) {
	// FIXTURES
	postID := uuid.New()
	validPost := fmt.Sprintf(`{
"title": "test post",
"accountId": "%s",
"authorID": "%s"
}`, uuid.New().String(), uuid.New().String())
	invalidJson := "abc{123}"
	invalidPost := `{"title": "test post"}`

	// TEST CASES
	cases := []struct {
		name            string
		postService     divulge.PostService
		body            string
		expectedCode    int
		expectedContent string
	}{
		{
			name: "Happy path",
			postService: &mock.PostService{
				SavePostFn: func(ctx context.Context, post divulge.Post) (uuid.UUID, error) {
					return postID, nil
				},
			},
			body:            validPost,
			expectedCode:    stdhttp.StatusOK,
			expectedContent: postID.String(),
		},
		{
			name:            "Invalid JSON",
			postService:     &mock.PostService{},
			body:            invalidJson,
			expectedCode:    stdhttp.StatusBadRequest,
			expectedContent: "malformed json",
		},
		{
			name:            "Invalid Post",
			postService:     &mock.PostService{},
			body:            invalidPost,
			expectedCode:    stdhttp.StatusBadRequest,
			expectedContent: "the post was invalid",
		},
		{
			name:            "PostService error",
			postService:     &mock.PostService{Error: errors.New("forced")},
			body:            validPost,
			expectedCode:    stdhttp.StatusInternalServerError,
			expectedContent: "failed to save post",
		},
	}

	for _, c := range cases {
		// SETUP
		t.Run(c.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, err := stdhttp.NewRequest(stdhttp.MethodPost, "http://localhost/post", bytes.NewBufferString(c.body))
			if err != nil {
				t.Fatalf("unexpected error while building request: %s", err)
			}

			// RUN
			http.BuildRouter(c.postService, logrus.New()).ServeHTTP(rec, req)

			// ASSERT
			if rec.Code != c.expectedCode {
				t.Fatalf("expected status code %d, but got %d", c.expectedCode, rec.Code)
			}

			content := rec.Body.String()
			if content != c.expectedContent {
				t.Fatalf("expected body '%s', but got '%s'", c.expectedContent, content)
			}
		})
	}
}
