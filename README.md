# Kubernetes CRD and Custom Resource

In Kubernetes, you can think of a "Custom Resource" as a way to create your own special objects to manage things that are unique to your application. These special objects are like the standard things Kubernetes knows how to handle, like Pods and Services, but they are tailored to your specific needs.

Imagine you're running a video game on Kubernetes, and you want to create a custom resource to represent a new type of in-game item. You can define the properties of this item, like its name, power, and special abilities, using something called a "Custom Resource Definition" (CRD).

# Current Environment

- Go (1.18.5)
- Kubectl (GitCommit : 1b4df30b3, Git Version: v1.27.0)
- KubeBuilder (3.5.0)
- Packages need to Install using apt such as make, build-essential
- Linux/AMD64

# Establish Current Environment
- Open Bash Terminal.
- Copy the config script command from the git repo main branch and then execute the script.
```
controlplane ~ ➜ vi config.sh
controlplane ~ ➜ # Copy config.sh from this git repo.

controlplane ~ ➜ bash config.sh
```

- Your Environment is Ready. Now we will see the next steps in the "p1 Branch".
