# K8S Auditor

Verify the service deployed into GCP is compliant with the service deployment standard we hve
## Getting Started

Run against all service namespaces in dev (default)

```
go run ./src/com/xmatters/auditor/auditor.go -a
```

Run against all service namespaces in active region

```
go run ./src/com/xmatters/auditor/auditor.go -a -r active
```

Run against all service namespaces in passive region

```
go run ./src/com/xmatters/auditor/auditor.go -a -r passive
```

Run against one service namespace e.g. xmapi

```
./bin/auditor xmapi
```