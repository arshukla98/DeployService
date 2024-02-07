package main

import (
        "flag"
        "fmt"
        "os"
        "os/signal"
        "syscall"

        "context"
        corev1 "k8s.io/api/core/v1"
        kubernetes "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/kubernetes/scheme"
        typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
        "k8s.io/client-go/tools/clientcmd"
        "k8s.io/client-go/tools/record"
        "k8s.io/klog/v2"

        clientset "github.com/DeployService/generated/clientset/versioned"
)

var kubeconfig = flag.String("kubeconfig", "", "kube config")

func main() {
        fmt.Println("Beginning in main....")
        _ = SetupSignalHandler()

        flag.Parse()
        fmt.Println("passing kubeconfig 1")
        cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
        if err != nil {
                fmt.Println("Error building kubeconfig")
                klog.FlushAndExit(klog.ExitFlushTimeout, 1)
        }
        fmt.Println("passing kubeconfig 2")
        kubeClient, err := kubernetes.NewForConfig(cfg)
        if err != nil {
                fmt.Println("Error building kubernetes clientset")
                klog.FlushAndExit(klog.ExitFlushTimeout, 1)
        }
        fmt.Println("passing kubeconfig 3")
        exampleClient, err := clientset.NewForConfig(cfg)
        if err != nil {
                fmt.Println("Error building kubernetes clientset")
                klog.FlushAndExit(klog.ExitFlushTimeout, 1)
        }

        eventBroadcaster := record.NewBroadcaster()
        eventBroadcaster.StartStructuredLogging(0)
        eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
        _ = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "sample-controller"})

        c := NewController(exampleClient, kubeClient)

        ctx := context.Background()

        fmt.Println("Run controller with 1 workers")

        if err = c.Run(ctx, 1); err != nil {
                fmt.Println("Error running controller")
                os.Exit(1)
        }
}

func SetupSignalHandler() (stopCh <-chan struct{}) {
        stop := make(chan struct{})
        c := make(chan os.Signal, 2)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

        go func() {
                <-c
                fmt.Println("Received stop signal, attempting graceful termination...")
                close(stop)
                <-c
                fmt.Println("Received stop signal, terminating immediately!")
                os.Exit(1) // second signal. Exit directly.
        }()
        return stop
}
