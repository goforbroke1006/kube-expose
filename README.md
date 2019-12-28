# kube-expose

Tool for exposing kubernetes cluster services ports to localhost.

Inspired by docker-composer.

Draft version.

## Config sample

```yaml
services:

  data-aggregator-svc:
    ports:
      - 9292:8080

  ws-broadcaster-svc:
    ports:
      - 18080:8080

```

where **data-aggregator-svc**  and **ws-broadcaster-svc** are "name"-labels for services from k8s deployment files.
