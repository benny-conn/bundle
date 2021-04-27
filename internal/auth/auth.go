package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
	"github.com/form3tech-oss/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	User api.User `json:"user"`
	jwt.StandardClaims
}

func NewAuthToken(user *api.User) (string, error) {

	secret := viper.GetString("ClientSecret")

	claims := CustomClaims{
		User: *user,
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

func CheckScope(tokenString string, scopes ...string) bool {

	tokenUser, err := GetUserFromToken(tokenString)
	if err != nil {
		return false
	}

	isAuthorized := true
	for _, scope := range scopes {
		if !bundle.Contains(tokenUser.Scopes, scope) {
			isAuthorized = false
		}
	}
	return isAuthorized
}

func GetUserFromToken(tokenString string) (*api.User, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &claims.User, nil

}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	secret := viper.GetString("ClientSecret")

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
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("expired token")
		return nil, err
	}

	return claims, nil

}

func AuthReqWithScope(next http.Handler, scopes ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		tokenString := token.Value

		canContinue := CheckScope(tokenString, scopes...)

		if !canContinue {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		RefreshToken(next).ServeHTTP(w, r)
	})
}

func AuthReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		_, err = ValidateToken(token.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		RefreshToken(next).ServeHTTP(w, r)
	})
}

func Refresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("access_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		RefreshToken(next).ServeHTTP(w, r)
	})
}

func AuthUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			userJSON := r.Header.Get("User")
			user := api.User{}

			err := json.Unmarshal([]byte(userJSON), &user)
			if err != nil {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			isValid := bundle.IsUserValid(user)

			if !isValid {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			dbUser, err := storage.GetUser(user)

			if err != nil {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

			if err != nil {
				http.Error(w, "incorrect password", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func NewAccessCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   600,
		HttpOnly: true,
	}
}

func RefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("access_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokenString := token.Value

		tokenUser, err := GetUserFromToken(tokenString)
		if err != nil {
			token.MaxAge = -1
			http.SetCookie(w, token)
			next.ServeHTTP(w, req)
			return
		}

		newToken, err := NewAuthToken(tokenUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newCookie := NewAccessCookie(newToken)

		http.SetCookie(w, newCookie)

		next.ServeHTTP(w, req)
	})
}
