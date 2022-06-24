## Simple app in kubernetes cluster

### Create docker image

#### Build binary with command

```
$ CGO_ENABLED=0 GOOS=linux go build .
``` 

Cope binary to Dockerfile (look Dockerfile)

#### Build image
```
$ sudo docker build --no-cache --tag simple_app_for_kube:v1.0.0 .
```

If you have problem with postgres, switch on log in postgres: 
in postgresql.conf change logging_collector, log_directory, 

### Run docker image with network(postgres running in another network)
```
$ sudo docker run -it --net=workdocker_workdocker_default -p 11101:11101 simple_app_for_kube:v1.0.0
```

or allow communication between networks(if you run without --net container will run in default network 'bridge' and interface docker0 usually)
```
$ sudo iptables -I DOCKER-USER -i docker0 -o br-82bedfb56fef -j ACCEPT
$ sudo iptables -I DOCKER-USER -i br-82bedfb56fef -o docker0 -j ACCEPT
```

To get bridge name use command:
```
$ ifconfig | grep 172.18.0. -B 1
```

### Docker compose
Checking app via docker-compose(look docker-compose.yml):
```
$ sudo docker-compose up -d
```
####!!! The host of app must be 0.0.0.0 instead 127.0.0.1 (in config.yaml)

example request: 
```
$ curl --request GET '0.0.0.0:11101/users?userId=5'
```

### Create kube cluster 
For use local image - look instruction in test_pod.yaml
Here I use external PG for kube cluster (how to connect app to external PG - look pg.yaml)
```
$ minikube start
$ kubectl apply -f ./k8s/pg.yaml
$ kubectl apply -f ./k8s/deployment.yaml
```
For checking should to get minikube ip:
```
$ minikube ip
192.168.59.100
```
example request:
```
$ curl --request GET '192.168.59.100:30001/users?userId=4'
"name"
```
#### Helm
Create the same cluster via helm (look folder k8s/helm-chart)
```
$ helm install simple-app ./k8s/helm-chart/
```
Set replica count by --set parameters
```
$ helm install simple-app ./k8s/helm-chart/ --set replicaCount=3
```

### Access by domain
Useful information: https://medium.com/@emirmujic/istio-and-metallb-on-minikube-242281b1134b

#### Add istio
```
istioctl install --set profile=demo -y
```
Vertify install
```
istioctl verify-install -f ~/istio-1.14.0/manifests/profiles/demo.yaml
```
Set label istio-injection
```
kubectl label namespace default istio-injection=enabled
```

Check
```
istioctl analyze
Warning [IST0103] (Pod default/simple-app-simple-deployment-854c8b84f9-5bxg5) The pod default/simple-app-simple-deployment-854c8b84f9-5bxg5 is missing the Istio proxy. This can often be resolved by restarting or redeploying the workload.
```

Restart deployment and analyze again
```
kubectl rollout restart deployment simple-app-simple-deployment
```
No warnings more

- Add gateway and virtual service (look gateway.yaml and virtual-service.yaml)
- Change service type to ClusterIP and comment field nodePort in template (look deployment.yaml)
```
$ helm upgrade simple-app ./k8s/helm-chart/
```

Check virtual services
```
$ kubectl get virtualservices
```

#### Install MetalLB
```
$ kubectl create namespace metallb-system
$ kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.12.1/manifests/metallb.yaml
```
Create configMap for metalLB(look config-map.yaml) and upgrade
```
$ minikube ip
192.168.59.100

$ helm upgrade simple-app ./k8s/helm-chart/
```

Check externalIp for service istio-ingressgateway 
```
$ kubectl get svc/istio-ingressgateway -n istio-system
istio-ingressgateway   LoadBalancer   10.98.2.241   192.168.59.97 
```

#### Get externalIP and add host to /ets/hosts/ 
In my case: 192.168.59.97 simple-app.org
```
$ sudo vi /etc/hosts
```

Send requests
```
$ curl --location --request GET 'simple-app.org/users?userId=6'
"na1d2me"
$ curl --location --request POST 'simple-app.org/users' --header 'Content-Type: application/json' --data-raw '{"name":"Test"}'
11
```