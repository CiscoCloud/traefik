## Traefik Metrics in Kafka
#### Overview
This repository contains updates to Traefik supporting writing metrics to Kafka.  Metrics are in [Monasco format](https://wiki.openstack.org/wiki/Monasca/Message_Schema#Metrics_Message).  A sample metric is:

    { 
      "metric": {
        "name":"elapsedMs-traefik",
        "dimensions": {
          "SHIPPED_ENVIRONMENT_ID":"70fe6d48-0622-11e6-bdad-0242ac110003",
          "SHIPPED_ENVIRONMENT_NAME":"david-staging",
          "SHIPPED_PROJECT_ID":"19c60964-0621-11e6-bd9b-0242ac110003",
          "SHIPPED_PROJECT_NAME":"traefik-kafka-3",
          "SHIPPED_SERVICE_ID":"3ab53f68-0621-11e6-bd9f-0242ac110003",
          "SHIPPED_SERVICE_NAME":"traefik-kafka-cmx-3",
        },
        "timestamp":1461198662000,
        "value":161
      },
      "meta":{
        "tenantId":"19c60964-0621-11e6-bd9b-0242ac110003",
        "region":""
      },
      "creation_time":1461198662000
    }

These follow the layout described in [cAdvisor documentation](https://github.com/CiscoCloud/shipped-monitoring/blob/monitoring/docs/cAdvisor.md#monasca-metrics-format), except without containerId, containerName, and hostName, which are not available to Traefik (see [Metric Dimensions](#dimensions))

#### Configuration
Metric output is configured in traefik.toml.  The `accessLogsFile` configuration parameter is deprecated (though still supported) and replaced by an [accessLog] group:

    [accessLog]
      Filename = "/var/log/traefik/access.log"
      Brokers = "173.39.246.177:31092"
      Topic = "traefik"

where:

**Filename** specifies the path to a log file where Traefik writes the status of each query.   Specifying an empty string for Filename suppresses writing to the log.  Log output takes the form:

    73.193.20.224 - - [21/Apr/2016:12:36:54 +0000] "GET /api/location/v1/history/clients/00:00:2a:01:00:30 HTTP/1.1" 200 1166221 "http://david-staging--traefik-kafka-3--traefik-kafka-cmx-3--807e8c.tx3.shipped-cisco.com/" "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36" 3 "david-staging--traefik-kafka-3--traefik-kafka-cmx-3--807e8c" "http://shipped-tx3-worker-103.node.consul:12687" 300.997468ms

where the last four tokens (request number, frontend name, backend name, and elapsed time) are added by this update.  This portion of the update adding these tokens has been submitted to the Traefik master as a pull request.

**Brokers** specifies a comma-separated list of Kafka brokers in the form *hostname:port*. Specifying an empty string for Brokers suppresses writing Kafka metrics.

**Topic** specifies the topic for Kafka metrics.  This is normally "metrics".

#### <a name="dimensions"></a>Metric Dimensions
Metric dimensions allow filtering metrics in queries.  The initial implementation includes six dimensions as noted above.  Other dimensions could be added if supported in flexdb.  

Traefik has access to information in Marathon labels on the frontend and to the values reported to the access log.  Any of this information could be added to metrics.

Field Name | Source | Description | Example
---------- | ------ | ----------- | -------
agent | Request | Requesting agent | Chrome/49.0.2623.112 Safari/537.36
backend URL | Request | Backend assigned | http://shipped-tx3-worker-103.node.consul:12687
config_id | Marathon | Shipped configuration UUID | c5d02d75-05b9-11e6-89ec-0242ac110003
deploy_target_id | Marathon | Shipped deploy target UUID | 31f03bc9-0314-11e5-b9c3-6c4008ad584e
deploy_target_name | Marathon | Shipped deploy target name | TX3
elapsed time | Request | Request time in nanoseconds | 3848332
env_id | Marathon | Shipped environment UUID | 70fe6d48-0622-11e6-bdad-0242ac110003
env_name | Marathon | Shipped environment name | david-staging
frontend name | Request | Frontend name | david-staging--traefik-kafka-3--traefik-kafka-cmx-3--807e8c
HAPROXY_HTTP | Marathon | | true/false
host | Request | Requesting host | 73.193.20.224
method | Request | HTTP method | GET
HTTP_PORT_IDX_0_NAME | Marathon | | david-staging--traefik-kafka-3--traefik-kafka-cmx-3--807e8c
project_id | Marathon | Shipped project UUID | 19c60964-0621-11e6-bd9b-0242ac110003
project_name | Marathon | Shipped project name | traefik-kafka-3
protocol | Request | HTTP protocol | HTTP/1.1
referer | Request | HTTP referrer |
request id | Request | Sequential request number | 123
service_id | Marathon | Shipped service UUID | 3ab53f68-0621-11e6-bd9f-0242ac110003
service_name | Marathon | Shipped service name | traefik-kafka-cmx-3
size | Request | Request size | 1166221
status | Request | HTTP status code | 200
timestamp | Request | Request Unix time | 1461198662000
URI | Request | Request URI | /api/location/v1/history/clients/00:00:2a:01:00:30
username | Request | Request user
