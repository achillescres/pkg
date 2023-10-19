package messagebroker

import "context"

type Message interface{}

type Callback[MessageType Message] func(MessageType)
type CancelSubscription func()

type Topic[MessageType Message] interface {
	SubTopic[MessageType]
	PubTopic[MessageType]
	//Pub(ctx context.Context, message MessageType) error
	//Sub(callback Callback[MessageType]) (CancelSubscription, error)
}

type SubTopic[MessageType Message] interface {
	Name() string
	Sub(callback Callback[MessageType]) (CancelSubscription, error)
}

type PubTopic[MessageType Message] interface {
	Name() string
	Pub(ctx context.Context, message MessageType) error
}
