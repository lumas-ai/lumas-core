package streamserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	api "github.com/lumas-ai/lumas-core/protos/golang/stream"
	"github.com/nitishm/go-rejson"
)

type StreamServer struct {
	sessions []*api.Session
	redis    *rejson.Handler
}

func Init(redisServer string, redisPass string) (*StreamServer, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisServer,
		Password: redisPass,
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	if pong != "PONG" || err != nil {
		return nil, err
	}

	redisJSONClient := rejson.NewReJSONHandler()
	redisJSONClient.SetGoRedisClient(redisClient)

	return &StreamServer{
		redis: redisJSONClient,
	}, nil
}

func result(errorKind string, errorMessage string) *api.Result {
	if errorKind == "" && errorMessage == "" {
		return &api.Result{
			Successful: true,
		}
	}
	return &api.Result{
		Successful: false,
		ErrorKind:  errorKind,
		Message:    errorMessage,
	}
}

func (s *StreamServer) delSession(session *api.Session) error {
	if _, exists := s.getSession(session.Id); !exists {
		errMsg := fmt.Sprintf("Could not find session with id %s", session.Id)
		return errors.New(errMsg)
	}

	res, err := s.redis.JSONDel("session_"+session.Id, ".")
	if res.(string) != "OK" || err != nil {
		return err
	}

	return nil
}

func (s *StreamServer) setSession(session *api.Session) error {
	res, err := s.redis.JSONSet("session_"+session.Id, ".", session)
	if res == nil || err != nil {
		errMsg := fmt.Sprintf("Could not set session to Redis: %v", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

func (s *StreamServer) getSession(id string) (*api.Session, bool) {
	var sessionJSON interface{}
	var err error
	session := &api.Session{}

	if sessionJSON, err = s.redis.JSONGet("session_"+id, "."); err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve session from Redis: %s", err.Error())
		log.Println(errMsg)
		return session, false
	}

	if sessionJSON == nil {
		return session, false
	}

	if err := json.Unmarshal(sessionJSON.([]byte), &session); err != nil {
		errMsg := fmt.Sprintf("Could not retrieve session from Redis: %v", err.Error())
		log.Println(errMsg)
		return session, false
	}

	//source := getSource(session.SourceID)

	return session, true
}

func (s *StreamServer) Setup(ctx context.Context, source *api.Source) (*api.Session, error) {
	var session *api.Session

	//Generate random IDs until we find one that isn't being used
	//It'll be very rare the first one doesn't work
	for {
		sessionID := uuid.New().String()
		if _, exists := s.getSession(sessionID); exists {
			continue
		} else {
			session = &api.Session{
				Id:     sessionID,
				Source: source,
			}
			break
		}
	}

	if err := s.setSession(session); err != nil {
		emptySession := &api.Session{}
		return emptySession, err
	}

	return session, nil
}

func (s *StreamServer) Play(ctx context.Context, source *api.Session) (*api.Result, error) {
	return &api.Result{Successful: true}, nil
}

func (s *StreamServer) Pause(ctx context.Context, session *api.Session) (*api.Result, error) {
	return &api.Result{Successful: true}, nil
}

func (s *StreamServer) Describe(ctx context.Context, source *api.Source) (*api.StreamSourceInfo, error) {
	return &api.StreamSourceInfo{Id: "12345"}, nil
}

func (s *StreamServer) Teardown(ctx context.Context, session *api.Session) (*api.Result, error) {
	if _, exists := s.getSession(session.Id); !exists {
		errMsg := fmt.Sprintf("Could not find session with ID %s", session.Id)
		return result("SessionNotFound", errMsg), errors.New(errMsg)
	}

	if err := s.delSession(session); err != nil {
		return result("SessionNotDeleted", err.Error()), err
	}

	return result("", ""), nil
}

func (s *StreamServer) Status(ctx context.Context, source *api.Session) (*api.SessionStatus, error) {
	return &api.SessionStatus{}, nil
}

func (s *StreamServer) Record(ctx context.Context, source *api.Source) (*api.Session, error) {
	session := &api.Session{
		Id: "12345",
		RTPConfig: {
			rtpAddress:   "lumas-processor",
			audioRTPPort: 9000,
			videoRTPPort: 9001,
		},
	}

	return session, nil
}
