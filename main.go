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
	defaultclientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("error from building watcher %s", err.Error())
	}
	// wiface, err := defaultclientset.CoreV1().Pods("").Watch(context.TODO(), v1.ListOptions{})
	// if err != nil {
	// 	klog.Fatalf("watch interface  : %s", err.Error())
	// }
	//genral conetxt
	context := context.TODO()
	if err := runtime.Runloop([]subscription.Isubscription{
		&subscription.ConfigmapSubscirption{
			Client:         defaultclientset,
			Ctx:            context,
			CompletionChan: make(chan bool),
		},
	}); err != nil {
		klog.Fatalf("this is the error => %v", err.Error())
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
