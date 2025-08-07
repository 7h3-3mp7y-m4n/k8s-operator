package subscription

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type PodSubcription struct {
	WatcherInterface      watch.Interface
	Client                kubernetes.Interface
	Ctx                   context.Context
	CompletionChan        chan bool
	ConfigMapSubscriptRef *ConfigmapSubscirption
}

func (p *PodSubcription) applyConfigMapChanges(pod *v1.Pod, event watch.EventType) {

	if p.ConfigMapSubscriptRef != nil {
		if p.ConfigMapSubscriptRef.PlatformConfigMapAnnotations != nil {
			updatedPod := pod.DeepCopy()
			if updatedPod.Annotations == nil {
				updatedPod.Annotations = make(map[string]string)
			}
			// Loop through and apply
			for _, annotation := range p.ConfigMapSubscriptRef.PlatformConfigMapAnnotations.Annotations {
				updatedPod.Annotations[annotation.Name] = annotation.Value
			}
			// Update the pod
			_, err := p.Client.CoreV1().Pods(pod.Namespace).Update(p.Ctx, updatedPod, metav1.UpdateOptions{})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
func (p *PodSubcription) Reconcile(object runtime.Object, event watch.EventType) {
	pod := object.(*v1.Pod)
	klog.Infof("Pod subscription update ->  %s", pod.Name)
	switch event {
	case watch.Added:
		if _, ok := pod.Annotations["type"]; !ok {
			updatePod := pod.DeepCopy()
			updatePod.Annotations = make(map[string]string)
			updatePod.Annotations["type"] = "Tim"
			_, err := p.Client.CoreV1().Pods(pod.Namespace).Update(p.Ctx, updatePod, metav1.UpdateOptions{})
			if err != nil {
				klog.Error(err)
			}
		}

	case watch.Deleted:
	case watch.Modified:
		if pod.Annotations["Type"] == "Tim" {
			klog.Info("This could be magic beyond just a CRD")
		}
	}
}

func (p *PodSubcription) Subscribe() (watch.Interface, error) {
	var err error
	p.WatcherInterface, err = p.Client.CoreV1().Pods("").Watch(p.Ctx, metav1.ListOptions{})
	if err != nil {
		klog.Fatalf("watch interface  : %s", err.Error())
	}
	return p.WatcherInterface, nil
}

func (p *PodSubcription) IsComplete() <-chan bool {

	return p.CompletionChan
}
