package sources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	api "github.com/lumas-ai/lumas-core/protos/golang/source"
	"github.com/nitishm/go-rejson"
)

type SourceServer struct {
	redis *rejson.Handler
}

type sourceType int

const (
	EVENT  sourceType = 0
	CAMERA sourceType = 1
)

type Source struct {
	Id            string
	SourceType    string
	HasLiveStream bool
	CameraID      string
	EventID       string
}

func Init(redisServer string, redisPass string) (*SourceServer, error) {
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

	//If the root source entry hasn't been set, we need to initialize it in Redis
	res, err := redisJSONClient.JSONGet("source", ".")
	if res == nil {
		aa := make(map[string]interface{})
		if _, err := redisJSONClient.JSONSet("source", ".", aa); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return &SourceServer{
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

func (s *SourceServer) delSource(sourceID string) error {
	if _, exists := s.getSource(sourceID); !exists {
		errMsg := fmt.Sprintf("Could not find source with id %s", sourceID)
		return errors.New(errMsg)
	}

	res, err := s.redis.JSONDel("source", ".id_"+sourceID)
	if res.(string) != "OK" || err != nil {
		return err
	}

	return nil
}

func (s *SourceServer) setSource(sourceInfo *api.SourceInfo) error {
	rsource := &Source{
		Id:         sourceInfo.Id,
		SourceType: "CAMERA",
		CameraID:   sourceInfo.CameraID,
		EventID:    sourceInfo.EventID,
	}

	res, err := s.redis.JSONSet("source", ".id_"+rsource.Id, rsource)
	if res == nil || err != nil {
		errMsg := fmt.Sprintf("Could not set source to Redis: %v", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

func (s *SourceServer) getSource(id string) (*Source, bool) {
	var sourceJSON interface{}
	var err error
	rsource := &Source{}

	if sourceJSON, err = s.redis.JSONGet("source", ".id_"+id); err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve source from Redis: %s", err.Error())
		log.Println(errMsg)
		return rsource, false
	}

	if sourceJSON == nil {
		return rsource, false
	}

	if err := json.Unmarshal(sourceJSON.([]byte), &rsource); err != nil {
		errMsg := fmt.Sprintf("Could not retrieve source from Redis: %v", err.Error())
		log.Println(errMsg)
		return rsource, false
	}

	return rsource, true
}

func (s *SourceServer) Add(ctx context.Context, source *api.SourceInfo) (*api.SourceID, error) {
	return &api.SourceID{Id: "12345"}, nil
}

func (s *SourceServer) List(listRequest *api.SourceListRequest, stream api.Source_ListServer) error {
	//var sourceJSON interface{}
	if _, err := s.redis.JSONGet("source", "."); err != nil {
		return err
	}

	return nil
}

func (s *SourceServer) Update(ctx context.Context, sourceInfo *api.SourceInfo) (*api.Result, error) {
	return result("", ""), nil
}

func (s *SourceServer) Delete(ctx context.Context, sourceID *api.SourceID) (*api.Result, error) {
	if _, exists := s.getSource(sourceID.Id); !exists {
		errMsg := fmt.Sprintf("Could not find source with ID %s", sourceID.Id)
		return result("SourceNotFound", errMsg), errors.New(errMsg)
	}

	if err := s.delSource(sourceID.Id); err != nil {
		return result("SourceNotDeleted", err.Error()), err
	}

	return result("", ""), nil
}

func (s *SourceServer) Describe(ctx context.Context, sourceID *api.SourceID) (*api.SourceInfo, error) {
	return &api.SourceInfo{Id: "12345"}, nil
}
