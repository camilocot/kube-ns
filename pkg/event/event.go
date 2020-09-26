package event

import (
	"fmt"

	"github.com/bitnami-labs/kubewatch/pkg/utils"
)

// Event represent an event got from k8s api server
// Events from different endpoints need to be casted to KubewatchEvent
// before being able to be handled by handler
type Event struct {
	Reason string
	Name   string
}

// New create new KubewatchEvent
func New(obj interface{}, action string) Event {
	var reason, name string

	objectMeta := utils.GetObjectMetaData(obj)
	name = objectMeta.Name
	reason = action

	kbEvent := Event{
		Reason: reason,
		Name:   name,
	}
	return kbEvent
}

// Message returns event message in standard format.
// included as a part of event package to enhance code reusability across handlers.
func (e *Event) Message() (msg string) {
	msg = fmt.Sprintf(
		"A namespace named `%s` has been `%s`",
		e.Name,
		e.Reason,
	)
	return msg
}
