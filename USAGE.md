### Bootstrap

- Set up Kubernetes cluster using the way you prefer. I use `minikube` for testing purposes.
- Create resources required
```
kubectl apply -f k8s/nats.yaml -f k8s/pg.yaml -f k8s/bank.yaml
```
- Create databases and run migrations
```
kubectl apply -f k8s/create_db.yaml -f k8s/migrations.yaml
```
- Create services
```
kubectl apply -f k8s/registry.yaml -f k8s/api.yaml
```

### Usage

It is required to set up port forwarding to connect services in cluster.
There are different ways this can be done:
- using minikube
```
minikube service api --url
```
- using kubectl
```
kubectl port-forward service/api <local-port>:8010
```

Then you can use SwaggerUI to interact with the API service
```
http://127.0.0.1:8010/docs/#!/default/
```
