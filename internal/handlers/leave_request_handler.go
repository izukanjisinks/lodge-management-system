package handlers

import (
	"net/http"

	"hr-system/internal/interfaces"
	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type LeaveRequestHandler struct {
	service    *services.LeaveRequestService
	empService *services.EmployeeService
}

func NewLeaveRequestHandler(svc *services.LeaveRequestService, empSvc *services.EmployeeService) *LeaveRequestHandler {
	return &LeaveRequestHandler{service: svc, empService: empSvc}
}

func (h *LeaveRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondError(w, http.StatusBadRequest, "No employee record linked to your account")
		return
	}

	var req models.LeaveRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	req.EmployeeID = emp.ID

	if err := h.service.Create(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, req)
}

func (h *LeaveRequestHandler) GetMyRequests(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondJSON(w, http.StatusOK, []interface{}{})
		return
	}
	pag := utils.ParsePagination(r)
	reqs, total, err := h.service.List(interfaces.LeaveRequestFilter{EmployeeID: &emp.ID}, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list requests")
		return
	}
	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{Data: reqs, Page: pag.Page, PageSize: pag.PageSize, Total: total})
}

func (h *LeaveRequestHandler) List(w http.ResponseWriter, r *http.Request) {
	pag := utils.ParsePagination(r)
	filter := interfaces.LeaveRequestFilter{
		Status: r.URL.Query().Get("status"),
	}
	if v := r.URL.Query().Get("employee_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.EmployeeID = &id
		}
	}
	reqs, total, err := h.service.List(filter, pag.Page, pag.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to list requests")
		return
	}
	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{Data: reqs, Page: pag.Page, PageSize: pag.PageSize, Total: total})
}

func (h *LeaveRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}
	req, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Leave request not found")
		return
	}
	utils.RespondJSON(w, http.StatusOK, req)
}

func (h *LeaveRequestHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondError(w, http.StatusBadRequest, "No employee record linked to your account")
		return
	}
	if err := h.service.Cancel(id, emp.ID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Leave request cancelled"})
}

func (h *LeaveRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondError(w, http.StatusBadRequest, "No employee record for reviewer")
		return
	}
	var body struct{ Comment string `json:"comment"` }
	_ = utils.DecodeJson(r, &body)
	if err := h.service.Approve(id, emp.ID, body.Comment); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Leave request approved"})
}

func (h *LeaveRequestHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	emp := h.getEmployeeByUser(userID)
	if emp == nil {
		utils.RespondError(w, http.StatusBadRequest, "No employee record for reviewer")
		return
	}
	var body struct{ Comment string `json:"comment"` }
	_ = utils.DecodeJson(r, &body)
	if err := h.service.Reject(id, emp.ID, body.Comment); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Leave request rejected"})
}

func (h *LeaveRequestHandler) getEmployeeByUser(userID uuid.UUID) *models.Employee {
	emps, _, err := h.empService.List(interfaces.EmployeeFilter{}, 1, 500)
	if err != nil {
		return nil
	}
	for i := range emps {
		e := &emps[i]
		if e.UserID != nil && *e.UserID == userID {
			return e
		}
	}
	return nil
}
