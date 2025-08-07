package main

import (
	"context"
	"fmt"
	"k8soperator/pkg/runtime"
	"k8soperator/pkg/subscription"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {
	kubeconfig := filepath.Join(
		homeDir(), ".kube", "config",
	)
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}
	fmt.Printf("listing to kube api %v \n", cfg.Host)
	defaultClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("error from building watcher %s", err.Error())
	}
	// wiface, err := defaultclientset.CoreV1().Pods("").Watch(context.TODO(), v1.ListOptions{})
	// if err != nil {
	// 	klog.Fatalf("watch interface  : %s", err.Error())
	// }
	//genral conetxt
	ctx := context.TODO()

	configMapSubscription := &subscription.ConfigmapSubscirption{
		Client:         defaultClientset,
		Ctx:            ctx,
		CompletionChan: make(chan bool),
	}

	podSubscription := &subscription.PodSubcription{
		Client:                defaultClientset,
		Ctx:                   ctx,
		CompletionChan:        make(chan bool),
		ConfigMapSubscriptRef: configMapSubscription,
	}

	if err := runtime.RunLoop([]subscription.Isubscription{
		configMapSubscription,
		podSubscription,
	}); err != nil {
		log.Fatalf("RunLoop error: %v", err)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
