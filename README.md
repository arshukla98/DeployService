# Kubernetes CRD and Custom Resource

In Kubernetes, you can think of a "Custom Resource" as a way to create your own special objects to manage things that are unique to your application. These special objects are like the standard things Kubernetes knows how to handle, like Pods and Services, but they are tailored to your specific needs.

Imagine you're running a video game on Kubernetes, and you want to create a custom resource to represent a new type of in-game item. You can define the properties of this item, like its name, power, and special abilities, using something called a "Custom Resource Definition" (CRD).

To establish current environment, execute steps in main branch.

# Current Environment

- Go (1.18.5)
- Kubectl (GitCommit : 1b4df30b3, Git Version: v1.27.0)
- KubeBuilder (3.5.0)
- Packages need to Install using apt such as make, build-essential
- Linux/AMD64

# Objective

The purpose of this repo is to create the following custom resource as shown below by creating CRD.

```
apiVersion: deploy.msys.io/v1
kind: DeployService
metadata:
  name: example-kube1
spec:
  deploymentTemplate:
    imageName: nginx:alpine
    name: nginx
    namespace: default
    replicas: 3
  serviceTemplate:
    name: nginx-svc
    type: NodePort
    servicePort: 80
```

The provided YAML code defines a custom resource named DeployService with the API version deploy.msys.io/v1. 
Custom Resources (CRs) in Kubernetes allow you to define and use your own custom objects alongside the built-in Kubernetes objects.

Let's break down the structure of this custom resource:

- **apiVersion**: Specifies the API version for the custom resource. In this case, it's deploy.msys.io/v1.

- **kind**: Specifies the kind of the custom resource, which is DeployService.

- **metadata**: Contains metadata information for the custom resource, including its name. In this example, the name is set to example-kube1.

- **spec**: Contains the specification for the custom resource. It defines the desired state for the deployment and service associated with this custom resource.

- **deploymentTemplate**: Describes the template for a Kubernetes Deployment.

    - **imageName**: Specifies the Docker image to be used for the deployment. In this case, it's nginx:alpine.

    - **name**: Specifies the name of the deployment, which is set to nginx.

    - **namespace**: Specifies the Kubernetes namespace where the deployment should be created. In this example, it's set to default.

    - **replicas**: Specifies the number of replicas (pods) for the deployment. In this case, it's set to 3.

- **serviceTemplate**: Describes the template for a Kubernetes Service.

    - **name**: Specifies the name of the service, which is set to nginx-svc.

    - **type**: Specifies the type of the service. In this case, it's set to NodePort, which exposes the service on each Node's IP at a static port.

    - **servicePort**: Specifies the port on which the service should listen. In this example, it's set to 80.

In summary, this custom resource is defining a deployment of three replicas using the nginx:alpine Docker image and a NodePort service exposing port 80. This allows you to encapsulate and manage deployment and service configurations using a custom resource in Kubernetes.

Note that this is a like an object such as objects in Java. To add behaviors and methods associated with the object i.e Custom Resource, we already created Custom controller in 'p2 branch'.