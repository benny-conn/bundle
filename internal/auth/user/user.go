package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/wrapper"
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

func checkScope(tokenString string, scopes ...string) bool {

	tokenUser, err := GetUserFromToken(tokenString)
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

func GetUserFromToken(tokenString string) (*api.User, error) {

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

	return &claims.User, nil

}

func ValidateToken(tokenString string) error {
	secret := viper.GetString("ClientSecret")

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

func AuthUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			userJSON := r.Header.Get("User")
			user := &api.User{}

			err := json.Unmarshal([]byte(userJSON), user)
			if err != nil {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			isValid := internal.IsUserValid(user)

			if !isValid {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			dbUser, err := wrapper.GetUserApi(user.Username, user.Email)

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
