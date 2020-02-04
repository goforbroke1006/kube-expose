# kube-expose

Tool for exposing kubernetes cluster services ports to localhost.

Inspired by docker-composer.

Draft version.

## Config sample

```yaml
namespaces:

  - name: default
    resources:
      
      data-aggregator-svc:
        type: service
        ports:
          - 9292:8080
      
      ws-broadcaster-svc:
        type: service
        ports:
          - 18080:8080

  - name: kubernetes-dashboard
    resources:
      
      kubernetes-dashboard:
        type: pod
        ports:
          - 8443:8443

```

where **data-aggregator-svc**  and **ws-broadcaster-svc** are "name"-labels for services from k8s deployment files.
