package endpoint_poller

import "net/source/userapi"

type EndPointPoll interface {
	Poll(c userapi.IClient) int
}
