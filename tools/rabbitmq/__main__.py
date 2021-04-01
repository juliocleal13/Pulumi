"""A Kubernetes Python Pulumi program"""

import pulumi, base64
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts
from pulumi_kubernetes.core.v1 import Secret



config = pulumi.Config()
data = config.require_object("set")
config_k8s = pulumi.Config("kubernetes")


rabbitmq = Chart('rabbitmq-julioleal', config=ChartOpts(
    chart='rabbitmq',
    namespace=config_k8s.get("namespace"),
    fetch_opts=FetchOpts(
        repo="https://charts.bitnami.com/bitnami"
    ),
    values={
        "replicaCount": data.get("replicaCount"),
        "resources": {
            "limits": {
                "memory": data.get("resourcesLimitsMemory")
            },
        },
        "resources": {
            "limits": {
                "cpu": data.get ("resourcesLimitsCpu")
            },
        },
        "resources": {
            "requests": {
                "memory": data.get("resourcesRequestsMemory"),
            },
        },
        "resources": {
            "requests": {
                "cpu": data.get("resourcesRequestCpu"),
            },
        },
        "memoryHighWatermark": {
            "enabled": data.get("memoryHighWatermarkEnabled"),
        },   
        "memoryHighWatermark": {
            "type": data.get("memoryHighWatermarkType"),
        },
        "memoryHighWatermark": {
            "value": data.get("memoryHighWatermarkValue"),
        },
        "extraPlugins": data.get("extraPlugins"),
        "auth": {
            "password": config.require_secret("rabbitmqAuthPassword"),
        },
        "auth": {
            "username": data.get("username"),
        },
        "service": {
            "type": data.get('serviceType'),
        },
        "clustering":{
            "forceBoot": data.get("forceBoot"),
        },
        "persistence": {
            "size": data.get("persistenceSize")
        },            
    },
))



# Export the public IP for WordPress.
#frontend = mongodb.get_resource('v1/Service', 'mongodb-julioleal')
#pulumi.export('frontend_ip', frontend.status.load_balancer.ingress[0].ip)
