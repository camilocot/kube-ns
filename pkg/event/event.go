package event

import (
	"fmt"

	"github.com/camilocot/kube-ns/pkg/utils"
)

// Event represent an event got from k8s api server
type Event struct {
	Reason      string
	Name        string
	Annotations map[string]string
}

// New create new KubewatchEvent
func New(obj interface{}, action string) Event {

	kbEvent := Event{
		Reason:      action,
		Name:        utils.GetObjectMetaData(obj).Name,
		Annotations: utils.GetObjectMetaData(obj).Annotations,
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
