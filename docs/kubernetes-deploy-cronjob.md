

# Goal

Example deployment of kubent as a cronjob periodally scanning kubernetes
cluster and publishing result in prometheus format. 


# Implementation

-   `kubent` run as kubernetes cronjob
-   result of `kubent` scan stored in json file
-   json file transformed into prometheus metrics format using `jq`
-   prometheus metrics pushed using curl to prometheus pushgateway


# Installation

    kubectl apply -f ./manifests/kubent-cronjob.yaml -n kubent-system

Assuming prometheus [pushgateway](https://github.com/prometheus/pushgateway) is installed, configured and scraped by some prometheus instance.
For installing the pushgateway in kubernetes see for example [this](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-pushgateway) helm chart.
Example grafana dashboard representing scraped data included in [./dashboards](dashboards/).

