package handlers

import (
	"log"
	"testing"

	"github.com/camilocot/kube-ns/pkg/event"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeletePreviousNetpols(t *testing.T) {
	e := event.Event{Name: "namespace", Reason: "Created"}
	client := fake.NewSimpleClientset(denyAllNetPol(e))

	DeletePreviousNetpols(e, client, DenyIngressNetpol)

	netpols, err := client.NetworkingV1().NetworkPolicies(e.Name).List(metav1.ListOptions{LabelSelector: "kubens/managed.netpod=true"})
	if err != nil {
		log.Printf("err %s getting managed netpols", err)
		return
	}

	if len(netpols.Items) > 0 {
		t.Fatalf("Unexpected number of network policies: %v", len(netpols.Items))
	}
}
