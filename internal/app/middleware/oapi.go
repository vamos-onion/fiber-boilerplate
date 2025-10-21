package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"fiber-boilerplate/internal/app/handlers"
	"fiber-boilerplate/internal/defs"
	api "fiber-boilerplate/internal/generated/serviceapi"
	logging "fiber-boilerplate/internal/pkg/logging"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	oapiRouter routers.Router
)

func init() {
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	swagger.Servers = []*openapi3.Server{{}}

	oapiRouter, err = gorillamux.NewRouter(swagger)
	if err != nil {
		panic(err)
	}
}

func oapiAuthenticationFunc(ctx *fiber.Ctx, skipAuth bool) func(c context.Context, input *openapi3filter.AuthenticationInput) error {
	return func(openapi3Ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// security가 필요없는 엔드포인트는 인증 스킵
		if skipAuth {
			return nil
		}

		// NOTE : jwt middleware can be placed here instead of working standalone
		if input.SecurityScheme.Scheme == "bearer" {
			if !strings.HasPrefix(input.RequestValidationInput.Request.Header.Get(fiber.HeaderAuthorization), "Bearer") {
				return defs.ErrUnauthorized
			}
		}
		if input.SecuritySchemeName == "appEngineCron" {
			if input.RequestValidationInput.Request.Header.Get("X-Appengine-Cron") != "true" {
				return defs.ErrUnauthorized
			}
		}
		return nil
	}
}

func oapiRequestValidate(ctx *fiber.Ctx) error {
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	// e.Use(middleware.OapiRequestValidator(swagger))

	// 위처럼 등록할 경우 모든 경로에 대해 openapi request validator 가 동작하고,
	// openapi spec 에 없는 경로는 모두 "Path not found" 에러가 리턴된다.
	// 참고 이슈: https://github.com/deepmap/oapi-codegen/issues/55
	// * 이슈 마지막 댓글대로 하면 validator 가 동작하지 않는다...

	logging.Debug("oapiRequestValidate called for path: %s", ctx.Path())

	// /api prefix를 제거하고 OpenAPI 라우터에 전달
	// OpenAPI spec의 base path가 /api이므로, 실제 path에서 /api를 제거해야 함
	pathWithoutAPI := ctx.Path()
	if len(pathWithoutAPI) >= 4 && pathWithoutAPI[:4] == "/api" {
		pathWithoutAPI = pathWithoutAPI[4:]
		if pathWithoutAPI == "" {
			pathWithoutAPI = "/"
		}
	}
	logging.Debug("Original path: %s, path for OpenAPI router: %s", ctx.Path(), pathWithoutAPI)

	rctx := ctx.Context()
	r := new(http.Request)

	err := fasthttpadaptor.ConvertRequest(rctx, r, true)
	if err != nil {
		logging.Debug("ConvertRequest error: %v", err)
		return err
	}

	// path를 OpenAPI base path 제외한 경로로 변경
	r.URL.Path = pathWithoutAPI

	route, _, errRoute := oapiRouter.FindRoute(r)
	if errors.Is(errRoute, routers.ErrPathNotFound) {
		logging.Debug("Path not found in OpenAPI spec: %s", ctx.Path())
		return ctx.Next()
	}
	if errRoute != nil {
		logging.Debug("FindRoute error for path %s: %v", ctx.Path(), errRoute)
	}

	// OpenAPI 스펙에서 security 요구사항 확인
	// security가 없거나 비어있으면 JWT 검증 스킵하도록 표시
	skipAuth := false
	if route != nil && route.Operation != nil {
		if route.Operation.Security == nil || len(*route.Operation.Security) == 0 {
			skipAuth = true
			ctx.Locals("oapi:skip_auth", true)
			logging.Debug("Skip auth for path: %s (no security required)", ctx.Path())
		} else {
			logging.Debug("Auth required for path: %s, security: %+v", ctx.Path(), route.Operation.Security)
		}
	} else {
		logging.Debug("Route or Operation is nil for path: %s, route=%v", ctx.Path(), route != nil)
	}

	// OpenAPI validation 수행 - 수정된 request 사용
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request: r,
		Route:   route,
		Options: &openapi3filter.Options{
			AuthenticationFunc: oapiAuthenticationFunc(ctx, skipAuth),
		},
	}

	err = openapi3filter.ValidateRequest(rctx, requestValidationInput)
	if err != nil {
		logging.Debug("OpenAPI validation error for path %s: %v", ctx.Path(), err)
		return handlers.SendError(ctx, http.StatusBadRequest, err)
	}

	return ctx.Next()
}
