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
func New(obj interface{}, action string) (kbEvent Event, err error) {

	metadata, err := utils.GetObjectMetaData(obj)

	if err != nil {
		return kbEvent, err
	}

	kbEvent = Event{
		Reason:      action,
		Name:        metadata.Name,
		Annotations: metadata.Annotations,
	}
	return kbEvent, err
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
