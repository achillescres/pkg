package kafkapubsub

// Topic is kafka publish and subscribe topic implementation.
// It uses kafka.Writer and kafka.Reader from kafka-go
// Specify your own MessageType with Message interface
// WARNING: barely useless and maybe harmful for you application logic
type Topic[MessageType Message] struct {
	PubTopic[MessageType]
	SubTopic[MessageType]
}

func NewTopic[MessageType Message](pubTopic PubTopic[MessageType], subTopic SubTopic[MessageType]) *Topic[MessageType] {
	return &Topic[MessageType]{PubTopic: pubTopic, SubTopic: subTopic}
}
