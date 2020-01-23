# SOXY

Simple websocket powered proxy



## Server mode

```
soxy -f server.yaml
```

```yaml
# server.yaml
http:
  port: 8080
```



## Client mode

```
soxy -f client.yml -t mysql 
```

```yaml
# client.yml
tunnels:
  - name: mysql
    remote: mysql-service:3306
    local: 3306
```