# Kubernetes Leader Election Demo

This is a practical demonstration of how to use Kubernetes leases for leader election in a distributed application. The demo shows how multiple instances of an application can coordinate to elect a leader using Kubernetes leases.

Working demonstration of https://medium.com/@sanjimoh/understanding-kubernetes-leases-from-concept-to-implementation-bc423e868276

## Prerequisites

- Go 1.21 or later
- Docker
- Kubernetes cluster (local or remote)
- kubectl configured to work with your cluster

## Building the Application

1. Create a Kind cluster:
```bash
kind create cluster --name leader-election-cluster
```

2. Build the Docker image:
```bash
docker build -t leader-election-demo:latest .
```

3. Loading the Docker image into leader-election-cluster cluster:
```bash
kind load docker-image leader-election-demo:latest --name leader-election-cluster
```

## Running the Demo

1. Apply the rbac and Kubernetes deployment:
```bash
kubectl apply -f k8s/rbac.yaml
kubectl apply -f k8s/deployment.yaml
```

2. Watch the pods and their logs:
```bash
kubectl get pods -w
kubectl logs -f deployment/leader-election-demo
```

## How it Works

The application uses Kubernetes leases to implement leader election:

1. Each pod attempts to acquire a lease using the Kubernetes API
2. Only one pod can hold the lease at a time
3. The pod holding the lease becomes the leader
4. If the leader pod fails or is deleted, another pod will acquire the lease and become the new leader

The lease has the following timing parameters:
- Lease Duration: 15 seconds
- Renew Deadline: 10 seconds
- Retry Period: 2 seconds

## Observing Leader Election

You can observe the leader election process by:

1. Watching the pod logs:
```bash
kubectl logs -f deployment/leader-election-demo
```

2. Checking the lease object:
```bash
kubectl get lease leader-election-demo -o yaml
```

3. Deleting the leader pod to see failover:
```bash
kubectl delete pod <leader-pod-name>
```

## Cleanup

To clean up the demo:
```bash
kubectl delete -f k8s/deployment.yaml
``` 
