# Service Monitor 
Service monitor is a solution in golang designed to run on a Kubernetes Cluster to monitor internet urls and provide prometheus metrics.

## Assumption
1. Internet URLs not reachable or taking longer that 5 seconds to respond are down.
2. Numbers of URLs to monitor are relatively small and can be passed as configuration to the application.
