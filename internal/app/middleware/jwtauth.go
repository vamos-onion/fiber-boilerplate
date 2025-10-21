package middleware

import (
	"errors"
	"strings"

	"fiber-boilerplate/internal/app/config"
	logging "fiber-boilerplate/internal/pkg/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"
)

const ContextKeyStore = "fiber-boilerplate/#/keyStore"

func validateAPIKey(c *fiber.Ctx, tokenString string) (bool, error) {
	// keyauth가 Authorization 헤더에서 Bearer를 제거하지 못했을 경우 수동 처리
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}

	// JWT 파서 생성 (자동 expiration 검증 포함)
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}), // HS256만 허용
		jwt.WithExpirationRequired(),            // exp 클레임 필수 및 자동 검증
	)

	// JWT 파싱 및 검증
	token, err := parser.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 서명 알고리즘 확인 (보안상 중요)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.Server.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		logging.Debug("JWT validation failed: err=%v, valid=%v", err, token.Valid)
		return false, err
	}

	// 클레임을 컨텍스트에 저장
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		/*
			{
				"uuid": uuid string
			}
		*/

		// 사용자 정보 컨텍스트에 저장
		c.Locals(ContextKeyStore, claims)
	}

	return true, nil
}
