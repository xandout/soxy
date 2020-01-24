# SOXY [![Build Status](https://travis-ci.com/xandout/soxy.svg?branch=master)](https://travis-ci.com/xandout/soxy)

Soxy has a simple design.  Run a `soxy` server behind your Ingress or other HTTP reverse proxy such as nginx.

When a client connects to the `soxy` server, the server opens a connection to the backend service and pipes the data back to a local port.  

My test rig is a k8s cluster with `ingress-nginx` providing HTTPS termination.  This makes it easy to secure the traffic on your tunnels.


## Server mode

This will start a `soxy` server which listens on port 8080 for HTTP connections.
```
soxy serve -p :8080
```



## Client mode

This will start a client on your local laptop, proxying connections on port 8479 to `mongodb-service:27017` as seen from the proxy.

```
soxy proxy -U ws://soxy-server.com:8080 -L :8479 -R mongodb-service:27017
```

If you have `soxy` behind an HTTPS ingress or reverse proxy, you need to use 

```
soxy proxy -U wss://soxy-server.com:8080 -L :8479 -R mongodb-service:27017
```
> Notice the extra `s` in `wss`


You can now connect to the remote service from your workstation as follows

```
mongo --host 127.0.0.1 --port 8479
```


## DEMO

### docker-compose
The included [docker-compose.yml](docker-compose.yml) sets up a demo environment with 

* Frontend and backend networks to mimic real world topologies
* Backend Redis DB service
* Backend `soxy` service
* Frontend nginx service
* Frontend `soxy` client exposes `backend-redis:6379`


Data path is redis-cli -> local TCP socket(soxy) -> websocket connection through nginx -> soxy server -> redis container.

```
✔ ~/go/src/github.com/xandout/soxy [master|✔] 
02:10 # docker-compose up -d
Creating network "soxy_back" with the default driver
Creating network "soxy_front" with the default driver
Creating soxy_backend-redis_1 ... done
Creating soxy_backend-soxy-server_1 ... done
Creating soxy_frontend-nginx_1      ... done
Creating soxy_frontend-soxy-client_1 ... done
✔ ~/go/src/github.com/xandout/soxy [master|✔] 
02:14 # docker run --rm -it redis redis-cli -h 192.168.0.250 INFO CPU
# CPU
used_cpu_sys:0.025931
used_cpu_user:0.153783
used_cpu_sys_children:0.001324
used_cpu_user_children:0.001367
```


### Just a quick local demo proxying redis
```
✔ ~/go/src/github.com/xandout/soxy [master|✔] 
01:24 # docker run --rm -d -p 6379:6379 redis
7cefb999a6a23e03883a41776b74304506249ae0368b81cd8308865f826fe404
✔ ~/go/src/github.com/xandout/soxy [master|✔] 
01:22 # docker run --rm -d -p 8080:8080 xandout/soxy serve -p :8080
6690096b456e64baf3223df80f590ea2e0962c8c3f07bcdebcb9b0f25dadb3e6
✔ ~/go/src/github.com/xandout/soxy [master|✔] 
01:22 # docker run --rm -d -p 8081:8081 xandout/soxy proxy -U ws://192.168.0.250:8080 -L :8081 -R 192.168.0.250:6379
423ccbbfa4b6d961f4fb1b740e93bf38bfda9a96d03c56e3f8188b87d4a88d5b
✔ ~/go/src/github.com/xandout/soxy [master|✔] 
01:22 # docker run --rm -it redis redis-cli -h 192.168.0.250 -p 8081
192.168.0.250:8081> INFO CPU
# CPU
used_cpu_sys:86.098951
used_cpu_user:74.249437
used_cpu_sys_children:0.001757
used_cpu_user_children:0.000798
192.168.0.250:8081> 
```

### Secured traffic over the Internet to a k8s cluster
```
01:25 # docker run --rm -d -p 8082:8082 xandout/soxy proxy -U wss://soxy.my-kubernetes-cluster.com -L :8082 -R mongodb-service:27017
f7393b5b5254bd5c4ad0b4c8cb0ed3ac1cd0dc7c73bef909eca4cdf896bb8865
✔ ~/go/src/github.com/xandout/soxy [master|✔] 
01:26 # mongo --host 192.168.0.250 --port 8082
MongoDB shell version v4.0.3
connecting to: mongodb://192.168.0.250:8082/
WARNING: No implicit session: Logical Sessions are only supported on server versions 3.6 and greater.
Implicit session: dummy session
MongoDB server version: 3.0.12
WARNING: shell and server versions do not match
> show databases;
local               0.078GB
> 
```
## Usage

```
# soxy -h
NAME:
   soxy - fight the loneliness!

USAGE:
   soxy [global options] command [command options] [arguments...]

COMMANDS:
   serve    Start proxying traffic(server)
   proxy    Start proxying traffic(client)
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)


# soxy proxy -h
NAME:
   soxy proxy - Start proxying traffic(client)

USAGE:
   soxy proxy [command options] [arguments...]

OPTIONS:
   --soxy-url value, -U value  ws://soxy-daemon.com:8080
   --local value, -L value     Which local port to listen on.
                               Example: :3306 or 0.0.0.0:3306
   --remote value, -R value    Where should the daemon proxy traffic to?
                               Example: mysql-service:3306
   --help, -h                  show help (default: false)   


# soxy serve -h
NAME:
   soxy serve - Start proxying traffic(server)

USAGE:
   soxy serve [command options] [arguments...]

OPTIONS:
   --port value, -p value  
   --help, -h              show help (default: false)

```