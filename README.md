# Gazan. Cloud native load balancer.

Gazan is a free, very fast and reliable reverse-proxy offering high availability, load balancing, authenticating and proxying for HTTP-based applications. 
It is particularly suited for cloud environments where realtime configuration via API call and high performance is needed.  
Gazan is single file Go application, which is designed to run in Linux machines.
It comes with simple json based API server for configuring upstreams or integrates with DNS servers and reads SRV records .   
Gazan can load balance applications based on first level URL path or do wildcard load balancing based on DNS name or IP address of request.  

To install Gazan you need just to download the latest binary and config.yml file from our GitHub and run it. 

Some sample config files for feeding API server are included in `cfgjson` folder of sourcetree.

```shell
curl -XPOST -u 'test:Te$ting' --data-binary @/tmp/balod.json 127.0.0.1:4141/config?cfg=new
curl -XPOST -u 'test:Te$ting' --data-binary @/tmp/valod.json 127.0.0.1:4141/config?cfg=append
curl -u 'test:Te$ting' 127.0.0.1:4141/config?cfg=get
```

Following is sample config file for Gazan from GitHub :  

```yaml
main:
  listen : 0.0.0.0:8080
  confurl : 127.0.0.1:4141
  healtchecks : 2
  authtype: none
  accesslog: off
dnsconfig:
  enabled : yes
  dnsservers:
    - 192.168.10.104:8600
    - 192.168.10.105:8600
    - 192.168.10.106:8600
  srvrecords:
    - name    : p1.netangels.net:8080
      address : nginx-consul-NginX-health.tcp.service.consul
    - name    : p2.netangels.net:8080
      address : devecho-devEcho-health.tcp.service.consul
tls:
  enabled : no
  certificate : /tmp/fullchain.cer
  privatekey : /tmp/netangels.net.key
  listen : 0.0.0.0:8443
client:
  clientauth : no
  clientuser : test
  clientpass : Te$ting
  maxidle : 100
  maxperhost : 10
  maxidleperhost : 10
  idletimeout : 90
  timeout : 10
monitoring:
  enabled : yes
  url:  127.0.0.1:9191
```
## Config Sections: 

### main
Configure basic setting, like bind address etc ...  
**listen** : Bind address port for running the server '0.0.0.0' means run on all interfaces.

**confurl** : Address and port for accepting API calls for configuring upstreams.

**healtchecks** : Upstreams health check interval in second.

**accesslog** : off

**authtype** : Enable and set authentication parameter.(apikey , basic, jwt, none) for client requests.

If authentication is enabled , you should expose one of following OS enviroment variables.

**$GAZANKEY** : API key for `apikey` authentication 

**$BASICUSER, $BASICPASS** : Username and password for basic auth

**$JWTSECRET** : if you want to use JWT authentication.  

### dnsconfig
Enable and configure DNS API integration. 

**enabled** : Should we use DNS servers for getting upstreams. (yes/no) 

**dnsservers** : list of DNS servers IP:PORT

**srvrecords** : SRV records to query. Where `name` is reques url which should be proxied,
`address` is SRV record which Gazan should get in order get address and port of upstreams.   

### tls

Configuration parameters for TLS server. 

**enabled** : Enable or disable TLS listener (yes/no)

**certificate** : Full path to TLS full chain file

**privatekey** : Full path to private kwy file

**listen** : Bind address port for TLS server '0.0.0.0' means run on all interfaces.

### client

Configuration parameters for HTTP proxy client:

**clientauth** : should we authenticate request to upstreams (yes/no). Basic auth only.

**clientuser** : Client username

**clientpass** : Client password

**maxidle** : Number of maximum idle connections

**maxperhost** : Maximum connections per upstream

**maxidleperhost** : Number of maximum idle connections per host

**idletimeout** : Timeout for idle connection

**timeout** : Timeout for client connect. 

## monitoring 
Collect and expose some metrics about Gazan server 

**enabled** : Run metrics server (yes/no)

**url**:  Address:Port to bind metrics server