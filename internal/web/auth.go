package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

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

func newSession(prof profile) (*api.Session, error) {
	gs := gate.NewGateService("", "")

	req := &api.Session{
		UserId: prof.Id,
	}

	res, err := gs.InsertSession(req)
	if err != nil {
		return nil, err
	}
	req.Id = res.Id
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

func getProfileFromToken(token string) (profile, error) {
	req := &api.Session{
		Id: token,
	}
	gs := gate.NewGateService("", "")

	ses, err := gs.GetSession(req)
	if err != nil {
		return profile{}, nil
	}
	user, err := gs.GetUser(&api.User{Id: ses.UserId})
	if err != nil {
		return profile{}, nil
	}
	return userToProfile(user), nil
}
