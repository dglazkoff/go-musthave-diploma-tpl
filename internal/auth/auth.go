package auth

import (
	"fmt"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

func BuildJWTString(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userId,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIDFromRequest(request *http.Request) (string, bool) {
	tokenString := request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		logger.Log.Error("Error while parse token: ", err)
		return "", false
	}

	return claims.UserID, true
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	logger.Log.Debug("Valid token")
	return nil
}

func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")
		if tokenString == "" {
			logger.Log.Error("No token")
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		err := verifyToken(tokenString)
		if err != nil {
			logger.Log.Error("Invalid token: ", err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(writer, request)
	}
}
