package wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
	auth "github.com/bennycio/bundle/internal/auth/client"
)

func UpdateUserApi(username string, updatedUser *api.User) error {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf("http://localhost:%v/api/users", port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	up := &api.UpdateUserRequest{
		Username:    username,
		UpdatedUser: updatedUser,
	}

	updatedBs, err := json.Marshal(up)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.Write(updatedBs)

	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}

	access, err := auth.GetClientToken()
	if err != nil {
		return err
	}
	req.Header.Add("authorization", "Bearer "+access)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func GetUserApi(username string, email string) (*api.User, error) {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf("http://localhost:%v/api/users", port)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("username", username)
	q.Set("email", email)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	access, err := auth.GetClientToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+access)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.User{}

	err = json.Unmarshal(bs, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func InsertUserApi(user *api.User) error {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf("http://localhost:%v/api/users", port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	bs, err := json.Marshal(user)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.Write(bs)

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return err
	}
	access, err := auth.GetClientToken()
	if err != nil {
		return err
	}
	req.Header.Add("authorization", "Bearer "+access)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
