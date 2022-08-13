package api

import (
	"context"
)

func (ac *AccountClient) ClearNotifications() error {
	return ac.Client.ClearNotifications(context.Background())
}
