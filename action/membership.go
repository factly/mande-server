package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// membership request body
type membership struct {
	Status    string `json:"status"`
	UserID    uint   `json:"user_id"`
	PaymentID uint   `json:"payment_id"`
	PlanID    uint   `json:"plan_id"`
}

// GetMemberships - Get all memberships
// @Summary Show all memberships
// @Description Get all memberships
// @Tags Membership
// @ID get-all-memberships
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Membership
// @Router /memberships [get]
func GetMemberships(w http.ResponseWriter, r *http.Request) {

	var memberships []model.Membership
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Preload("User").Preload("Plan").Preload("Payment").Preload("Payment.Currency").Model(&model.Membership{}).Find(&memberships)

	json.NewEncoder(w).Encode(memberships)
}

// GetMembership - Get membership by id
// @Summary Show a membership by id
// @Description Get membership by ID
// @Tags Membership
// @ID get-membership-by-id
// @Produce  json
// @Param id path string true "Membership ID"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships/{id} [get]
func GetMembership(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Membership{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Membership{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&req).Association("User").Find(&req.User)
	model.DB.Model(&req).Association("Plan").Find(&req.Plan)
	model.DB.Model(&req).Association("Payment").Find(&req.Payment)
	model.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)
	json.NewEncoder(w).Encode(req)
}

// CreateMembership - Create membership
// @Summary Create membership
// @Description Create membership
// @Tags Membership
// @ID add-membership
// @Consume json
// @Produce  json
// @Param Membership body membership true "Membership object"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships [post]
func CreateMembership(w http.ResponseWriter, r *http.Request) {

	req := &model.Membership{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "User", "Plan", "Payment")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Membership{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&req).Association("User").Find(&req.User)
	model.DB.Model(&req).Association("Plan").Find(&req.Plan)
	model.DB.Model(&req).Association("Payment").Find(&req.Payment)
	model.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)

	json.NewEncoder(w).Encode(req)
}

// UpdateMembership - Update membership by id
// @Summary Update a membership by id
// @Description Update membership by ID
// @Tags Membership
// @ID update-membership-by-id
// @Produce json
// @Consume json
// @Param id path string true "Membership ID"
// @Param Membership body membership false "Membership"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships/{id} [put]
func UpdateMembership(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	membership := &model.Membership{
		ID: uint(id),
	}
	req := &model.Membership{}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&membership).Updates(model.Membership{
		UserID:    req.UserID,
		PaymentID: req.PaymentID,
		PlanID:    req.PlanID,
		Status:    req.Status,
	})

	model.DB.First(&membership)
	model.DB.Model(&membership).Association("User").Find(&membership.User)
	model.DB.Model(&membership).Association("Plan").Find(&membership.Plan)
	model.DB.Model(&membership).Association("Payment").Find(&membership.Payment)
	model.DB.Model(&membership.Payment).Association("Currency").Find(&membership.Payment.Currency)
	json.NewEncoder(w).Encode(membership)
}

// DeleteMembership - Delete membership by id
// @Summary Delete a membership
// @Description Delete membership by ID
// @Tags Membership
// @ID delete-membership-by-id
// @Consume  json
// @Param id path string true "Membership ID"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships/{id} [delete]
func DeleteMembership(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	membership := &model.Membership{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&membership).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&membership)

	json.NewEncoder(w).Encode(membership)
}
