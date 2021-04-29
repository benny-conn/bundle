package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

func UsersHandlerFunc(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		req.ParseForm()

		userName := req.FormValue("username")
		email := req.FormValue("email")

		r := &api.GetUserRequest{
			Username: userName,
			Email:    email,
		}
		user, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bs, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(bs)
	case http.MethodPost:
		bs, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newUser := &api.User{}
		err = json.Unmarshal(bs, newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = insertUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		bs, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newUser := &api.UpdateUserRequest{}
		err = json.Unmarshal(bs, newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = updateUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func insertUser(user *api.User) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	_, err = client.InsertUser(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

func getUser(req *api.GetUserRequest) (*api.User, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	user, err := client.GetUser(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func updateUser(req *api.UpdateUserRequest) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	_, err = client.UpdateUser(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}
