package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/andreyxaxa/Delayed-Notifier/internal/controller/restapi/v1/request"
	"github.com/andreyxaxa/Delayed-Notifier/internal/controller/restapi/v1/response"
	"github.com/andreyxaxa/Delayed-Notifier/internal/entity"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/channel"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/errs"
	"github.com/andreyxaxa/Delayed-Notifier/pkg/types/status"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary 	Create notification
// @Description Creates new notification
// @ID 			create
// @Tags 		notifications
// @Accept 		json
// @Produce 	json
// @Param 		request body request.CreateNotificationRequest true "Notification"
// @Success 	200 {object} response.CreateNotificationResponse
// @Failure 	400 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/notify [post]
func (r *V1) createNotification(ctx *fiber.Ctx) error {
	now := time.Now()
	var body request.CreateNotificationRequest

	decoder := json.NewDecoder(bytes.NewReader(ctx.Body()))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&body)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	err = r.v.Struct(body)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if body.Payload.Channel == channel.Email {
		if body.Payload.Email == nil {
			return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
		}
	} else {
		if body.Payload.Telegram == nil {
			return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
		}
	}

	ntf := entity.Notification{
		UID:       uuid.New(),
		Payload:   body.Payload,
		CreatedAt: now,
		UpdatedAt: now,
		SendAt:    body.SendAt,
		Status:    status.Pending,
	}

	err = r.n.CreateNotification(ctx.UserContext(), ntf)
	if err != nil {
		r.l.Error(err, "restapi - v1 - createNotification")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	resp := response.CreateNotificationResponse{
		UID:    ntf.UID,
		Status: ntf.Status,
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

// @Summary     Get notification status
// @Description Returns notification status by UUID
// @ID 			get
// @Tags        notifications
// @Produce     json
// @Param 		id path string true "Notification UUID"
// @Success 	200 {object} response.GetNotificationStatusResponse
// @Failure 	400 {object} response.Error
// @Failure 	404 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/notify/{id} [get]
func (r *V1) getStatus(ctx *fiber.Ctx) error {
	uidStr := ctx.Params("id")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	s, err := r.n.GetStatus(ctx.UserContext(), uid)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "notification not found")
		}
		r.l.Error(err, "restapi - v1 - getStatus")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	resp := response.GetNotificationStatusResponse{
		Status: s,
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

// @Summary 	Cancel notification
// @Description Cancel notification by UUID
// @ID 			cancel
// @Tags 		notifications
// @Produce 	json
// @Param 		id path string true "Notification UUID"
// @Success 	200 {object} response.CancelNotificationResponse
// @Failure 	400 {object} response.Error
// @Failure 	404 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/notify/{id} [delete]
func (r *V1) cancelNotification(ctx *fiber.Ctx) error {
	uidStr := ctx.Params("id")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	err = r.n.CancelNotification(ctx.UserContext(), uid)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "notification not found")
		}
		if errors.Is(err, errs.ErrAlreadyCancelled) {
			return errorResponse(ctx, http.StatusBadRequest, "already cancelled")
		}
		if errors.Is(err, errs.ErrAlreadySentOrFailed) {
			return errorResponse(ctx, http.StatusBadRequest, "already sent or failed")
		}
		r.l.Error(err, "restapi - v1 - cancelNotification")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	resp := response.CancelNotificationResponse{
		UID:    uid,
		Status: status.Cancelled,
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
