package routes

import (
	"context"
	"net/http"

	pb "github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

func UsersHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method == http.MethodPost {

		req.ParseForm()
		newUser := &pb.User{
			Username: req.FormValue("username"),
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		}
		err := insertUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

}

func insertUser(user *pb.User) error {
	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewUsersServiceClient(conn)
	client.InsertUser(context.Background(), user)
	return nil
}
