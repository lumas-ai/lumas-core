package streamserver

import (
	"context"
	"errors"
	"fmt"

  "google.golang.org/grpc"
	"github.com/google/uuid"
	processor "github.com/lumas-ai/lumas-core/protos/golang/processor"
	camera "github.com/lumas-ai/lumas-core/protos/golang/camera"
	api "github.com/lumas-ai/lumas-core/protos/golang/stream"
	//"github.com/go-redis/redis"
	//"github.com/nitishm/go-rejson"
)

type StreamServer struct {
	sessions []*api.Session
	//redis    *rejson.Handler
}

//func Init(redisServer string, redisPass string) (*StreamServer, error) {
//	redisClient := redis.NewClient(&redis.Options{
//		Addr:     redisServer,
//		Password: redisPass,
//		DB:       0,
//	})
//
//	pong, err := redisClient.Ping().Result()
//	if pong != "PONG" || err != nil {
//		return nil, err
//	}
//
//	redisJSONClient := rejson.NewReJSONHandler()
//	redisJSONClient.SetGoRedisClient(redisClient)
//
//	return &StreamServer{
//		redis: redisJSONClient,
//	}, nil
//}

func cameraClient() (camera.CameraClient, error) {
  var opts []grpc.DialOption

  opts = append(opts, grpc.WithInsecure())
  conn, err := grpc.Dial("camera", opts...)
  if err != nil {
    return nil, err
  }
  client := camera.NewCameraClient(conn)

  return client, nil
}

func processorClient() (processor.ProcessorClient, error) {
  var opts []grpc.DialOption

  opts = append(opts, grpc.WithInsecure())
  conn, err := grpc.Dial("processor", opts...)
  if err != nil {
    return nil, err
  }
  client := processor.NewProcessorClient(conn)

  return client, nil
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

	//res, err := s.redis.JSONDel("session_"+session.Id, ".")
	//if res.(string) != "OK" || err != nil {
	//	return err
	//}

	return nil
}

func (s *StreamServer) setSession(session *api.Session) error {
	//res, err := s.redis.JSONSet("session_"+session.Id, ".", session)
	//if res == nil || err != nil {
	//	errMsg := fmt.Sprintf("Could not set session to Redis: %v", err.Error())
	//	log.Println(errMsg)
	//	return err
	//}

	return nil
}

func (s *StreamServer) getSession(id string) (*api.Session, bool) {
	//var sessionJSON interface{}
	session := &api.Session{}

	//if sessionJSON, err = s.redis.JSONGet("session_"+id, "."); err != nil {
	//	errMsg := fmt.Sprintf("Unable to retrieve session from Redis: %s", err.Error())
	//	log.Println(errMsg)
	//	return session, false
	//}

	//if sessionJSON == nil {
	//	return session, false
	//}

	//if err := json.Unmarshal(sessionJSON.([]byte), &session); err != nil {
	//	errMsg := fmt.Sprintf("Could not retrieve session from Redis: %v", err.Error())
	//	log.Println(errMsg)
	//	return session, false
	//}

	//source := getSource(session.SourceID)

	return session, true
}

func (s *StreamServer) newSession(source *api.Source) (*api.Session, error) {
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

  return session, nil
}

func (s *StreamServer) Setup(ctx context.Context, source *api.Source) (*api.Session, error) {
	sess, err := s.newSession(source)
  if err != nil {
		emptySession := &api.Session{}
		return emptySession, err
	}

	return sess, nil
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

func (s *StreamServer) Record(ctx context.Context, source *api.Source) (*api.RecordSession, error) {
  //Create a new session
	sess, err := s.newSession(source)
  if err != nil {
		r := &api.RecordSession{}
		return r, err
	}

  //Get SDP information about the camera
  cameraClient, err := cameraClient()
  if err == nil {
    m := fmt.Sprintf("Could not get camera information: %v", err)
    r := &api.RecordSession{}
    return r, errors.New(m)
  }
  cID := &camera.CameraID{ Id: source.Id}
  cameraInfo, err := cameraClient.Describe(context.Background(), cID)
  if err != nil {
    m := fmt.Sprintf("Could not get camera information: %v", err)
    r := &api.RecordSession{}
    return r, errors.New(m)
  }
  if cameraInfo.AudioSDP == "" && cameraInfo.VideoSDP == "" {
    m := fmt.Sprintf("Camera %s has no video SDP information", source.Id)
    r := &api.RecordSession{}
    return r, errors.New(m)
  }

  //Request a processor to handle the recording feed
  processorSession := &processor.Session{
    Id: sess.Id,
    VideoSDP: cameraInfo.VideoSDP,
    AudioSDP: cameraInfo.AudioSDP,
  }

  processorClient, err := processorClient()
  rtpConfig, err := processorClient.Setup(context.Background(), processorSession)
  if err == nil {
    m := fmt.Sprintf("Could not setup feed processor: %v", err)
    r := &api.RecordSession{}
    return r, errors.New(m)
  }

  //We need to translate the processor.RTPConfig to api.RTPConfig
  sRTPConfig := &api.RTPConfig{
    RtpAddress: rtpConfig.RtpAddress,
    AudioRTPPort: rtpConfig.AudioRTPPort,
    VideoRTPPort: rtpConfig.VideoRTPPort,
  }

  r := &api.RecordSession{
    Id: sess.Id,
    RtpConfig: sRTPConfig,
  }

	return r, nil
}
