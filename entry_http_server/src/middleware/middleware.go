package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	constant "http_server/src/common"
	"http_server/src/logger"

	"github.com/unrolled/secure"
	"time"
)

// jwt 相关结构体

type Claims struct {
	UserId string
	jwt.StandardClaims
}

func GenerateJwtToken(userId string) string {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			Issuer:   "zhifeng.wei",      // 签名颁发者
			Subject:  "user login token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	tokenString, err := token.SignedString(constant.JwtKey)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("Generate token failed,%v", err))
		return ""
	}
	return tokenString
}
func PhaseJwtToken(token string) string {
	claims := &Claims{}
	phasedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return constant.JwtKey, nil
	})
	if err != nil || !phasedToken.Valid {
		return ""
	}
	userid := phasedToken.Claims.(*Claims).UserId
	return userid
}
func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
