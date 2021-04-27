package routes

import "net/http"

func LogoutHandlerFunc(w http.ResponseWriter, req *http.Request) {

	accessCookie, err := req.Cookie("access_token")
	if err == nil {
		accessCookie.MaxAge = -1
	}

	http.SetCookie(w, accessCookie)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}
