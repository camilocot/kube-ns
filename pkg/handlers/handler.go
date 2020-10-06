package handlers

import (
	"log"

	"github.com/camilocot/kube-ns/config"
	"github.com/camilocot/kube-ns/pkg/event"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init(c *config.Config) error
	Handle(e event.Event, c kubernetes.Interface)
}

// HandlerConfig hadler configuration
type HandlerConfig struct {
	NetPol config.NetPol
}

// Init initializes handler configuration
func (h *HandlerConfig) Init(c *config.Config) error {
	h.NetPol = c.NetPol
	return nil
}

// Handle handles an event.
func (h *HandlerConfig) Handle(e event.Event, client kubernetes.Interface) {
	log.Printf("Processing ns %s", e.Name)

	if h.NetPol.Enabled {
		switch e.Annotations[h.NetPol.Annotation] {
		case "deny-all":
			_, err := client.NetworkingV1().NetworkPolicies(e.Name).Create(h.newNetPol(e))
			if err != nil {
				log.Printf("err %s creating netpol", err)
				return
			}
		default:
			log.Printf("Cant find annotation: %s in %v", h.NetPol.Annotation, e.Annotations)
		}

	}
}

func (h *HandlerConfig) newNetPol(e event.Event) *netv1.NetworkPolicy {

	return &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      e.Annotations[h.NetPol.Annotation],
			Namespace: e.Name,
		},
	}
}
