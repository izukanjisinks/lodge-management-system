package handlers

import (
	"net/http"

	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type ClientHandler struct {
	service *services.ClientService
}

func NewClientHandler(service *services.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

// ─── Individual ───────────────────────────────────────────────────────────────

func (h *ClientHandler) ListIndividual(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	clients, total, err := h.service.ListIndividual(search, status, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     clients,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *ClientHandler) GetIndividualByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	client, err := h.service.GetIndividualByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Individual client not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) CreateIndividual(w http.ResponseWriter, r *http.Request) {
	var c models.IndividualClient
	if err := utils.DecodeJson(r, &c); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.CreateIndividual(&c); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, c)
}

func (h *ClientHandler) UpdateIndividual(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	var updates models.IndividualClient
	if err := utils.DecodeJson(r, &updates); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client, err := h.service.UpdateIndividual(id, &updates)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) DeleteIndividual(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	if err := h.service.DeleteIndividual(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Individual client deleted successfully"})
}

// ─── Corporate ────────────────────────────────────────────────────────────────

func (h *ClientHandler) ListCorporate(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	clients, total, err := h.service.ListCorporate(search, status, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     clients,
		Page:     pag.Page,
		PageSize: pag.PageSize,
		Total:    total,
	})
}

func (h *ClientHandler) GetCorporateByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	client, err := h.service.GetCorporateByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Corporate client not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) CreateCorporate(w http.ResponseWriter, r *http.Request) {
	var c models.CorporateClient
	if err := utils.DecodeJson(r, &c); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.CreateCorporate(&c); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, c)
}

func (h *ClientHandler) UpdateCorporate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	var updates models.CorporateClient
	if err := utils.DecodeJson(r, &updates); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client, err := h.service.UpdateCorporate(id, &updates)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) DeleteCorporate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	if err := h.service.DeleteCorporate(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Corporate client deleted successfully"})
}
