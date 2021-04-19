package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type jsonWebKeys struct {
	Keys []jsonWebKey `json:"keys"`
}

type jsonWebKey struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type claimsWithScope struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		aud := "https://bundlemc.io/auth/users"
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("invalid audience")
		}
		// Verify 'iss' claim
		iss := "https://bundle.us.auth0.com/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("invalid issuer")
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

func GetAuthToken() (*Authorization, error) {
	authURL := "https://bundle.us.auth0.com/oauth/token"

	secret := viper.GetString("Auth0Secret")

	values := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {"22oXY4A0h9Rfbo3XEAn8Fbptx715dBe4"},
		"client_secret": {secret},
		"audience":      {"https://bundlemc.io/auth/users"},
	}

	authRes, err := http.PostForm(authURL, values)
	if err != nil {
		return nil, err
	}

	defer authRes.Body.Close()
	body, _ := ioutil.ReadAll(authRes.Body)

	auth := &Authorization{}

	err = json.Unmarshal(body, auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

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
				asJSON, err := json.Marshal(user)
				if err != nil {
					WriteResponse(w, err.Error(), http.StatusBadRequest)
					return
				}
				userJSON = string(asJSON)
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
				w.Write([]byte("password does not match stored password"))
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

	var jwks = jsonWebKeys{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func CheckScope(scope string, tokenString string) bool {

	token, _ := jwt.ParseWithClaims(tokenString, &claimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims := token.Claims.(*claimsWithScope)

	return strings.Contains(claims.Scope, scope)
}
