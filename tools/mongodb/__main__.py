"""A Kubernetes Python Pulumi program"""

import pulumi
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts

mongodb = Chart('mongodb-julioleal', config=ChartOpts(
    chart='mongodb',
    namespace="julioleal",
    fetch_opts=FetchOpts(
        repo="https://charts.bitnami.com/bitnami"
    ),
    values={
        "auth": {
            "enabled": False,
        },
        "persistence": {
            "size": "20Gi",
        },
    },
))

# Export the public IP for WordPress.
#frontend = mongodb.get_resource('v1/Service', 'mongodb-julioleal')
#pulumi.export('frontend_ip', frontend.status.load_balancer.ingress[0].ip)