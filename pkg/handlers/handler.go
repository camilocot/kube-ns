package handlers

import (
	"log"

	"github.com/camilocot/kubernetes-ns-default-netpol/config"
	"github.com/camilocot/kubernetes-ns-default-netpol/pkg/event"
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

// Map maps each event handler function to a name for easily lookup
var Map = map[string]interface{}{
	"default": &Default{},
}

// Default handler implements Handler interface,
// print each event with JSON format
type Default struct {
	NetPolRecipe string
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init(c *config.Config) error {
	d.NetPolRecipe = c.NetPol.Recipe
	return nil
}

// Handle handles an event.
func (d *Default) Handle(e event.Event, client kubernetes.Interface) {
	log.Printf("Processing ns %s", e.Name)

	switch d.NetPolRecipe {
	case "deny-all":
		_, err := client.NetworkingV1().NetworkPolicies(e.Name).Create(d.newNetPol(e))
		if err != nil {
			log.Printf("err %s creating netpol", err)
			return
		}
	default:
		log.Printf("err could not get netpol action: %s", d.NetPolRecipe)
	}
}

func (d *Default) newNetPol(e event.Event) *netv1.NetworkPolicy {

	return &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.NetPolRecipe,
			Namespace: e.Name,
		},
	}
}
