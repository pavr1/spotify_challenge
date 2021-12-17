package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"spotify_challenge.com/src/application"
)

type CustomRenderer struct {
	StatusCode int
	Status     string
	Data       interface{}
}

func (c CustomRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, c.StatusCode)
	render.Respond(w, r, c)

	return nil
}

func NewRenderer(statusCode int, err error, data interface{}) render.Renderer {
	status := ""
	if err != nil {
		status = err.Error()
	} else {
		status = "Ok"
	}

	return CustomRenderer{
		StatusCode: statusCode,
		Status:     status,
		Data:       data,
	}
}

type HandlerImpl struct {
	App application.Application
}

func NewHandler(app application.Application) HandlerImpl {
	return HandlerImpl{
		App: app,
	}
}

func (h HandlerImpl) Write(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	isrc := vars["ISRC"]

	err := h.App.Write(ctx, isrc)
	if err != nil {
		render := NewRenderer(http.StatusInternalServerError, err, nil)
		render.Render(w, r)
	} else {
		render := NewRenderer(http.StatusOK, nil, nil)
		render.Render(w, r)
	}
}

func (h HandlerImpl) ReadByISRC(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	isrc := vars["ISRC"]

	track, err := h.App.ReadByISRC(ctx, isrc)
	if err != nil {
		render := NewRenderer(http.StatusInternalServerError, err, nil)
		render.Render(w, r)
	} else {
		if track == nil {
			render := NewRenderer(http.StatusOK, errors.New("track not found"), nil)
			render.Render(w, r)
		} else {
			render := NewRenderer(http.StatusOK, nil, nil)
			render.Render(w, r)
		}
	}
}

func (h HandlerImpl) ReadByArtist(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	name := vars["Name"]

	track, err := h.App.ReadByArtist(ctx, name)
	if err != nil {
		render := NewRenderer(http.StatusInternalServerError, err, nil)
		render.Render(w, r)
	} else {
		render := NewRenderer(http.StatusOK, nil, track)
		render.Render(w, r)
	}
}
