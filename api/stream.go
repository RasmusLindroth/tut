package api

import (
	"context"
	"sync"

	"github.com/RasmusLindroth/go-mastodon"
)

type MastodonType uint

const (
	StatusType MastodonType = iota
	UserType
	ProfileType
	NotificationType
	ListsType
)

type StreamType uint

const (
	HomeStream StreamType = iota
	LocalStream
	FederatedStream
	DirectStream
	TagStream
	ListStream
)

type Stream struct {
	id        string
	receivers []*Receiver
	incoming  chan mastodon.Event
	closed    bool
	mux       sync.Mutex
}

type Receiver struct {
	Ch     chan mastodon.Event
	Closed bool
	mux    sync.Mutex
}

func (s *Stream) ID() string {
	return s.id
}

func (s *Stream) AddReceiver() *Receiver {
	ch := make(chan mastodon.Event)
	rec := &Receiver{
		Ch:     ch,
		Closed: false,
	}
	s.receivers = append(s.receivers, rec)
	return rec
}

func (s *Stream) RemoveReceiver(r *Receiver) {
	index := -1
	for i, rec := range s.receivers {
		if rec.Ch == r.Ch {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	s.receivers[index].mux.Lock()
	if !s.receivers[index].Closed {
		close(s.receivers[index].Ch)
		s.receivers[index].Closed = true
	}
	s.receivers[index].mux.Unlock()
	s.receivers = append(s.receivers[:index], s.receivers[index+1:]...)
}

func (s *Stream) listen() {
	for e := range s.incoming {
		switch e.(type) {
		case *mastodon.UpdateEvent, *mastodon.NotificationEvent, *mastodon.DeleteEvent, *mastodon.ErrorEvent:
			for _, r := range s.receivers {
				go func(rec *Receiver) {
					rec.mux.Lock()
					if rec.Closed {
						return
					}
					rec.mux.Unlock()
					rec.Ch <- e
				}(r)
			}
		}
	}
}

func newStream(id string, inc chan mastodon.Event) (*Stream, *Receiver) {
	stream := &Stream{
		id:       id,
		incoming: inc,
	}
	rec := stream.AddReceiver()
	go stream.listen()
	return stream, rec
}

func (ac *AccountClient) NewGenericStream(st StreamType, data string) (rec *Receiver, err error) {
	var id string
	switch st {
	case HomeStream:
		id = "HomeStream"
	case LocalStream:
		id = "LocalStream"
	case FederatedStream:
		id = "FederatedStream"
	case DirectStream:
		id = "DirectStream"
	case TagStream:
		id = "TagStream" + data
	case ListStream:
		id = "ListStream" + data
	default:
		panic("invalid StreamType")
	}
	for _, s := range ac.Streams {
		if s.ID() == id {
			rec = s.AddReceiver()
			return rec, nil
		}
	}
	var ch chan mastodon.Event
	switch st {
	case HomeStream:
		ch, err = ac.Client.StreamingUser(context.Background())
	case LocalStream:
		ch, err = ac.Client.StreamingPublic(context.Background(), true)
	case FederatedStream:
		ch, err = ac.Client.StreamingPublic(context.Background(), false)
	case DirectStream:
		ch, err = ac.Client.StreamingDirect(context.Background())
	case TagStream:
		ch, err = ac.Client.StreamingHashtag(context.Background(), data, false)
	case ListStream:
		ch, err = ac.Client.StreamingList(context.Background(), mastodon.ID(data))
	default:
		panic("invalid StreamType")
	}
	if err != nil {
		return nil, err
	}
	stream, rec := newStream(id, ch)
	ac.Streams[stream.ID()] = stream
	return rec, nil
}

func (ac *AccountClient) NewHomeStream() (*Receiver, error) {
	return ac.NewGenericStream(HomeStream, "")
}

func (ac *AccountClient) NewLocalStream() (*Receiver, error) {
	return ac.NewGenericStream(LocalStream, "")
}

func (ac *AccountClient) NewFederatedStream() (*Receiver, error) {
	return ac.NewGenericStream(FederatedStream, "")
}

func (ac *AccountClient) NewDirectStream() (*Receiver, error) {
	return ac.NewGenericStream(DirectStream, "")
}

func (ac *AccountClient) NewListStream(id mastodon.ID) (*Receiver, error) {
	return ac.NewGenericStream(ListStream, string(id))
}

func (ac *AccountClient) NewTagStream(tag string) (*Receiver, error) {
	return ac.NewGenericStream(TagStream, tag)
}

func (ac *AccountClient) RemoveGenericReceiver(rec *Receiver, st StreamType, data string) {
	var id string
	switch st {
	case HomeStream:
		id = "HomeStream"
	case LocalStream:
		id = "LocalStream"
	case FederatedStream:
		id = "FederatedStream"
	case TagStream:
		id = "TagStream" + data
	case ListStream:
		id = "ListStream" + data
	default:
		panic("invalid StreamType")
	}
	stream, ok := ac.Streams[id]
	if !ok {
		return
	}
	stream.RemoveReceiver(rec)
	stream.mux.Lock()
	if len(stream.receivers) == 0 && !stream.closed {
		stream.closed = true
		delete(ac.Streams, id)
	}
	stream.mux.Unlock()
}

func (ac *AccountClient) RemoveHomeReceiver(rec *Receiver) {
	ac.RemoveGenericReceiver(rec, HomeStream, "")
}

func (ac *AccountClient) RemoveLocalReceiver(rec *Receiver) {
	ac.RemoveGenericReceiver(rec, LocalStream, "")
}

func (ac *AccountClient) RemoveFederatedReceiver(rec *Receiver) {
	ac.RemoveGenericReceiver(rec, FederatedStream, "")
}

func (ac *AccountClient) RemoveListReceiver(rec *Receiver, id mastodon.ID) {
	ac.RemoveGenericReceiver(rec, ListStream, string(id))
}

func (ac *AccountClient) RemoveTagReceiver(rec *Receiver, tag string) {
	ac.RemoveGenericReceiver(rec, TagStream, tag)
}
