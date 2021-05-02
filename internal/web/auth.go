package web

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/form3tech-oss/jwt-go"
)

type CustomClaims struct {
	User *api.User `json:"user"`
	jwt.StandardClaims
}

func newAuthToken(user *api.User) (string, error) {

	secret := os.Getenv("AUTH0_SECRET")

	claims := CustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
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

func checkScope(tokenString string, scopes ...string) bool {

	tokenUser, err := getUserFromToken(tokenString)
	if err != nil {
		return false
	}

	isAuthorized := true
	for _, scope := range scopes {
		if !internal.Contains(tokenUser.Scopes, scope) {
			isAuthorized = false
		}
	}
	return isAuthorized
}

func getUserFromToken(tokenString string) (*api.User, error) {

	secret := os.Getenv("AUTH0_SECRET")
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}

	return claims.User, nil

}

func validateToken(tokenString string) error {
	secret := os.Getenv("AUTH0_SECRET")

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

func newAccessCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   600,
		HttpOnly: true,
	}
}

func loginGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("access_token")
		if err != nil {
			http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
			return
		}
		tokenString := token.Value

		tokenUser, err := getUserFromToken(tokenString)
		if err != nil {
			token.MaxAge = -1
			http.SetCookie(w, token)
			http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
			return
		}

		newToken, err := newAuthToken(tokenUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newCookie := newAccessCookie(newToken)

		http.SetCookie(w, newCookie)

		next.ServeHTTP(w, req)
	})
}

func noGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("access_token")
		if err != nil {
			next.ServeHTTP(w, req)
			return
		}
		tokenString := token.Value

		tokenUser, err := getUserFromToken(tokenString)
		if err != nil {
			token.MaxAge = -1
			http.SetCookie(w, token)
			next.ServeHTTP(w, req)
			return
		}

		newToken, err := newAuthToken(tokenUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newCookie := newAccessCookie(newToken)

		http.SetCookie(w, newCookie)

		next.ServeHTTP(w, req)
	})
}

func scopeGate(next http.Handler, scopes ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("access_token")
		if err != nil {
			http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
			return
		}
		tokenString := token.Value

		tokenUser, err := getUserFromToken(tokenString)
		if err != nil {
			token.MaxAge = -1
			http.SetCookie(w, token)
			http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
			return
		}

		newToken, err := newAuthToken(tokenUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newCookie := newAccessCookie(newToken)

		http.SetCookie(w, newCookie)

		next.ServeHTTP(w, req)
	})
}
