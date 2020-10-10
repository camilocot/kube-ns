package handlers

import (
	"fmt"
	"log"

	"github.com/camilocot/kube-ns/config"
	"github.com/camilocot/kube-ns/pkg/event"
	"github.com/camilocot/kube-ns/pkg/utils"
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

// DenyAllNetpol defines the deny all netpol name
const DenyAllNetpol = "deny-all"

// DenyIngressNetpol defines the deny all netpol name
const DenyIngressNetpol = "deny-ingress"

// Init initializes handler configuration
func (h *HandlerConfig) Init(c *config.Config) error {
	h.NetPol = c.NetPol
	return nil
}

// Handle handles an event.
func (h *HandlerConfig) Handle(e event.Event, client kubernetes.Interface) {
	log.Printf("Processing ns %s, reason: %s", e.Name, e.Reason)

	if h.NetPol.Enabled {
		switch e.Annotations[h.NetPol.Annotation] {
		case DenyAllNetpol:
			_, err := client.NetworkingV1().NetworkPolicies(e.Name).Create(denyAllNetPol(e))
			if err != nil {
				log.Printf("err %s creating deny-all netpol", err)
				return
			}
			err = DeletePreviousNetpols(e, client, DenyAllNetpol)
		case DenyIngressNetpol:
			_, err := client.NetworkingV1().NetworkPolicies(e.Name).Create(denyIngressNetPol(e))
			if err != nil {
				log.Printf("err %s creating deny-ingress netpol", err)
				return
			}
			err = DeletePreviousNetpols(e, client, DenyIngressNetpol)

		default:
			log.Printf("Cant find annotation: %s in %v", h.NetPol.Annotation, e.Annotations)
		}
	}
}

func denyAllNetPol(e event.Event) *netv1.NetworkPolicy {

	return &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DenyAllNetpol,
			Namespace: e.Name,
			Labels:    map[string]string{"kubens/managed.netpod": "true"},
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []netv1.PolicyType{"Ingress", "Egress"},
		},
	}
}

func denyIngressNetPol(e event.Event) *netv1.NetworkPolicy {

	return &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DenyIngressNetpol,
			Namespace: e.Name,
			Labels:    map[string]string{"kubens/managed.netpod": "true"},
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []netv1.PolicyType{"Ingress"},
		},
	}
}

// DeletePreviousNetpols deletes previous applied network policies
func DeletePreviousNetpols(e event.Event, client kubernetes.Interface, current string) (err error) {
	netpols, err := client.NetworkingV1().NetworkPolicies(e.Name).List(metav1.ListOptions{LabelSelector: "kubens/managed.netpod=true"})
	if err != nil {
		return fmt.Errorf("err %s getting managed netpols", err)
	}

	for _, netpol := range netpols.Items {

		metadata, err := utils.GetObjectMetaData(&netpol)

		if err != nil {
			return fmt.Errorf("err %s getting metadata info", err)
		}

		if metadata.Name != current {
			log.Printf("deleted %s previous netpol %s", metadata.Name, e.Name)
			client.NetworkingV1().NetworkPolicies(e.Name).Delete(metadata.Name, &metav1.DeleteOptions{})
		}

	}
	return err
}
