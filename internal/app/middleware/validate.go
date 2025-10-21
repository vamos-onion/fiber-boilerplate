package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"fiber-boilerplate/internal/models"
	"fiber-boilerplate/internal/pkg/session"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func validate(ctx *fiber.Ctx, sessionKey string) (int, *session.DataBlock, error) {
	entity, err := models.Appuser.SearchAppusers(ctx.Context(), models.SearchAppusersParams{
		UUID: null.StringFrom(sessionKey),
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return http.StatusInternalServerError, nil, err
	}
	if len(entity) == 0 {
		return http.StatusUnauthorized, nil, nil
	}

	return http.StatusOK, session.New(&entity[0]), nil
}
