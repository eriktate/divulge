package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eriktate/divulge"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GetPost handles HTTP requests for fetching a specific Post.
func GetPosts(ps divulge.PostService, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryAccountID := r.URL.Query().Get("accountId")
		accountID, err := uuid.Parse(queryAccountID)
		if err != nil {
			log.WithError(err).WithField("query_account_id", queryAccountID).Error("failed to parse uuid")
			badRequest(w, "accountId is improperly formatted")
			return
		}

		req := divulge.ListPostsReq{AccountID: accountID}

		posts, err := ps.ListPosts(r.Context(), req)
		if err != nil {
			log.WithError(err).WithField("req", req).Error("failed to list posts")
			serverError(w, "failed to list posts")
			return
		}

		data, err := json.Marshal(posts)
		if err != nil {
			log.WithError(err).Error("failed to marshal post")
			serverError(w, "something went wrong while fetching post")
			return
		}

		ok(w, data)
	}
}

// GetPost handles HTTP requests for fetching a specific Post.
func GetPost(ps divulge.PostService, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		id, err := uuid.Parse(postID)
		if err != nil {
			log.WithError(err).WithField("id", id).Error("failed to parse uuid")
			badRequest(w, "id is improperly formatted")
			return
		}

		post, err := ps.FetchPost(r.Context(), id)
		if err != nil {
			log.WithError(err).WithField("id", id).Error("failed to fetch post")
			serverError(w, "failed to fetch post")
			return
		}

		data, err := json.Marshal(post)
		if err != nil {
			log.WithError(err).Error("failed to marshal post")
			serverError(w, "something went wrong while fetching post")
			return
		}

		ok(w, data)
	}
}

// PostPost handles HTTP requests for posting a Post.
func PostPost(ps divulge.PostService, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("failed to read body")
			badRequest(w, "failed to read body")
			return
		}

		var post divulge.Post
		if err := json.Unmarshal(data, &post); err != nil {
			log.WithError(err).Error("failed to unmarshal post")
			badRequest(w, "malformed json")
			return
		}

		if err := post.Validate(); err != nil {
			log.WithError(err).Error("invalid post")
			badRequest(w, "the post was invalid")
			return
		}

		id, err := ps.SavePost(r.Context(), post)
		if err != nil {
			log.WithError(err).Error("failed to save post")
			serverError(w, "failed to save post")
			return
		}

		ok(w, []byte(id.String()))
	}
}
