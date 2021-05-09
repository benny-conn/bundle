package gate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

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

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		aud := "https://bundlemc.io/auth"
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

func simpleAuth(next http.Handler) http.Handler {
	return jwtMiddleware.Handler(next)
}

func scopedAuth(next http.Handler, scopes ...string) http.Handler {
	return jwtMiddleware.Handler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
		token := authHeaderParts[1]

		hasScope := checkScope(token, scopes...)

		if !hasScope {
			internal.WriteResponse(rw, "insufficient scope", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r)
	}))
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

func getAccessToken() (string, error) {

	u := "https://bundle.us.auth0.com/oauth/token"

	id := os.Getenv("AUTH0_ID")
	secret := os.Getenv("AUTH0_SECRET")
	aud := os.Getenv("AUTH0_AUD")

	form := url.Values{}

	form.Set("grant_type", "client_credentials")
	form.Set("client_id", id)
	form.Set("client_secret", secret)
	form.Set("audience", aud)

	res, err := http.PostForm(u, form)
	if err != nil {
		return "", nil
	}

	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	j := &struct {
		AccessToken string `json:"access_token"`
	}{}

	err = json.Unmarshal(bs, j)
	if err != nil {
		return "", nil
	}

	return j.AccessToken, nil
}

func checkScope(tokenString string, scopes ...string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return false
	}

	hasScope := true
	for _, v := range scopes {
		if !strings.Contains(claims.Scope, v) {
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
