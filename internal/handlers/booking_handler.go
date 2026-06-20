package handlers

import (
	"net/http"

	"lodge-system/internal/middleware"
	"lodge-system/internal/models"
	"lodge-system/internal/services"
	"lodge-system/pkg/utils"

	"github.com/google/uuid"
)

type BookingHandler struct {
	service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

// ─── Individual booking ───────────────────────────────────────────────────────

// CreateIndividual handles POST /api/v1/bookings/individual
func (h *BookingHandler) CreateIndividual(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	branchID := middleware.GetBranchIDFromContext(r.Context())

	var req models.CreateIndividualBookingRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking, err := h.service.CreateIndividual(orgID, branchID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, booking)
}

// CreateFromRequest handles POST /api/v1/booking-requests/{request_id}/materialise
func (h *BookingHandler) CreateFromRequest(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	branchID := middleware.GetBranchIDFromContext(r.Context())

	requestID, err := uuid.Parse(r.PathValue("request_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request ID")
		return
	}

	var matReq models.MaterialiseRequest
	if err := utils.DecodeJson(r, &matReq); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking, err := h.service.CreateFromRequest(orgID, branchID, requestID, &matReq)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, booking)
}

// ─── List & get ───────────────────────────────────────────────────────────────

// List handles GET /api/v1/bookings
func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	p := utils.ParsePagination(r)
	bookerType := r.URL.Query().Get("booker_type")
	bookingType := r.URL.Query().Get("booking_type")
	status := r.URL.Query().Get("status")

	bookings, total, err := h.service.List(orgID, bookerType, bookingType, status, p.Page, p.PageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.PaginatedResponse{
		Data:     bookings,
		Page:     p.Page,
		PageSize: p.PageSize,
		Total:    total,
	})
}

// GetByID handles GET /api/v1/bookings/{id}
func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	booking, err := h.service.GetByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, booking)
}

// ─── Status ───────────────────────────────────────────────────────────────────

// UpdateStatus handles PUT /api/v1/bookings/{id}/status
func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	var req models.UpdateBookingStatusRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateStatus(id, orgID, req.Status); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Status updated"})
}

// CheckIn handles PUT /api/v1/bookings/{id}/checkin
func (h *BookingHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	if err := h.service.CheckIn(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := h.service.GetByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, booking)
}

// CheckOut handles PUT /api/v1/bookings/{id}/checkout
func (h *BookingHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	if err := h.service.CheckOut(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := h.service.GetByID(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, booking)
}

// Cancel handles DELETE /api/v1/bookings/{id}
func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	if err := h.service.Cancel(id, orgID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Booking cancelled"})
}

// ─── Room assignments ─────────────────────────────────────────────────────────

// AssignRoom handles POST /api/v1/bookings/{id}/assignments
func (h *BookingHandler) AssignRoom(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	var req models.CreateRoomAssignmentRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	assignment, err := h.service.AssignRoom(id, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, assignment)
}

// ListAssignments handles GET /api/v1/bookings/{id}/assignments
func (h *BookingHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	assignments, err := h.service.ListAssignments(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, assignments)
}

// UpdateAssignment handles PUT /api/v1/bookings/{id}/assignments/{assign_id}
func (h *BookingHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	assignID, err := uuid.Parse(r.PathValue("assign_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	var req models.UpdateRoomAssignmentRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	assignment, err := h.service.UpdateAssignment(id, orgID, assignID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, assignment)
}

// RemoveAssignment handles DELETE /api/v1/bookings/{id}/assignments/{assign_id}
func (h *BookingHandler) RemoveAssignment(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	assignID, err := uuid.Parse(r.PathValue("assign_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	if err := h.service.RemoveAssignment(id, orgID, assignID); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Assignment removed"})
}

// CheckInAssignment handles PUT /api/v1/bookings/{id}/assignments/{assign_id}/checkin
func (h *BookingHandler) CheckInAssignment(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	assignID, err := uuid.Parse(r.PathValue("assign_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	if err := h.service.CheckInAssignment(id, orgID, assignID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Room checked in"})
}

// CheckOutAssignment handles PUT /api/v1/bookings/{id}/assignments/{assign_id}/checkout
func (h *BookingHandler) CheckOutAssignment(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	assignID, err := uuid.Parse(r.PathValue("assign_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid assignment ID")
		return
	}

	if err := h.service.CheckOutAssignment(id, orgID, assignID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Room checked out"})
}

// ─── Attendees ────────────────────────────────────────────────────────────────

// ListAttendees handles GET /api/v1/bookings/{id}/attendees
func (h *BookingHandler) ListAttendees(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	attendees, err := h.service.ListAttendees(id, orgID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, attendees)
}

// AddAttendee handles POST /api/v1/bookings/{id}/attendees
func (h *BookingHandler) AddAttendee(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	var req models.CreateAttendeeRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	attendee, err := h.service.AddAttendee(id, orgID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, attendee)
}

// UpdateAttendee handles PUT /api/v1/bookings/{id}/attendees/{attendee_id}
func (h *BookingHandler) UpdateAttendee(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	attendeeID, err := uuid.Parse(r.PathValue("attendee_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid attendee ID")
		return
	}

	var req models.UpdateAttendeeRequest
	if err := utils.DecodeJson(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	attendee, err := h.service.UpdateAttendee(id, orgID, attendeeID, &req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, attendee)
}

// RemoveAttendee handles DELETE /api/v1/bookings/{id}/attendees/{attendee_id}
func (h *BookingHandler) RemoveAttendee(w http.ResponseWriter, r *http.Request) {
	orgID, _ := middleware.GetOrgIDFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}
	attendeeID, err := uuid.Parse(r.PathValue("attendee_id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid attendee ID")
		return
	}

	if err := h.service.RemoveAttendee(id, orgID, attendeeID); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Attendee removed"})
}
