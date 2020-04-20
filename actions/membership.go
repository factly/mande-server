package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/go-chi/chi"
)

// membership request body
type membership struct {
	Status    string `json:"status"`
	UserID    uint   `json:"user_id"`
	PaymentID uint   `json:"payment_id"`
	PlanID    uint   `json:"plan_id"`
}

// GetMembership - Get membership by id
// @Summary Show a membership by id
// @Description Get membership by ID
// @Tags Membership
// @ID get-membership-by-id
// @Produce  json
// @Param id path string true "Membership ID"
// @Success 200 {object} models.Membership
// @Router /memberships/{id} [get]
func GetMembership(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Membership{
		ID: uint(id),
	}

	models.DB.Model(&models.Membership{}).First(&req)

	models.DB.Model(&req).Association("User").Find(&req.User)
	models.DB.Model(&req).Association("Plan").Find(&req.Plan)
	models.DB.Model(&req).Association("Payment").Find(&req.Payment)
	models.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)
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
// @Success 200 {object} models.Membership
// @Router /memberships [post]
func CreateMembership(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.Membership{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	err := models.DB.Model(&models.Membership{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	models.DB.Model(&req).Association("User").Find(&req.User)
	models.DB.Model(&req).Association("Plan").Find(&req.Plan)
	models.DB.Model(&req).Association("Payment").Find(&req.Payment)
	models.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)

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
// @Success 200 {object} models.Membership
// @Router /memberships/{id} [put]
func UpdateMembership(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		log.Fatal(err)
	}

	membership := &models.Membership{
		ID: uint(id),
	}
	req := &models.Membership{}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&membership).Updates(models.Membership{
		UserID:    req.UserID,
		PaymentID: req.PaymentID,
		PlanID:    req.PlanID,
		Status:    req.Status,
	})

	models.DB.First(&membership)
	models.DB.Model(&membership).Association("User").Find(&membership.User)
	models.DB.Model(&membership).Association("Plan").Find(&membership.Plan)
	models.DB.Model(&membership).Association("Payment").Find(&membership.Payment)
	models.DB.Model(&membership.Payment).Association("Currency").Find(&membership.Payment.Currency)
	json.NewEncoder(w).Encode(membership)
}

// DeleteMembership - Delete membership by id
// @Summary Delete a membership
// @Description Delete membership by ID
// @Tags Membership
// @ID delete-membership-by-id
// @Consume  json
// @Param id path string true "Membership ID"
// @Success 200 {object} models.Membership
// @Router /memberships/{id} [delete]
func DeleteMembership(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		log.Fatal(err)
	}

	membership := &models.Membership{
		ID: uint(id),
	}

	models.DB.First(&membership)
	models.DB.Delete(&membership)

	json.NewEncoder(w).Encode(membership)
}
