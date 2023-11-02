![c97b66731796c6df734fc0e290b856e](https://github.com/leetomlee123/book_api_public/assets/19498940/cba9f7aa-3621-41b0-944b-f9a6c88b1244)

# mongodb config
```aidl
## content
systemLog:
  destination: file
  logAppend: true
  path: /www/server/mongodb/log/config.log
 
# Where and how to store data.
storage:
  dbPath: /www/server/mongodb/data
  directoryPerDB: true

  journal:
    enabled: true
  wiredTiger:
    engineConfig:
      cacheSizeGB: 0.5
# how the process runs
processManagement:
  fork: true
  pidFilePath: /www/server/mongodb/log/configsvr.pid
 
# network interfaces
net:
  port: 27017
  bindIp: 0.0.0.0
 
#operationProfiling:
#replication:
#    replSetName: bt_main   
security:
  authorization: enabled
  javascriptEnabled: false

#sharding:
#    clusterRole: shardsvr
```
# docker elasticsearch
```aidl
docker run --name es -p 8088:9200 -p 9300:9300 -e "discovery.type=single-node" -e ES_JAVA_OPTS="-Xms64m -Xmx128m"  -v /root/lx/es/data:/usr/share/elasticsearch/data -v /root/lx/es/plugins:/usr/share/elasticsearch/plugins -d elasticsearch:7.8.0
```
