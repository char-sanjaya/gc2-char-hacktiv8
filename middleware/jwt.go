package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"gc2-university-api/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	StudentID uint   `json:"student_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(studentID uint, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := JWTClaims{
		StudentID: studentID,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "missing authorization token",
			})
		}

		// Mendukung format "Bearer <token>" maupun "<token>" langsung
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		secret := os.Getenv("JWT_SECRET")
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "invalid or expired authorization token",
			})
		}

		c.Set("student_id", claims.StudentID)
		return next(c)
	}
}
