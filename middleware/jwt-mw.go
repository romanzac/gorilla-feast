package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/romanzac/gorilla-feast/infra/config"
	"net/http"
	"os"
	"strings"
	"time"
)

// JWTToken represents token string
type JWTToken struct {
	Token string `json:"token"`
}

// GenerateJWT creates and signs new token
func GenerateJWT(acct, fullname string) (JWTToken, error) {
	signingKeyBytes, _ := os.ReadFile(config.Cfg.Web.JWTPrivKey)
	signingKey, _ := jwt.ParseRSAPrivateKeyFromPEM(signingKeyBytes)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"acct": acct,
		"name": fullname,
	})
	signedToken, err := token.SignedString(signingKey)
	return JWTToken{signedToken}, err
}

// VerifyJWTToken checks for expiry and signature
func VerifyJWTToken(signedToken string) (jwt.Claims, error) {
	verifyKeyBytes, _ := os.ReadFile(config.Cfg.Web.JWTPubKey)
	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM(verifyKeyBytes)
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

// JWTHandler protects routes with JWT token
func JWTHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Stripe away the Bearer string
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyJWTToken(tokenString)
		if err != nil {
			http.Error(w, "Error verifying JWT token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Extract account info from the claims
		acct := claims.(jwt.MapClaims)["acct"].(string)
		r.Header.Set("acct", acct)
		next.ServeHTTP(w, r)
	})
}
