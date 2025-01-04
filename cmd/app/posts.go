package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/indiecodermm/go-social/internal/store"
)

type CreatePostInput struct {
	Title   string   `json:"title" validate:"required,min=1,max=100"`
	Content string   `json:"content" validate:"required,min=1,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostInput struct {
	Title   string   `json:"title" validate:"omitempty,min=1,max=100"`
	Content string   `json:"content" validate:"omitempty,min=1,max=1000"`
	Tags    []string `json:"tags"`
}

type postContextKey string

const (
	postKey postContextKey = "post"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostInput
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: Add UserID after auth
		UserID: 1,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) listPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.store.Posts.List(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := getPostFromContext(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := getPostFromContext(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var payload UpdatePostInput
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Title != "" {
		post.Title = payload.Title
	}

	if payload.Content != "" {
		post.Content = payload.Content
	}

	if len(payload.Tags) > 0 {
		post.Tags = payload.Tags
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := getPostFromContext(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.Delete(r.Context(), post.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		post, err := app.store.Posts.Get(r.Context(), postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				app.notFound(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) (*store.Post, error) {
	post, ok := r.Context().Value(postKey).(*store.Post)
	if !ok {
		return nil, errors.New("post not found in context")
	}

	return post, nil
}
