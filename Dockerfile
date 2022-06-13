# syntax=docker/dockerfile:1
FROM alpine:latest

COPY simple_app_for_kube /
COPY cmd/cfg/config.yaml /

EXPOSE 11101
ENV POSTGRES_HOST 172.18.0.3

#RUN apk add --no-cache netcat-openbsd
#RUN apk add --no-cache bash

CMD ["/simple_app_for_kube"]

# before build image
# CGO_ENABLED=0 GOOS=linux go build .

# build image
# sudo docker build --no-cache --tag simple_app_for_kube:v1.0.0 .

# check
# sudo docker run simple_app_for_kube:v1.0.0

# sudo docker run -it --net=workdocker_workdocker_default -p 11101:11101 simple_app_for_kube:v1.0.0 - if pg in another docker network
# sudo docker network inspect workdocker_workdocker_default
