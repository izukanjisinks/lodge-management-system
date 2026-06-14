package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	CorporateBookingTypeAccommodation = "accommodation"
	CorporateBookingTypeMeals         = "meals"
	CorporateBookingTypeConference    = "conference"
	CorporateBookingTypeEvent         = "event"

	CorporateBookingStatusPending   = "pending"
	CorporateBookingStatusApproved  = "approved"
	CorporateBookingStatusRejected  = "rejected"
	CorporateBookingStatusCancelled = "cancelled"
)

type CorporateBookingRequest struct {
	ID                   uuid.UUID  `json:"id"`
	OrgID                uuid.UUID  `json:"org_id"`
	BranchID             *uuid.UUID `json:"branch_id,omitempty"`
	CorProfileID         *uuid.UUID `json:"cor_profile_id,omitempty"`
	CompanyID            *uuid.UUID `json:"company_id,omitempty"`
	BookingType          string     `json:"booking_type"`
	Status               string     `json:"status"`
	ReasonForBooking     string     `json:"reason_for_booking,omitempty"`
	Notes                string     `json:"notes,omitempty"`
	AuthoriserName       string     `json:"authoriser_name,omitempty"`
	AuthoriserEmail      string     `json:"authoriser_email,omitempty"`
	AuthoriserPhone      string     `json:"authoriser_phone,omitempty"`
	AuthoriserTitle      string     `json:"authoriser_title,omitempty"`
	AuthoriserDepartment string     `json:"authoriser_department,omitempty"`
	AuthoriserGLCode     string     `json:"authoriser_gl_code,omitempty"`
	Documents            []string   `json:"documents,omitempty"`
	Payload              json.RawMessage   `json:"payload,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	// Joined fields
	CompanyName string `json:"company_name,omitempty"`
	BranchName  string `json:"branch_name,omitempty"`
	ProfileName string `json:"profile_name,omitempty"`
}

// ── Submission payloads per booking type ─────────────────────────────────────

type CorBookingCompanyInput struct {
	CompanyName string `json:"company_name"`
	RegNumber   string `json:"reg_number,omitempty"`
	TPIN        string `json:"tpin,omitempty"`
	Industry    string `json:"industry,omitempty"`
	Country     string `json:"country,omitempty"`
}

type CorBookingBranchInput struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

type CorBookingProfileInput struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	JobTitle   string `json:"job_title,omitempty"`
	Department string `json:"department,omitempty"`
}

type CorBookingAuthoriser struct {
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Title      string `json:"title,omitempty"`
	Department string `json:"department,omitempty"`
	GLCode     string `json:"gl_code,omitempty"`
}

type CorBookingGuestInput struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email,omitempty"`
	Phone              string `json:"phone,omitempty"`
	IdentificationCard string `json:"identification_card"`
	CheckIn            string `json:"check_in,omitempty"`
	CheckOut           string `json:"check_out,omitempty"`
	RoomType           string `json:"room_type,omitempty"`
}

type SubmitAccommodationRequest struct {
	OrgID            uuid.UUID              `json:"org_id"`
	BranchID         *uuid.UUID             `json:"branch_id,omitempty"`
	Company          CorBookingCompanyInput  `json:"company"`
	Branch           *CorBookingBranchInput  `json:"branch,omitempty"`
	Profile          CorBookingProfileInput  `json:"booked_by"`
	ReasonForBooking string                 `json:"reason_for_booking,omitempty"`
	Notes            string                 `json:"notes,omitempty"`
	Authoriser       *CorBookingAuthoriser  `json:"authoriser,omitempty"`
	Guests           []CorBookingGuestInput  `json:"guests"`
	Documents        []string               `json:"documents,omitempty"`
}

type CorMealItemInput struct {
	MenuItemID string  `json:"menu_item_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type CorMealGuestInput struct {
	FirstName          string             `json:"first_name"`
	LastName           string             `json:"last_name"`
	Email              string             `json:"email,omitempty"`
	IdentificationCard string             `json:"identification_card,omitempty"`
	MealItems          []CorMealItemInput `json:"meal_items,omitempty"`
}

type SubmitMealsRequest struct {
	OrgID            uuid.UUID             `json:"org_id"`
	BranchID         *uuid.UUID            `json:"branch_id,omitempty"`
	Company          CorBookingCompanyInput `json:"company"`
	Branch           *CorBookingBranchInput `json:"branch,omitempty"`
	Profile          CorBookingProfileInput `json:"booked_by"`
	ReasonForBooking string                `json:"reason_for_booking,omitempty"`
	PlanType         string                `json:"plan_type"`
	From             string                `json:"from"`
	To               string                `json:"to"`
	DietaryNotes     string                `json:"dietary_notes,omitempty"`
	Authoriser       *CorBookingAuthoriser `json:"authoriser,omitempty"`
	Guests           []CorMealGuestInput   `json:"guests"`
	Documents        []string              `json:"documents,omitempty"`
}

type CorConferenceGuestInput struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email,omitempty"`
	IdentificationCard string `json:"identification_card,omitempty"`
}

type SubmitConferenceRequest struct {
	OrgID            uuid.UUID                 `json:"org_id"`
	BranchID         *uuid.UUID                `json:"branch_id,omitempty"`
	Company          CorBookingCompanyInput     `json:"company"`
	Branch           *CorBookingBranchInput     `json:"branch,omitempty"`
	Profile          CorBookingProfileInput     `json:"booked_by"`
	ReasonForBooking string                    `json:"reason_for_booking,omitempty"`
	StartDate        string                    `json:"start_date"`
	EndDate          string                    `json:"end_date,omitempty"`
	StartTime        string                    `json:"start_time"`
	EndTime          string                    `json:"end_time,omitempty"`
	Attendees        int                       `json:"attendees"`
	Equipment        []string                  `json:"equipment,omitempty"`
	Notes            string                    `json:"notes,omitempty"`
	Authoriser       *CorBookingAuthoriser     `json:"authoriser,omitempty"`
	Guests           []CorConferenceGuestInput  `json:"guests,omitempty"`
	Documents        []string                  `json:"documents,omitempty"`
}

// MaterialiseGuestAssignment maps a guest (by index in payload.guests) to a real room picked by staff.
type MaterialiseGuestAssignment struct {
	GuestIndex int       `json:"guest_index"`
	RoomID     uuid.UUID `json:"room_id"`
}

// MaterialiseRequest is the body for POST /api/v1/booking-requests/:id/materialise
type MaterialiseRequest struct {
	Assignments []MaterialiseGuestAssignment `json:"assignments"`
}

type SubmitEventRequest struct {
	OrgID            uuid.UUID                 `json:"org_id"`
	BranchID         *uuid.UUID                `json:"branch_id,omitempty"`
	Company          CorBookingCompanyInput     `json:"company"`
	Branch           *CorBookingBranchInput     `json:"branch,omitempty"`
	Profile          CorBookingProfileInput     `json:"booked_by"`
	ReasonForBooking string                    `json:"reason_for_booking,omitempty"`
	EventType        string                    `json:"event_type"`
	StartDate        string                    `json:"start_date"`
	EndDate          string                    `json:"end_date,omitempty"`
	StartTime        string                    `json:"start_time"`
	EndTime          string                    `json:"end_time,omitempty"`
	Headcount        int                       `json:"headcount"`
	CateringRequired bool                      `json:"catering_required"`
	Notes            string                    `json:"notes,omitempty"`
	Authoriser       *CorBookingAuthoriser     `json:"authoriser,omitempty"`
	Guests           []CorConferenceGuestInput  `json:"guests,omitempty"`
	Documents        []string                  `json:"documents,omitempty"`
}
