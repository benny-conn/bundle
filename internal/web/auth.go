package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

// type CustomClaims struct {
// 	Id string `json:"id"`
// 	jwt.StandardClaims
// }

// func newAuthToken(profile Profile) (string, error) {

// 	secret := os.Getenv("AUTH0_SECRET")

// 	claims := CustomClaims{
// 		Id: profile.Id,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
// 			Issuer:    "bundle",
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedToken, err := token.SignedString([]byte(secret))
// 	if err != nil {
// 		return "", err
// 	}
// 	return signedToken, nil
// }

// func getProfileFromToken(tokenString string) (Profile, error) {

// 	secret := os.Getenv("AUTH0_SECRET")
// 	token, err := jwt.ParseWithClaims(
// 		tokenString,
// 		&CustomClaims{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			return []byte(secret), nil
// 		},
// 	)
// 	if err != nil {
// 		return Profile{}, err
// 	}

// 	claims, ok := token.Claims.(*CustomClaims)
// 	if !ok {
// 		err = errors.New("couldn't parse claims")
// 		return Profile{}, err
// 	}

// 	gs := gate.NewGateService("", "")
// 	user, err := gs.GetUser(&api.User{Id: claims.Id})
// 	if err != nil {
// 		return Profile{}, err
// 	}

// 	return userToProfile(user), nil

// }

// func validateToken(tokenString string) error {
// 	secret := os.Getenv("AUTH0_SECRET")

// 	token, err := jwt.ParseWithClaims(
// 		tokenString,
// 		&CustomClaims{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			return []byte(secret), nil
// 		},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	claims, ok := token.Claims.(*CustomClaims)
// 	if !ok {
// 		err = errors.New("couldn't parse claims")
// 		return err
// 	}
// 	if claims.ExpiresAt < time.Now().Local().Unix() {
// 		err = errors.New("expired token")
// 		return err
// 	}

// 	return nil
// }

func newAccessCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   900,
		Secure:   true,
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

		ses := &api.Session{
			Id: tokenString,
		}

		err = checkSession(ses)

		if err != nil {
			fmt.Println(err)
			token.MaxAge = -1
			gs := gate.NewGateService("", "")
			gs.DeleteSession(ses)
			http.SetCookie(w, token)
		}

		next.ServeHTTP(w, req)
	})
}

func newSession(prof Profile) (*api.Session, error) {
	gs := gate.NewGateService("", "")

	req := &api.Session{
		UserId: prof.Id,
	}

	err := gs.InsertSession(req)
	if err != nil {
		return nil, err
	}
	ses, err := gs.GetSession(req)
	if err != nil {
		return nil, err
	}

	return ses, nil
}

func checkSession(req *api.Session) error {
	gs := gate.NewGateService("", "")

	ses, err := gs.GetSession(req)
	if err != nil {
		return err
	}

	if time.Since(time.Unix(ses.LastRetrieved, 0)) > (15 * time.Minute) {
		return errors.New("idle time exceeded")
	}
	return nil
}

func getProfileFromToken(token string) (Profile, error) {
	req := &api.Session{
		Id: token,
	}
	gs := gate.NewGateService("", "")

	ses, err := gs.GetSession(req)
	if err != nil {
		return Profile{}, nil
	}
	user, err := gs.GetUser(&api.User{Id: ses.UserId})
	if err != nil {
		return Profile{}, nil
	}
	return userToProfile(user), nil
}
