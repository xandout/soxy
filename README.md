# SOXY

Soxy has a simple design.  Run a `soxy` server behind your Ingress or other HTTP reverse proxy such as nginx.

When a client connects to the `soxy` server, it asks the server to establish a websocket connection to a back-end service.  Once the connection is established, the client will accept connections on a local port and forward traffic to the remote end.  

My test rig is a k8s cluster with `ingress-nginx` providing HTTPS termination.  This makes it easy to secure the traffic on your tunnels.


## Server mode

This will start a `soxy` server which listens on port 8080 for HTTP connections.
```
soxy serve -p :8080
```



## Client mode

This will start a client on your local laptop, proxying connections on port 8479 to `mongodb-service:27017` as seen from the proxy.

```
soxy proxy -U ws://soxy-server.com -L :8479 -R mongodb-service:27017
```

If you have `soxy` behind an HTTPS ingress or reverse proxy, you need to use 

```
soxy proxy -U wss://soxy-server.com -L :8479 -R mongodb-service:27017
```
> Notice the extra `s` in `wss`


You can now connect to the remote service from your workstation as follows

```
mongo --host 127.0.0.1 --port 8479
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