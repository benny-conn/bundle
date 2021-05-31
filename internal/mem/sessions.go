package mem

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/logger"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
)

type sessionsServer struct {
	client *redis.Client
	api.UnimplementedSessionServiceServer
}

func (s *sessionsServer) Get(ctx context.Context, req *api.Session) (*api.Session, error) {

	ses := &api.Session{}

	res, err := s.client.Get(ctx, req.Id).Result()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(res), ses)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	copy := &api.Session{
		Id:            ses.Id,
		UserId:        ses.UserId,
		LastRetrieved: time.Now().Unix(),
	}

	bs, err := json.Marshal(copy)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	err = s.client.Set(ctx, copy.Id, string(bs), redis.KeepTTL).Err()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	return ses, nil
}

func (s *sessionsServer) Insert(ctx context.Context, req *api.Session) (*api.SessionInsertResponse, error) {
	if req.Id == "" {
		req.Id = uuid.NewV4().String()
	}
	if req.LastRetrieved == 0 {
		req.LastRetrieved = time.Now().Unix()
	}

	bs, err := json.Marshal(req)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	err = s.client.Set(ctx, req.Id, string(bs), 24*time.Hour).Err()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	return &api.SessionInsertResponse{Id: req.Id}, nil
}

func (s *sessionsServer) Delete(ctx context.Context, req *api.Session) (*api.Empty, error) {

	err := s.client.Del(ctx, req.Id).Err()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}
	return &api.Empty{}, nil
}

func newSessionsServer() *sessionsServer {
	s := &sessionsServer{client: newClient()}
	return s
}
