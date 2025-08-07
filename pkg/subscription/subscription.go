package subscription

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type Isubscription interface {
	Subscribe() (watch.Interface, error)
	Reconcile(object runtime.Object, event watch.EventType)
	IsComplete() <-chan bool
}
