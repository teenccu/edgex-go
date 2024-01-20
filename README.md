# Core-data with Influx DB
Please refer the documentation of https://github.com/edgexfoundry/edgex-go for any details/licencing of edgex-go

This fork adds the implemenation of Core-data with InfluxDB and keeps the Core-Metadata still using the redisDB. 
A Hybrid client is created which implements the metadata and core-data interface and re-uses the redis db implementation for core-metadata and adds influx db for core-data in the database implementation.
All the core-data rest APIs are implemented and should work.

The reason of this implementation is due to some limitations of redisdb being in memory consumption lot of memory when lots of datas pumped inside core-data.

To activate influx inside core data 
1. Clone the repository and build a local image of core-data with the current code
```makefile
make ddata
```
2. Put the following changes in the Edgex compose file. 
3. For Core-data service add the following environment variables
   Please note the HYBRID: "TRUE" flag actyually tells the core-data to create a hybrid client and if this flag is not there core-data will continue to use redis as per the existing implementation
```
INFLUXDB_URL: http://influxdb:8086
INFLUXDB_ORG: TestORG
INFLUXDB_BUCKET: Edgex
HYBRID: "TRUE"
```
4. Use the local image normally tagged as
```
   image: edgexfoundry/core-data:0.0.0-dev
```
5. Add the influx image inside the compose file
```
influxdb:
    image: influxdb:2.7.0-alpine
    hostname: edgex-influx
    container_name: edgex-influx   
    networks:
      edgex-network: null
    ports:
    - mode: ingress
      target: 8086
      published: "8086"
      protocol: tcp
    volumes:
    - type: volume
      source: influxVolume
      target: /var/lib/influxdb2
      volume: {}
```
6. In the volume section add the influx volume
```
   influxVolume:
    name: edgex_influxVolume
```
7. Apart from core-data APIs itself the influx UI or influx cli can be also used to query the data stored.
   The usename is 'admin'. In secured mode the redis password is used to secure the influx db as well. So the password can be recovered from secret store as explained inside the link https://docs.edgexfoundry.org/3.0/security/Ch-SecretStore/#:~:text=is%20true.-,Using%20the%20Vault%20CLI,-Execute%20a%20shell

    In unsecured mode the password is "admin1234".
   

[Apache-2.0](LICENSE)


