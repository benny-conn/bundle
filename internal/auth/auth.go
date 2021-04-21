package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/form3tech-oss/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func GetAuthToken(scope ...string) (string, error) {

	secret := viper.GetString("ClientSecret")

	scopes := strings.Join(scope, ",")
	claims := CustomClaims{
		Scope: scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 43200,
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

func checkRole(role string, tokenString string) bool {

	secret := viper.GetString("ClientSecret")

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

	claims := token.Claims.(*CustomClaims)

	return role == claims.Scope
}

func AuthWithScope(next http.Handler, scopes ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenString := token.Value

		canContinue := true
		for _, scope := range scopes {
			ok := checkRole(scope, tokenString)
			if !ok {
				canContinue = false
			}
		}

		if !canContinue {
			http.Error(w, "unautorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ImplicitLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			var userJSON string

			userJSON = r.Header.Get("User")
			if userJSON == "" {
				err := r.ParseForm()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				user := bundle.User{
					r.FormValue("username"),
					r.FormValue("email"),
					r.FormValue("password"),
				}
				asJSON, err := json.Marshal(user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				userJSON = string(asJSON)
			}

			validatedUser, err := bundle.ValidateAndReturnUser(userJSON)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			dbUser, err := bundle.GetUser(*validatedUser)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(validatedUser.Password))

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Password does not match stored password"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
