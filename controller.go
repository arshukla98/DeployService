package main

import (
        "context"
        "fmt"
        "reflect"
        "strings"
        "time"

        //kubeinformers "k8s.io/client-go/informers"

        clientset "github.com/DeployService/generated/clientset/versioned"
        informers "github.com/DeployService/generated/informers/externalversions"
        appsv1 "k8s.io/api/apps/v1"
        corev1 "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        utilruntime "k8s.io/apimachinery/pkg/util/runtime"
        wait "k8s.io/apimachinery/pkg/util/wait"
        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/tools/cache"
        "k8s.io/client-go/util/workqueue"

        upgradev1 "github.com/DeployService/pkg/apis/deploy/v1"
        listers "github.com/DeployService/generated/listers/deploy/v1"
)

const (
        DeployServiceFinalizers string = "Deploy-Service-Finalizers"
)

type Controller struct {
        podsSynced cache.InformerSynced

        queue workqueue.RateLimitingInterface

        deployServiceLister listers.DeployServiceLister

        kubeClient *kubernetes.Clientset

        newClient *clientset.Clientset
}

func NewController(client *clientset.Clientset, kubeClient *kubernetes.Clientset) *Controller {
        fmt.Println("creating WorkQueue")
        queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "controller-name")

        fmt.Println("creating Shared Informers Factory")
        //kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*60)

        exampleInformerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
        deployInformer := exampleInformerFactory.Deploy().V1().DeployServices()

        fmt.Println("Setting up event handlers")
        deployInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
                        key, err := cache.MetaNamespaceKeyFunc(obj)
                        if err == nil {
                                queue.Add(key)
                        }
                        fmt.Println("Handling Add Resource Event",key)
                },
                UpdateFunc: func(old interface{}, newo interface{}) {
                        if !reflect.DeepEqual(old, newo) {
                                key, err := cache.MetaNamespaceKeyFunc(newo)
                                if err == nil {
                                        queue.Add(key)
                                }
                                fmt.Println("Handling Update Resource Event", key)
                        } else {
                                fmt.Println("not updated")
                        }
                },
                // It is not needed because we want to handler
                // cleanup tasks while resource is being deleted.
                /*DeleteFunc: func(obj interface{}) {
                        key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
                        if err == nil {
                                queue.Add(key)
                        }
                        fmt.Println("Handling Delete Resource Event:", key)
                },*/
        })

        exampleInformerFactory.Start(context.Background().Done())

        c := &Controller{}
        c.newClient = client
        c.kubeClient = kubeClient
        c.deployServiceLister = deployInformer.Lister()
        c.podsSynced = deployInformer.Informer().HasSynced
        c.queue = queue
        return c
}

func (c *Controller) Run(ctx context.Context, threadiness int) error {
        // don't let panics crash the process
        defer utilruntime.HandleCrash()
        // make sure the work queue is shutdown which will trigger workers to end
        defer c.queue.ShutDown()

        fmt.Println("Starting <NAME> controller")

        // wait for your secondary caches to fill before starting your work
        if !cache.WaitForCacheSync(ctx.Done(), c.podsSynced) {
                return fmt.Errorf("failed to wait for caches to sync")
        }

        for i := 0; i < threadiness; i++ {
                // runWorker will loop until "something bad" happens.  The .Until will
                // then rekick the worker after one second
                go wait.UntilWithContext(ctx, c.runWorker, time.Second)
        }

        fmt.Println("Started workers")
        <-ctx.Done()
        fmt.Println("Shutting down workers")

        return nil
}

func (c *Controller) runWorker(ctx context.Context) {
        // hot loop until we're told to stop.  processNextWorkItem will
        // automatically wait until there's work available, so we don't worry
        i := 0
        fmt.Println("Begin runWorker")
        for c.processNextWorkItem(ctx) {
                fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
                i += 1
                fmt.Println("i:", i)
        }
        fmt.Println("End runWorker")
}

// processNextWorkItem deals with one key off the queue.  It returns false
// when it's time to quit.
func (c *Controller) processNextWorkItem(ctx context.Context) bool {
        fmt.Println("Begin processNextWorkItem")
        key, quit := c.queue.Get()
        if quit {
                return false
        }
        defer c.queue.Done(key)

        fmt.Println("processNextWorkItem key:", key.(string))

        keyparts := strings.Split(key.(string),"/")

        ds, err := c.deployServiceLister.DeployServices(keyparts[0]).Get(keyparts[1])
        if err != nil {
                fmt.Println("Error Getting deployService CR", err.Error())
                return true
        }
        fmt.Printf("%+v\n",ds)
        fmt.Printf("%p\n",ds)

        // do your work on the key.  This method will contains your "do stuff" logic
        up, err := c.ProcessDeployService(context.Background(), ds.DeepCopy())
        if err == nil {
                fmt.Printf("%+v\n",up)
                fmt.Printf("%p\n",up)
                fmt.Println("processNextWorkItem Forget key:", key.(string))
                c.queue.Forget(key)
        } else if c.queue.NumRequeues(key) < 5 {
                // err != nil and retry
                fmt.Println("processNextWorkItem if 2nd block retry")
                c.queue.AddRateLimited(key)
        } else {
                // err != nil and too many retries
                fmt.Println("processNextWorkItem final else block")
                c.queue.Forget(key)
                utilruntime.HandleError(fmt.Errorf("%v failed with : %v", key, err))
        }
        fmt.Println("End processNextWorkItem")
        return true
}

func (c *Controller) ProcessDeployService(ctx context.Context, ds *upgradev1.DeployService) (*upgradev1.DeployService, error) {
        defer fmt.Println("Exiting ProcessDeployService")
        fmt.Println("Entering ProcessDeployService")

        updated := false

        if len(ds.ObjectMeta.Finalizers) == 0 {
                fmt.Println("Add Finalizers for Proper Cleanup Tasks")
                ds = c.ensureFinalizers(ctx, ds)
                updated = true
                fmt.Println("Resource Finalizers needs to be updated.")
                updatedDS, err := c.newClient.DeployV1().DeployServices(ds.ObjectMeta.Namespace).Update(context.Background(), ds, metav1.UpdateOptions{})
                if err != nil {
                        fmt.Println("Error Updating DeployServices", err.Error())
                        return nil, err
                }
                return updatedDS, nil
        }

        if !ds.ObjectMeta.DeletionTimestamp.IsZero() {
                namespace, name := ds.Spec.DepTemp.Namespace, ds.Spec.DepTemp.Name
                err := c.kubeClient.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
                if err != nil {
                        return nil, err
                }

                namespace, name = ds.Spec.DepTemp.Namespace, ds.Spec.ServiceTemp.Name
                err = c.kubeClient.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
                if err != nil {
                        return nil, err
                }
                ds = c.removeFinalizers(ctx, ds)

                updatedDS, err := c.newClient.DeployV1().
                        DeployServices(ds.Namespace).
                        Update(ctx, ds, metav1.UpdateOptions{})
                if err != nil {
                        return nil, err
                }
                return updatedDS, nil
        }

        depTemp := &ds.Spec.DepTemp
        /*
        if ds.Status.DepCreated {
                if depTemp.ImageName != ds.Status.DepImage {
                        namespace, name := ds.Spec.DepTemp.Namespace, ds.Spec.DepTemp.Name
                        err := c.kubeClient.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
                        if err != nil {
                                return nil, err
                        }

                        ds.Status.DepCreated = false
                        ds.Status.DepImage = ""
                        ds.Status.DepReplicasCount = 0

                        fmt.Println("Deployment Image needs to be updated.")
                        updatedDS, err := c.newClient.DeployV1().DeployServices(ds.ObjectMeta.Namespace).UpdateStatus(context.Background(), ds, metav1.UpdateOptions{})
                        if err != nil {
                                fmt.Println("Error Updating DeployServices", err.Error())
                                return nil, err
                        }
                        return updatedDS, nil
                }

        }
        */
        if !ds.Status.DepCreated && depTemp != nil {
                //time.Sleep(30*time.Second)
                obj := newDeployment(depTemp)
                deployment, err := c.kubeClient.AppsV1().Deployments(depTemp.Namespace).Create(ctx, obj, metav1.CreateOptions{})
                if err != nil {
                        return nil, err
                }
                if deployment == nil {
                        return nil, fmt.Errorf("nil deployment")
                }
                ds.Status.DepCreated = true
                ds.Status.DepImage = depTemp.ImageName
                ds.Status.DepReplicasCount = depTemp.Replicas
                updated = true
        }

        servTemp := &ds.Spec.ServiceTemp

        if !ds.Status.SvcCreated && servTemp != nil {
                obj := newService(servTemp, depTemp)
                service, err := c.kubeClient.CoreV1().Services(depTemp.Namespace).Create(ctx, obj, metav1.CreateOptions{})
                if err != nil {
                        return nil, err
                }
                if service == nil {
                        return nil, fmt.Errorf("nil service")
                }
                ds.Status.SvcCreated = true
                updated = true
        }
        if updated {
                fmt.Println("Resource needs to be updated.")
                updatedDS, err := c.newClient.DeployV1().DeployServices(ds.ObjectMeta.Namespace).UpdateStatus(context.Background(), ds, metav1.UpdateOptions{})
                if err != nil {
                        fmt.Println("Error Updating DeployServices", err.Error())
                        return nil, err
                }
                return updatedDS, nil
        }
        fmt.Println("Resource does not need to be updated.")
        return ds, nil
}

func newDeployment(dep *upgradev1.DeploymentTemplate) *appsv1.Deployment {
        labels := map[string]string{
                "app":        dep.Name,
                "controller": "sample-controller",
        }
        return &appsv1.Deployment{
                ObjectMeta: metav1.ObjectMeta{
                        Name: dep.Name,
                },
                Spec: appsv1.DeploymentSpec{
                        Replicas: &dep.Replicas,
                        Selector: &metav1.LabelSelector{
                                MatchLabels: labels,
                        },
                        Template: corev1.PodTemplateSpec{
                                ObjectMeta: metav1.ObjectMeta{
                                        Labels: labels,
                                },
                                Spec: corev1.PodSpec{
                                        Containers: []corev1.Container{
                                                {
                                                        Name:  dep.Name,
                                                        Image: dep.ImageName,
                                                },
                                        },
                                },
                        },
                },
        }
}

func newService(serviceTemp *upgradev1.ServiceTemplate, dep *upgradev1.DeploymentTemplate) *corev1.Service {
        labels := map[string]string{
                "app":        dep.Name,
                "controller": "sample-controller",
        }
        return &corev1.Service{
                ObjectMeta: metav1.ObjectMeta{
                        Name: serviceTemp.Name,
                },
                Spec: corev1.ServiceSpec{
                        Selector: labels,
                        Type:     corev1.ServiceType(serviceTemp.Type),
                        Ports: []corev1.ServicePort{
                                {
                                        Port:     serviceTemp.ServicePort,
                                        Protocol: corev1.ProtocolTCP,
                                },
                        },
                },
        }
}

func (c *Controller) ensureFinalizers(ctx context.Context, ds *upgradev1.DeployService) *upgradev1.DeployService {
        fmt.Println("Entering ensureFinalizers")
        defer fmt.Println("Exiting ensureFinalizers")

        if ds.ObjectMeta.DeletionTimestamp.IsZero() {
                if len(ds.ObjectMeta.Finalizers) == 0 {
                        fmt.Println("Adding Finalizers to DeployService")
                        ds.ObjectMeta.Finalizers = append(ds.ObjectMeta.Finalizers, DeployServiceFinalizers)
                }
        }
        return ds
}

func (c *Controller) removeFinalizers(ctx context.Context, ds *upgradev1.DeployService) *upgradev1.DeployService {
        fmt.Println("Entering removeFinalizers")
        defer fmt.Println("Exiting removeFinalizers")

        finalizers := ds.ObjectMeta.Finalizers

        if len(finalizers) > 0 {
                fmt.Println("Removing Finalizers from DeployService")
                for i, finalizer := range finalizers {
                        if finalizer == DeployServiceFinalizers {
                                ds.ObjectMeta.Finalizers = append(finalizers[0:i], finalizers[i+1:]...)
                        }
                }
        }
        return ds
}
