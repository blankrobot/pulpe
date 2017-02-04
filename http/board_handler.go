package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/blankrobot/pulpe"
	"github.com/julienschmidt/httprouter"
)

var (
	defaultSettings = json.RawMessage([]byte(`{}`))
)

// NewBoardHandler returns a new instance of BoardHandler.
func NewBoardHandler(c pulpe.Client) *BoardHandler {
	h := BoardHandler{
		Router: httprouter.New(),
		Client: c,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.GET("/v1/boards", h.handleGetBoards)
	h.POST("/v1/boards", h.handlePostBoard)
	h.GET("/v1/boards/:id", h.handleGetBoard)
	h.DELETE("/v1/boards/:id", h.handleDeleteBoard)
	h.PATCH("/v1/boards/:id", h.handlePatchBoard)
	return &h
}

// BoardHandler represents an HTTP API handler for boards.
type BoardHandler struct {
	*httprouter.Router

	Client pulpe.Client

	Logger *log.Logger
}

// handlePostBoard handles requests to create a new board.
func (h *BoardHandler) handleGetBoards(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	boards, err := h.Client.Connect().BoardService().Boards()
	switch err {
	case nil:
		encodeJSON(w, boards, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handlePostBoard handles requests to create a new board.
func (h *BoardHandler) handlePostBoard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req pulpe.BoardCreate

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	if req.Settings == nil {
		req.Settings = &defaultSettings
	}

	board, err := h.Client.Connect().BoardService().CreateBoard(&req)
	switch err {
	case nil:
		w.WriteHeader(http.StatusCreated)
		encodeJSON(w, board, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handleGetBoard handles requests to fetch a single board.
func (h *BoardHandler) handleGetBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.Client.Connect()

	// Get the board
	board, err := session.BoardService().Board(pulpe.BoardID(id))
	if err != nil {
		if err == pulpe.ErrBoardNotFound {
			NotFound(w)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	if board.Settings == nil {
		board.Settings = &defaultSettings
	}

	// Get the board's lists
	board.Lists, err = session.ListService().ListsByBoard(pulpe.BoardID(id))
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	// Get the board's cards
	board.Cards, err = session.CardService().CardsByBoard(pulpe.BoardID(id))
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	encodeJSON(w, board, h.Logger)
}

// handleDeleteBoard handles requests to delete a single board and all of its content.
func (h *BoardHandler) handleDeleteBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	err := h.Client.Connect().BoardService().DeleteBoard(pulpe.BoardID(id))
	if err != nil {
		if err == pulpe.ErrBoardNotFound {
			NotFound(w)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePatchBoard handles requests to update a board.
func (h *BoardHandler) handlePatchBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req pulpe.BoardUpdate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	board, err := h.Client.Connect().BoardService().UpdateBoard(pulpe.BoardID(id), &req)
	switch err {
	case nil:
		encodeJSON(w, board, h.Logger)
	case pulpe.ErrBoardNotFound:
		NotFound(w)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}