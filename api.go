package main

import (
	"context"
	"errors"

	"github.com/mattn/go-mastodon"
)

type TimelineType uint

const (
	TimelineHome TimelineType = iota
	TimelineDirect
	TimelineLocal
	TimelineFederated
)

type API struct {
	Client *mastodon.Client
}

type AccountRegister struct {
	Account
	AuthURI string
}

func (api *API) SetClient(c *mastodon.Client) {
	api.Client = c
}

func (api *API) getStatuses(tl TimelineType, pg *mastodon.Pagination) ([]*mastodon.Status, error) {
	var statuses []*mastodon.Status
	var err error

	switch tl {
	case TimelineHome:
		statuses, err = api.Client.GetTimelineHome(context.Background(), pg)
	case TimelineDirect:
		statuses, err = api.Client.GetTimelineDirect(context.Background(), pg)
	case TimelineLocal:
		statuses, err = api.Client.GetTimelinePublic(context.Background(), true, pg)
	case TimelineFederated:
		statuses, err = api.Client.GetTimelinePublic(context.Background(), false, pg)
	default:
		err = errors.New("No timeline selected")
	}

	return statuses, err
}

func (api *API) GetStatuses(tl TimelineType) ([]*mastodon.Status, error) {
	return api.getStatuses(tl, nil)
}

func (api *API) GetStatusesOlder(tl TimelineType, s *mastodon.Status) ([]*mastodon.Status, bool, error) {
	pg := &mastodon.Pagination{
		MaxID: s.ID,
	}

	statuses, err := api.getStatuses(tl, pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MinID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) GetStatusesNewer(tl TimelineType, s *mastodon.Status) ([]*mastodon.Status, bool, error) {
	pg := &mastodon.Pagination{
		MinID: s.ID,
	}

	statuses, err := api.getStatuses(tl, pg)
	if err != nil {
		return statuses, false, err
	}

	if pg.MaxID == "" {
		return statuses, false, err
	}

	return statuses, true, err
}

func (api *API) GetThread(s *mastodon.Status) ([]*mastodon.Status, int, error) {
	cont, err := api.Client.GetStatusContext(context.Background(), s.ID)
	if err != nil {
		return nil, 0, err
	}
	thread := cont.Ancestors
	thread = append(thread, s)
	thread = append(thread, cont.Descendants...)
	return thread, len(cont.Ancestors), nil
}

func (api *API) Boost(s *mastodon.Status) error {
	_, err := api.Client.Reblog(context.Background(), s.ID)
	return err
}

func (api *API) Unboost(s *mastodon.Status) error {
	_, err := api.Client.Unreblog(context.Background(), s.ID)
	return err
}

func (api *API) Favorite(s *mastodon.Status) error {
	_, err := api.Client.Favourite(context.Background(), s.ID)
	return err
}

func (api *API) Unfavorite(s *mastodon.Status) error {
	_, err := api.Client.Unfavourite(context.Background(), s.ID)
	return err
}

func (api *API) DeleteStatus(s *mastodon.Status) error {
	//TODO: check user here?
	return api.Client.DeleteStatus(context.Background(), s.ID)
}
