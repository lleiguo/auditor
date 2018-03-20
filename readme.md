# K8S Auditor

Verify the service deployed into GCP is compliant with the service deployment standard we hve
## Getting Started

Build it

```
go build src/com/xmatters/auditor/auditor.go
```

Run against all service namespaces

```
./bin/auditor -a
```

Run against one service namespace e.g. xmapi

```
./bin/auditor xmapi
```