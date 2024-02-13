# Introduction to Custom Controller

The repository contains a Kubernetes controller written in Go that manages custom resources of typeDeployService. The controller watches for changes to DeployService resources and takes corresponding actions to ensure the desired state is maintained.

## Demo Recording

The demo is recorded in "CRD_Part2.mp4".

## Main Code

- **Import Statements**: The code imports various packages required for interacting with Kubernetes, such as flag for command-line arguments, os for handling operating system signals, k8s.io packages for Kubernetes client and API objects, and github.com for the custom clientset.

- **Global Variables**: Defines a command-line flag -kubeconfig to specify the path to the Kubernetes configuration file.

- **Main Function**:The main function is the entry point of the program.

- **SetupSignalHandler Function**:This function sets up a signal handler to catch signals like SIGINT and SIGTERM and returns a channel that can be used to stop the application gracefully.

- **Inside main Function**:
    - The program starts by printing a message and setting up a signal handler.
    - Parses command-line flags, specifically the Kubernetes configuration file path.
    - Builds a Kubernetes client configuration from the provided config file.
    - Initializes Kubernetes clientsets for both the standard client and the custom client.
    - Sets up event broadcasting for logging and recording events related to the controller.
    - Creates a new instance of the controller (c) using the custom client and standard client.
    - Runs the controller with one worker, handling any errors that may occur during the process.

Overall, this program is a Kubernetes controller that watches for changes in resources and takes actions accordingly using a custom client.

## Controller Code

This code represents a Kubernetes controller written in Go using the client-go library. Kubernetes controllers are programs that watch for changes to resources in a Kubernetes cluster and take action based on those changes.

Here's a breakdown of the main components in this code:

- **Controller Struct**:This struct contains various fields that hold information needed for the controller to function.
    - **podsSynced**: An InformerSynced object that tracks the synchronization status of the DeployService resources.
    - **queue**: A work queue for handling resource events.
    - **deployServiceLister**: A lister for DeployService resources.
    - **kubeClient**: A Kubernetes client used for interacting with the core Kubernetes API.
    - **newClient**: A client for interacting with the custom API (DeployV1).

- **NewController Function**:
    - Creates a new instance of the Controller struct.
    - Sets up a work queue and shared informers for the DeployService resources.
    - Defines event handlers for resource addition and update events.
    - Starts the informer factory and initializes the controller fields.

- **Run Function**:
    - Runs the controller by starting worker goroutines.
    - Ensures secondary caches are filled before starting work.
    - Workers process items from the work queue until the context is canceled.

- **RunWorker and ProcessNextWorkItem Functions**:
    - RunWorker is a worker goroutine that continuously calls ProcessNextWorkItem until instructed to stop.
    - ProcessNextWorkItem retrieves an item from the work queue, processes it, and handles retries or errors.

- **ProcessDeployService Function**:
    - **Start**: The function begins by checking if the DeployService resource has any finalizers. If not, it adds the finalizer for proper cleanup during deletion.
    - **Deletion Check**: If the DeployService resource has a deletion timestamp set, it means the resource is being deleted.The function proceeds to delete associated Deployments and Services.
    - **Ensure Deployment and Service**: If the Deployment associated with the DeployService is not created (Status.DepCreated is false), the function creates the Deployment based on the provided template. Similarly, if the Service associated with the DeployService is not created (Status.SvcCreated is false), the function creates the Service based on the provided template.
    - **Update Status**: If any updates are made (new resources created or status changed), the function updates the DeployService status with the latest information. This includes information such as Deployment creation status, image, and replica count.
    - **End**: The function either updates the status and returns the updated DeployService or returns the original DeployService if no updates were made.

- **newDeployment and newService Functions**:
    - Helper functions to create Kubernetes Deployment and Service objects based on template specifications.

- **EnsureFinalizers and RemoveFinalizers Functions**:  
    - EnsureFinalizers adds finalizers to the DeployService resource to ensure proper cleanup during deletion.
    - RemoveFinalizers removes finalizers from the DeployService resource during cleanup.

- **Finalizer Constant**:
    - DeployServiceFinalizers is a constant string representing the finalizer used by the controller.

In summary, this controller ensures the proper deployment, updating, and deletion of associated Kubernetes resources based on the state of DeployService custom resources, following a pattern of shared informers, event handlers, and a work queue.

## Makefile

The provided Makefile contains a set of targets to simplify various tasks related to building, running, and managing a Go application that includes a Kubernetes controller for managing custom resources. Let's break down each target:

- **init**:
    - Creates necessary directories in the user's home directory for Go development (${HOME}/go/bin, ${HOME}/go/pkg, ${HOME}/go/src/github.com/DeployService).
    - Copies the contents of the current directory to ${HOME}/go/src/github.com/DeployService.
    - Prints a message with instructions on changing to the created directory.
- **env**: Sets environment variables for Go development:
    - GOPATH is set to ${HOME}/go.
    - GO111MODULE is set to off.
- **deps**: Uses go get to fetch and install the dependencies specified in the go.mod file.
- **create_crd**: Uses kubectl to create a custom resource definition (CRD) based on the artifacts/crd.yaml file.
- **build**: Builds the Go application, producing an executable named main.
- **run**: Executes the built main binary, providing the Kubernetes configuration file path as an argument.
- **get_cr_yaml**: Prints the content of the example custom resource (CR) YAML file (artifacts/example.yaml).
- **get_all**: Displays information about DeployServices, Deployments, and Services using kubectl get commands.
- **create_cr**: Uses kubectl to create an instance of the custom resource defined in artifacts/example.yaml.
- **desc**: Uses watch and kubectl describe to continuously monitor and display detailed information about a specific DeployService (example-kube1).
- **cleanup**: Uses kubectl to delete a specific DeployService (example-kube1).

These targets provide a convenient way to initialize the development environment, manage dependencies, build and run the application, interact with Kubernetes resources, and perform cleanup tasks. Users can run these targets individually based on their needs.
