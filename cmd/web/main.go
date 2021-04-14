package web

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/internal/web"
)

func main() {
	fmt.Println("Started server on port 8080")

	http.HandleFunc("/", web.HandleRoot)
	http.HandleFunc("/signup", web.HandleSignup)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("../../assets/public"))))

	http.ListenAndServe(":8080", nil)
}
