package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
	"github.com/form3tech-oss/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Profile bundle.Profile `json:"profile"`
	jwt.StandardClaims
}

func NewAuthToken(profile bundle.Profile) (string, error) {

	secret := viper.GetString("ClientSecret")

	claims := CustomClaims{
		Profile: profile,
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

	tokenUser, err := GetProfileFromToken(tokenString)
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

func GetProfileFromToken(tokenString string) (bundle.Profile, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return bundle.Profile{}, err
	}

	return claims.Profile, nil

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

func AuthWithScope(next http.Handler, scopes ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		tokenString := token.Value

		canContinue := CheckScope(tokenString, scopes...)

		if !canContinue {
			http.Error(w, "unautorized", http.StatusUnauthorized)
			return
		}

		RefreshToken(next).ServeHTTP(w, r)
	})
}

func AuthWithoutScope(next http.Handler) http.Handler {
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

			var userJSON string

			userJSON = r.Header.Get("User")
			if userJSON == "" {
				err := r.ParseForm()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				user := bundle.User{
					Username: r.FormValue("username"),
					Email:    r.FormValue("email"),
					Password: r.FormValue("password"),
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

			dbUser, err := storage.GetUser(*validatedUser)

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

		tokenUser, err := GetProfileFromToken(tokenString)
		if err != nil {
			http.Redirect(w, req, "/logout", http.StatusUnauthorized)
			return
		}

		newToken, err := NewAuthToken(tokenUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		newCookie := NewAccessCookie(newToken)

		http.SetCookie(w, newCookie)

		next.ServeHTTP(w, req)
	})
}
