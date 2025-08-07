package subscription

import (
	"context"
	"errors"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

var (
	configmapName      = "default-configmap"
	configmapNamespace = "default"
)

type ConfigmapSubscirption struct {
	WatcherInterface             watch.Interface
	Client                       kubernetes.Interface
	Ctx                          context.Context
	CompletionChan               chan bool
	PlatformConfigMapAnnotations *platformConfig
	PlatformConfigPhase          string
}
type platformAnnotation struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
type platformConfig struct {
	Annotations []platformAnnotation `yaml:"annotations"`
}

func isPlatformConfigmap(configmap *v1.ConfigMap) (bool, error) {
	if configmap == nil {
		return false, errors.New("No/empty configmap here")
	}
	if configmap.Name == configmapName {
		klog.Info("we caught our configmap")
		return true, nil
	}
	return false, nil
}

func (c *ConfigmapSubscirption) Reconcile(object runtime.Object, event watch.EventType) {
	configmap := object.(*v1.ConfigMap)
	klog.Infof("config Map subcription %s from %s", event, configmap.Name)
	if ok, err := isPlatformConfigmap(configmap); !ok {
		if err != nil {
			klog.Error(err)
		}
		return
	}
	switch event {
	case watch.Added:
		//populate
		c.PlatformConfigPhase = string(event)
		rawString := configmap.Data["default"]
		var unMarshalData platformConfig
		err := yaml.Unmarshal([]byte(rawString), &unMarshalData)
		if err != nil {
			klog.Error(err)
			return
		}
		c.PlatformConfigMapAnnotations = &unMarshalData
	case watch.Deleted:
		//wipe it
		c.PlatformConfigPhase = string(event)
		c.PlatformConfigMapAnnotations = nil
	case watch.Modified:
	}
}

func (c *ConfigmapSubscirption) Subscribe() (watch.Interface, error) {
	var err error
	c.WatcherInterface, err = c.Client.CoreV1().ConfigMaps("").Watch(c.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	klog.Info("started working stream for configmap")
	return c.WatcherInterface, nil
}

func (c *ConfigmapSubscirption) IsComplete() <-chan bool {
	return c.CompletionChan
}
