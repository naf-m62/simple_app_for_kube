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