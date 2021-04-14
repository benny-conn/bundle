package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type ClaimsWithScope struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		aud := "https://bundlemc.io/auth/users"
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("Invalid audience.")
		}
		// Verify 'iss' claim
		iss := "https://bundle.us.auth0.com/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			var userJSON string

			userJSON = r.Header.Get("User")
			if userJSON == "" {
				err := r.ParseForm()
				if err != nil {
					WriteResponse(w, err.Error(), http.StatusBadRequest)
					return
				}
				user := User{
					r.FormValue("username"),
					r.FormValue("email"),
					r.FormValue("password"),
				}
				asJson, err := json.Marshal(user)
				if err != nil {
					WriteResponse(w, err.Error(), http.StatusBadRequest)
					return
				}
				userJSON = string(asJson)
			}

			validatedUser, err := ValidateAndReturnUser(userJSON)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			dbUser, err := GetUserFromDatabase(validatedUser)

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

func AuthClient(next http.Handler) http.Handler {
	return jwtMiddleware.Handler(next)
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://bundle.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func CheckScope(scope string, tokenString string) bool {

	token, _ := jwt.ParseWithClaims(tokenString, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims := token.Claims.(*ClaimsWithScope)

	return strings.Contains(claims.Scope, scope)
}
