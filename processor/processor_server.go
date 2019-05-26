package processor

import (
  "fmt"
  "context"
  "errors"
  api "github.com/lumas-ai/lumas-core/protos/golang/processor"
)

type ProcessorServer struct {
  feeds []*Feed
}

func (s *ProcessorServer) FindFeed(id string) (*Feed, int ) {
  for i, feed := range s.feeds {
    if feed.Id == id {
      return feed, i
    }
  }

  return nil, -1
}

func (s *ProcessorServer) Setup(ctx context.Context, session *api.Session) (*api.RTPConfig, error) {
  //Feeds are created with the same ID as the calling session so their IDs should always match
  feed, _ := s.FindFeed(session.Id)
  if feed != nil {
    m := fmt.Sprintf("Feed with ID %s is already being processed." + feed.Id)
    return feed.RTPConfig, errors.New(m)
  }

  feed, err := NewFeed(session.Id, session.VideoSDP, session.AudioSDP)
  if err != nil {
    r := &api.RTPConfig{}
    return r, err
  }

  if feed.RTPConfig == nil {
    r := &api.RTPConfig{}
    return r, err
  }

  return feed.RTPConfig, nil
}

func (s *ProcessorServer) Teardown(ctx context.Context, session *api.Session) (*api.Result, error) {
  feed, _ := s.FindFeed(session.Id)
  if feed != nil {
    m := fmt.Sprintf("Feed with ID %s not found." + feed.Id)
    r := api.Result{Successful: false, ErrorKind: "FeedNotFound", Message: m}
    return &r, errors.New(m)
  }

  return &api.Result{Successful: true}, nil
}
