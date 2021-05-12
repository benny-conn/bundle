package gate

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Scopes []string `json:"scopes"`
	jwt.StandardClaims
}

func newAuthToken(scopes ...string) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	claims := CustomClaims{
		Scopes: scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Second).Unix(),
			Issuer:    "bundle",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func validateToken(tokenString string) error {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("expired token")
		return err
	}

	return nil
}

func simpleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
		var token string
		if len(authHeaderParts) == 3 {
			token = authHeaderParts[2]
		}

		err := validateToken(token)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func scopedAuth(next http.Handler, scopes ...string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
		token := authHeaderParts[1]

		hasScope := checkScope(token, scopes...)

		if !hasScope {
			internal.WriteResponse(rw, "insufficient scope", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

func checkScope(tokenString string, scopes ...string) bool {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return false
	}

	hasScope := true
	for _, v := range scopes {
		if !internal.Contains(claims.Scopes, v) {
			hasScope = false
		}
	}

	return hasScope
}

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			err := r.ParseMultipartForm(10 << 20)

			if err != nil {
				err = r.ParseForm()
				if err != nil {
					http.Error(w, "invalid user", http.StatusBadRequest)
					return
				}
			}

			user := &api.User{
				Username: r.FormValue("username"),
			}

			gs := NewGateService("", "")

			dbUser, err := gs.GetUser(user)

			if err != nil {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(r.FormValue("password")))

			if err != nil {
				http.Error(w, "incorrect password", http.StatusUnauthorized)
				return
			}

		}
		next.ServeHTTP(w, r)
	})
}
