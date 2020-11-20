package membership

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/data-portal-server/util/razorpay"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create membership
// @Summary Create membership
// @Description Create membership
// @Tags Membership
// @ID add-membership
// @Consume json
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param Membership body membership true "Membership object"
// @Success 201 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := util.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	membership := &membership{}
	err = json.NewDecoder(r.Body).Decode(&membership)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(membership)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Membership{
		Status: "created",
		UserID: uint(uID),
		PlanID: membership.PlanID,
	}

	tx := model.DB.Begin()

	// Check if the plan is not deleted
	plan := model.Plan{}
	plan.ID = membership.PlanID
	if err = tx.Model(&model.Plan{}).First(&plan).Error; err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	err = tx.Model(&model.Membership{}).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	// Fetch Currency
	currency := model.Currency{}
	model.DB.Model(&model.Currency{}).First(&currency)

	// Create a razorpay order and get razorpay orderID
	razorpayRequest := map[string]interface{}{
		"amount":          plan.Price * 100,
		"currency":        currency.IsoCode,
		"receipt":         fmt.Sprint(result.ID),
		"payment_capture": 1,
	}

	orderBody, err := razorpay.Client.Order.Create(razorpayRequest, nil)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	// Change membership status to initiated and add razorpay_id in membership table
	if _, found := orderBody["id"]; !found {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	result.Status = "processing"
	result.RazorpayOrderID = orderBody["id"].(string)

	err = tx.Model(&result).Updates(result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Preload("Plan").Preload("Plan.Catalogs").Preload("Plan.Catalogs.Products").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":         result.ID,
		"kind":       "membership",
		"status":     result.Status,
		"user_id":    result.UserID,
		"payment_id": result.PaymentID,
		"plan_id":    result.PlanID,
	}

	err = meili.AddDocument(meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()

	renderx.JSON(w, http.StatusCreated, result)
}
