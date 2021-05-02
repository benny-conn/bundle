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
	"github.com/bennycio/bundle/internal/gate"
)

func UpdateUserApi(updatedUser *api.User) error {
	port := os.Getenv("GATE_PORT")
	host := os.Getenv("GATE_HOST")
	addr := fmt.Sprintf("http://%v:%v/api/users", host, port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	updatedBs, err := json.Marshal(updatedUser)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.Write(updatedBs)

	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}

	access, err := gate.GetAccessToken()
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
	port := os.Getenv("GATE_PORT")
	host := os.Getenv("GATE_HOST")
	addr := fmt.Sprintf("http://%v:%v/api/users", host, port)
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
	access, err := gate.GetAccessToken()
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
	port := os.Getenv("GATE_PORT")
	host := os.Getenv("GATE_HOST")
	addr := fmt.Sprintf("http://%v:%v/api/users", host, port)
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
	access, err := gate.GetAccessToken()
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
