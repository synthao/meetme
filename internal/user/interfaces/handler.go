package interfaces

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/synthao/meetme/internal/user/application"
	"github.com/synthao/meetme/internal/user/domain"
	"net/http"
	"strconv"
)

type Handler struct {
	service *application.Service
}

type errorResponse struct {
	Error string `json:"error"`
}

type H map[string]any

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func NewHandler(s *application.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var dto application.CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	id, err := h.service.Create(dto)
	if err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	JSON(w, http.StatusCreated, H{"id": id})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var dto application.CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	var (
		vars = mux.Vars(r)
		id   = vars["id"]
	)

	userID, err := strconv.Atoi(id)
	if err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	user := &domain.User{
		ID:        userID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Gender:    domain.Gender(dto.Gender),
	}

	if err := h.service.Update(user); err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	user, err := h.service.GetByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	dto := getByIDResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	JSON(w, http.StatusOK, dto)
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetList(10, 0)
	if err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	dto := make([]getListResponse, len(users))

	for i, u := range users {
		dto[i] = getListResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		}
	}

	JSON(w, http.StatusOK, dto)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	err = h.service.Delete(userID)
	if err != nil {
		JSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
