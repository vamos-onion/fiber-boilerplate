package v1

import (
	"errors"
	"fmt"
	"net/http"

	"fiber-boilerplate/internal/defs"
	api "fiber-boilerplate/internal/generated/serviceapi"
	"fiber-boilerplate/internal/models"
	logging "fiber-boilerplate/internal/pkg/logging"
	"fiber-boilerplate/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

const (
	// MaxPaginationLimit is the maximum number of items that can be requested per page
	MaxPaginationLimit = 1000
	// MaxPaginationOffset is the maximum offset allowed for pagination
	MaxPaginationOffset = 1000000
)

// EntityResponse :
func EntityResponse(entity interface{}) (response api.EntityResponse) {
	refl := util.Struct.Reflect(entity, "db")

	for i := 0; i < refl.Type.NumField(); i++ {
		field := util.Struct.Field(refl, i)
		switch field.Tag {
		case "uuid":
			response.UUID = field.Value.Interface().(null.String).String
		case "created_at":
			cts := util.Time.UnixMilli(field.Value.Interface().(null.Time).Time)
			response.CreatedAt = &cts
		case "modified_at":
			mts := util.Time.UnixMilli(field.Value.Interface().(null.Time).Time)
			response.ModifiedAt = &mts
		}
	}

	return
}

// SendGeneric :
func SendGeneric(ctx *fiber.Ctx, code int, messageOpt ...string) error {
	response := &api.GenericResponse{
		Code: code,
	}

	switch len(messageOpt) {
	case 0:
		response.Message = http.StatusText(code)
	case 1:
		response.Message = messageOpt[0]
	default:
		panic(defs.ErrInvalid)
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// SendResponse :
func SendResponse(ctx *fiber.Ctx, code int, data interface{}) error {
	response := &api.GenericResponse{
		Code:    code,
		Message: http.StatusText(code),
		Data:    &data,
	}
	return ctx.Status(code).JSON(response)
}

// SendError :
func SendError(ctx *fiber.Ctx, code int, errs ...error) error {
	response := &api.GenericResponse{
		Code: code,
	}

	var err error
	if len(errs) == 0 || errs[0] == nil {
		err = defs.NewError(code)
	} else {
		err = errs[0]
	}

	var fiberError *fiber.Error
	switch {
	case errors.As(err, &fiberError):
		var e *fiber.Error
		errors.As(err, &e)
		response.Message = e.Message
		logging.Warn(defs.NewError(e.Code), "%+v", err)
	default:
		response.Message = err.Error()
		logging.Warn(err, "")
	}

	// Use the actual error code instead of StatusTeapot
	return ctx.Status(code).JSON(response)
}

// EntityListParam :
func EntityListParam(sortingParam *api.SortingQueryParam, paginationParam *api.PaginationQueryParam) (*models.SortingBlock, *models.PaginationBlock, error) {
	var sorting models.SortingBlock
	var pagination models.PaginationBlock

	if sortingParam != nil {
		var keys []string
		if sortingParam.Keys != nil {
			for _, key := range *sortingParam.Keys {
				keys = append(keys, key)
			}
		}

		var dirs []string
		if sortingParam.Dirs != nil {
			for _, dir := range *sortingParam.Dirs {
				dirs = append(dirs, dir)
			}
		}

		if len(keys) > 0 && len(dirs) > 0 {
			if len(keys) != len(dirs) {
				return nil, nil, fmt.Errorf("sorting parameter mismatch: keys=%d, dirs=%d", len(keys), len(dirs))
			}
			sorting.Provided = true
			for i := range keys {
				sorting.Orders = append(sorting.Orders, util.String.Words(keys[i], dirs[i]))
			}
		}
	}

	if paginationParam != nil {
		if paginationParam.Limit != nil && paginationParam.Page != nil {
			limit := *paginationParam.Limit
			page := *paginationParam.Page

			// Validate pagination parameters
			if limit <= 0 || page <= 0 {
				return nil, nil, fmt.Errorf("invalid pagination parameters: limit=%d, page=%d (must be positive)", limit, page)
			}
			if limit > MaxPaginationLimit {
				return nil, nil, fmt.Errorf("pagination limit too large: %d (max: %d)", limit, MaxPaginationLimit)
			}

			// Check for potential overflow
			offset := limit * (page - 1)
			if offset < 0 || offset > MaxPaginationOffset {
				return nil, nil, fmt.Errorf("pagination offset out of range: %d (max: %d)", offset, MaxPaginationOffset)
			}

			pagination.Provided = true
			pagination.Limit = limit
			pagination.Offset = offset
		}
	}

	return &sorting, &pagination, nil
}
