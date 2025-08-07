package runtime

import (
	"k8soperator/pkg/subscription"
	"sync"
)

// for {
// 	select {
// 	case msg := <-wiface.ResultChan():
// 		klog.Infof("%+v", msg)
// 	}
// }

var wg sync.WaitGroup

func Runloop(subscription []subscription.Isubscription) error {
	for _, subsubscription := range subscription {
		wiface, err := subsubscription.Subscribe()
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case msg := <-wiface.ResultChan():
					subsubscription.Reconcile(msg.Object, msg.Type)
					//TODO: want a way to escape
				}
			}
		}()
	}
	for _, subscription := range subscription {
		select {
		case _ = <-subscription.IsComplete():
			break
		}
	}
	return nil
}
