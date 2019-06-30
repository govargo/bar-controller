package main

import (
	"flag"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	barclientset "github.com/govargo/bar-controller/pkg/generated/clientset/versioned"
	barinformers "github.com/govargo/bar-controller/pkg/generated/informers/externalversions"
	fooclientset "k8s.io/sample-controller/pkg/generated/clientset/versioned"
	fooinformers "k8s.io/sample-controller/pkg/generated/informers/externalversions"
	"k8s.io/sample-controller/pkg/signals"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	fooClient, err := fooclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building foo clientset: %s", err.Error())
	}

	barClient, err := barclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building bar clientset: %s", err.Error())
	}

	fooInformerFactory := fooinformers.NewSharedInformerFactory(fooClient, time.Second*30)
	barInformerFactory := barinformers.NewSharedInformerFactory(barClient, time.Second*30)

	controller := NewController(kubeClient, fooClient, barClient,
		fooInformerFactory.Samplecontroller().V1alpha1().Foos(),
		barInformerFactory.Samplecontroller().V1alpha1().Bars())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	fooInformerFactory.Start(stopCh)
	barInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
