package runtime

import (
	"k8soperator/pkg/subscription"
	"sync"

	"k8s.io/klog/v2"
)

// for {
// 	select {
// 	case msg := <-wiface.ResultChan():
// 		klog.Infof("%+v", msg)
// 	}
// }

var wg sync.WaitGroup

func RunLoop(subscriptions []subscription.Isubscription) error {
	var wg sync.WaitGroup

	for _, sub := range subscriptions {
		wg.Add(1)

		go func(subscription subscription.Isubscription) {
			defer wg.Done()

			wiface, err := subscription.Subscribe()
			if err != nil {
				klog.Error("Subscription error: %v", err)
				return
			}

			for event := range wiface.ResultChan() {
				subscription.Reconcile(event.Object, event.Type)
			}
		}(sub)
	}

	wg.Wait()
	return nil
}
