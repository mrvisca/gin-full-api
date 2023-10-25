package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

func IsAuth() gin.HandlerFunc {
	return checkJWT()
}

func checkJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Token string di dapat dari header postman
		authHeader := c.Request.Header.Get("Authorization")
		// Mengambil token dari "Bearer <token>"
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			// Parse takes the token string and a function for looking up the key. The latter is especially
			// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
			// head of the token to identify which key to use, but the parsed token (head and claims) is provided
			// to the callback, providing flexibility.
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// fmt.Println(claims["user_id"], claims["user_role"])
				c.Set("jwt_user_id", claims["user_id"]) // Menentukan key yang akan di oper
				c.Set("jwt_isAdmin", claims["user_role"])
			} else {
				c.JSON(422, gin.H{"msg": "Invalid Token", "error": err})
				c.Abort()
				return
				// fmt.Println(err)
			}
		} else {
			c.JSON(422, gin.H{"msg": "Authorization token not provided"})
			c.Abort()
			return
		}
	}
}
