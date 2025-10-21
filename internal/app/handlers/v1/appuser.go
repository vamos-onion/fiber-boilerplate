package v1

import (
	"fmt"
	"net/http"

	api "fiber-boilerplate/internal/generated/serviceapi"
	"fiber-boilerplate/internal/models"
	"fiber-boilerplate/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func CreateAppuser(ctx *fiber.Ctx) error {
	var body api.CreateAppuserRequest
	if err := ctx.BodyParser(&body); err != nil {
		return SendError(ctx, http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
	}

	entity, err := models.Appuser.CreateAppuser(nil, ctx.Context(), models.CreateAppuserParams{
		Name:     null.StringFrom(body.Name),
		Birthday: null.TimeFrom(util.Time.ToTimeFromOapiDate(body.Birthday)),
		Gender:   models.GenderToNullString(string(body.Gender)),
		Withdraw: null.BoolFrom(false),
	})
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to create appuser: %w", err))
	}

	entityResp := EntityResponse(entity)
	return SendResponse(ctx, http.StatusOK, &api.Appuser{
		CreatedAt:  entityResp.CreatedAt,
		ModifiedAt: entityResp.ModifiedAt,
		UUID:       entityResp.UUID,
		Birthday:   util.Time.ToOapiDate(entity.Birthday.Time),
		Gender:     entity.Gender.String,
		Name:       entity.Name.String,
		Withdraw:   entity.Withdraw.Bool,
	})
}

func ListAppusers(ctx *fiber.Ctx, params api.ListAppusersParams) error {
	sorting, pagination, err := EntityListParam(params.Sorting, params.Pagination)
	if err != nil {
		return SendError(ctx, http.StatusBadRequest, fmt.Errorf("invalid list parameters: %w", err))
	}

	list, err := models.Appuser.SearchAppusers(ctx.Context(), models.SearchAppusersParams{
		UUID:     null.StringFromPtr(params.Uuid),
		Name:     null.StringFromPtr(params.Name),
		Gender:   null.StringFromPtr(params.Gender),
		Withdraw: null.BoolFromPtr(params.Withdraw),
		Options:  models.MakeListOptions(sorting, pagination).Parameterize(),
	})
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to search appusers: %w", err))
	}

	return SendResponse(ctx, http.StatusOK, list)
}

func UpdateAppuser(ctx *fiber.Ctx) error {
	var body api.Appuser
	if err := ctx.BodyParser(&body); err != nil {
		return SendError(ctx, http.StatusBadRequest, fmt.Errorf("failed to parse request body: %w", err))
	}

	/* Tx Begin */
	tx, qctx, err := models.SQL.BeginxContext(ctx.Context())
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer tx.Rollback()
	qtx := models.New(tx)
	entity, err := models.Appuser.UpdateAppuser(qtx, qctx, models.UpdateAppuserParams{
		UUID:     null.StringFrom(body.UUID),
		Name:     null.StringFrom(body.Name),
		Birthday: null.TimeFrom(util.Time.ToTimeFromOapiDate(body.Birthday)),
		Gender:   models.GenderToNullString(body.Gender),
		Withdraw: null.BoolFrom(body.Withdraw),
	})
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to update appuser: %w", err))
	}

	err = tx.Commit()
	if err != nil {
		return SendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction: %w", err))
	}
	/* Tx Commit */

	entityResp := EntityResponse(entity)
	return SendResponse(ctx, http.StatusOK, &api.Appuser{
		CreatedAt:  entityResp.CreatedAt,
		ModifiedAt: entityResp.ModifiedAt,
		UUID:       entityResp.UUID,
		Birthday:   util.Time.ToOapiDate(entity.Birthday.Time),
		Gender:     entity.Gender.String,
		Name:       entity.Name.String,
		Withdraw:   entity.Withdraw.Bool,
	})
}
