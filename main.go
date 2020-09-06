package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type flags struct {
	Kubeconfig    string
	ResyncPeriodS string
	ResyncPeriod  time.Duration
	StatusAddr    string
}

var f flags

func init() {
	var err error
	flag.StringVar(&f.Kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.StringVar(&f.ResyncPeriodS, "resync-period", "30m", "resynchronization period")
	flag.Parse()

	f.ResyncPeriod, err = time.ParseDuration(f.ResyncPeriodS)
	if err != nil {
		panic(err)
	}
}

func main() {
	var config *rest.Config
	var err error
	var client kubernetes.Interface

	if f.Kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", f.Kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", f.Kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	client = kubernetes.NewForConfigOrDie(config)

	handler := NewNSNetPolHandler(client, f.ResyncPeriod)
	log.Printf("using in-cluster configuration")

	go handler.Run()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type nsNetPolHandler struct {
	client     kubernetes.Interface
	store      cache.Store
	controller cache.Controller
}

type NSNetPolHandler interface {
	Run()
	Synced() bool
}

const (
	NetPolAnnotation = "netpol.v1.camilocot/action"
)

func NewNSNetPolHandler(client kubernetes.Interface, resyncPeriod time.Duration) NSNetPolHandler {

	handler := nsNetPolHandler{
		client: client,
	}

	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Namespaces().List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Namespaces().Watch(lo)
			},
		},
		&v1.Namespace{},
		resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: handler.NamespaceAdded,
		},
	)

	handler.store = store
	handler.controller = controller

	return &handler
}

func (r *nsNetPolHandler) Synced() bool {
	return r.controller.HasSynced()
}

func (r *nsNetPolHandler) Run() {
	log.Printf("running namespace controller")
	r.controller.Run(wait.NeverStop)
}

func (r *nsNetPolHandler) NamespaceAdded(obj interface{}) {
	ns := obj.(*v1.Namespace)
	log.Printf("Added ns: %s", ns.Name)

	val, ok := ns.Annotations[NetPolAnnotation]
	if !ok {
		log.Printf("ns %s no default netpol", ns.Name)
		return
	}

	log.Printf("ns %s default netpol %s", ns.Name, val)

	switch val {
	case "deny-all":
		_, err := r.client.NetworkingV1().NetworkPolicies(ns.Name).Create(netNetPol(ns.Name))
		if err != nil {
			log.Printf("err %s default netpol", err)
			return
		}
	default:
		log.Printf("could not get netpol action: %s", val)
	}

}

func netNetPol(nsName string) *networkingv1.NetworkPolicy {

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deny-all",
			Namespace: nsName,
		},
	}
}
