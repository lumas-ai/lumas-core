package events

import (
  "fmt"
  "log"
  "context"
  "errors"

  _struct "github.com/golang/protobuf/ptypes/struct"
  api "github.com/lumas-ai/lumas-core/protos/golang"
)

type frame struct {
  imageFormat string
  base64Image string
  presentationTimeStamp int
}

type video struct {
  startPTS int
  endPTS int
  source string
}

type artifact struct {
  name string
  video *video
  frame *frame
  parameters *_struct.Struct
}

type action struct {
  name string
  subscribers map[string]*subscriber //map[service name]*event
  parameters *_struct.Struct
  artifacts []*artifact
}

type camera struct {
  id int
  name string
  events map[string]*event //map[event type]*event
}

type event struct {
  eventType string
  artifact *artifact
  actions map[string]*action //map[action name]*action
}

type subscriber struct {
  id int
  name string
  channel chan *event
}

type EventsServer struct {
  cameras map[int]*camera
  subscribers map[string]*subscriber //map[service name]*subscriber
}

func getArtifact(event *api.Event) (*artifact, error) {
  art := &artifact{
    frame: &frame{
      imageFormat: "jpeg",
      base64Image: "oijweofijawoiefj",
    },
  }

  return art, nil
}

func (a *artifact) toAPIArtifact() *api.Artifact {
  var f *api.Frame
  var v *api.Video

  if a.video != nil {
    v = &api.Video {
      StartPTS: int64(a.video.startPTS),
      EndPTS: int64(a.video.endPTS),
      Source: a.video.source,
    }
  }

  if a.frame != nil {
    f = &api.Frame{
      Image: &api.Image {
        Base64Image: a.frame.base64Image,
        Format: a.frame.imageFormat,
      },
      PresentationTimeStamp: int64(a.frame.presentationTimeStamp),
    }
  }

  r := &api.Artifact{
    Name: a.name,
    Frame: f,
    Video: v,
    Parameters: a.parameters,
  }

  return r
}


func (a *artifact) fromAPIArtifact(apiArtifact *api.Artifact) {
  if apiArtifact.Frame != nil {
    a.frame = &frame{
      base64Image: apiArtifact.Frame.Image.Base64Image,
      imageFormat: apiArtifact.Frame.Image.Format,
      presentationTimeStamp: int(apiArtifact.Frame.PresentationTimeStamp),
    }
  }

  if apiArtifact.Video != nil {
    a.video = &video{
      startPTS: int(apiArtifact.Video.StartPTS),
      endPTS: int(apiArtifact.Video.EndPTS),
      source: apiArtifact.Video.Source,
    }
  }
}

func (e *EventsServer) RegisterAction(ctx context.Context, ea *api.EventAction) (*api.Result, error) {
  t_event := ea.Event
  t_action := ea.Action
  t_camera := ea.Camera
  t_service := ea.Receiver

  var event_artifacts []*artifact
  for _, _artifact := range t_action.Artifact {
    a := &artifact{}
    a.fromAPIArtifact(_artifact)
    event_artifacts = append(event_artifacts, a)
  }

  for _, cam := range t_camera {

    //Create the camera hash if one doesn't exist
    if e.cameras[int(cam.Id)] == nil {
      e.cameras[int(cam.Id)] = &camera{
        name: cam.Name,
        id: int(cam.Id),
      }
    }

    camera := e.cameras[int(cam.Id)]

    //Create the events hash if it doesn't exist
    if camera.events[t_event.Type] == nil {
      camera.events[t_event.Type] = &event{
        eventType: t_event.Type,
      }
    }

    camera_event := camera.events[t_event.Type]

    //Create the action hash if it doesn't exist
    if camera_event.actions[t_action.Name] == nil {
      camera_event.actions[t_action.Name] = &action{
        name: t_action.Name,
        parameters: t_action.Parameters,
      }
    }

    event_action := camera_event.actions[t_action.Name]

    //Create the subsriber hash if it doesn't exist
    if event_action.subscribers[t_service.Name] == nil {
      //Check to see if the service has already subscribed and point
      //to it if it has
      if e.subscribers[t_service.Name] != nil {
        event_action.subscribers[t_service.Name] = e.subscribers[t_service.Name]
      } else {
        s := &subscriber{
          id: int(t_service.Id),
          name: t_service.Name,
        }

        event_action.subscribers[t_service.Name] = s

        //The subscriber also needs to be referenced in a global subscriber
        //list so that when the subscriber subscribes for events, we'll be able to
        //associate the subscriber with its configured camera and event types
        //This will need to be secured in the future using cert auth so that one service
        //cannot subscribe to another service's events by stealing its name
        e.subscribers[t_service.Name] = s
      }
    }

  }

  r := &api.Result{
    Successful: true,
  }

  return r, nil
}

func (e *EventsServer) Subscribe(newSubscriber *api.Subscriber, stream api.Events_SubscribeServer) error {

  if e.subscribers[newSubscriber.Name] == nil {
    s := &subscriber {
      channel: make(chan *event),
      name:    newSubscriber.Name,
    }

    e.subscribers[newSubscriber.Name] = s
  } else {
    e.subscribers[newSubscriber.Name].channel = make(chan *event)
  }

  go func() {
    for s_event := range e.subscribers[newSubscriber.Name].channel {
      e := &api.Event{
        Type: s_event.eventType,
        Artifact: s_event.artifact.toAPIArtifact(),
      }

      stream.Send(e)
    }
  }()

  return nil
}

func (e *EventsServer) Publish(ctx context.Context, ev *api.Event) (*api.Result, error) {
  cam := e.cameras[int(ev.Camera.Id)]

  //a := &artifact{}
  //a.fromAPIArtifact(ev.Artifact)
  //eventArtifact := a

  if cam == nil {
    msg := fmt.Sprintf("Ignoring unkonwn event %s for unknown camera %d", ev.Type, ev.Camera.Id)
    log.Println(msg)
    r := &api.Result{
      Successful: false,
      ErrorKind: "UnknownCamera",
      Message: msg,
    }
    return r, errors.New(msg)
  }

  if cam.events[ev.Type] == nil {
    msg := fmt.Sprintf("Ignoring unkonwn event type %s", ev.Type)
    log.Println(msg)
    r := &api.Result{
      Successful: false,
      ErrorKind: "UnknownEventType",
      Message: msg,
    }
    return r, errors.New(msg)
  }

  cam_event := cam.events[ev.Type]

  for _, eventAction := range cam_event.actions {
    for _, subscriber := range eventAction.subscribers {
      art, _ := getArtifact(ev)
      sendEvent := &event{
        eventType: cam_event.eventType,
        artifact: art,
      }
      subscriber.channel <- sendEvent
    }
  }

  r := &api.Result{
    Successful: true,
  }

  return r, nil
}
