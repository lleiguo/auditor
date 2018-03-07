**Things to consider**
- support following parameters
  - --all-namespaces
  - --namespace


Based on the namespace it passed in, it'll go through the main pod, ignore the monitoring ones for now, 
1. Run "k get pod hyrax-dev-1-186-0-4-v000-2bg75 -n hyrax -o yaml" to create the pod yaml file
2. Parse the yaml file to generate the following output:
- Annotation: 
    - prometheus.io/scrape: "true"
- Labels:
    -    app: hyrax
    -    cluster: hyrax-dev-1-186-0-4
    -    detail: 1-186-0-4
    -    hyrax-dev-1-186-0-4: "true"
    -    pod-template-hash: "549930784"
    -    replication-controller: hyrax-dev-1-186-0-4-v000
    -    stack: dev
    -    version: "0"
- ownerReferences:
    -  "kind"
- spec/containers
    - consul
    ```
      - args:
        - agent
        - -advertise=$(POD_IP)
        - -bind=0.0.0.0
        - -retry-join=$(CONSUL_CLUSTER)
        - -disable-host-node-id
        - -datacenter=$(CONSUL_DC)
        - -node-meta=version:1-186-0
        env:
        - name: CONSUL_CLUSTER
          value: consul.service.consul
        - name: POD_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        - name: CONSUL_LOCAL_CONFIG
          value: '{   "service": {     "checks": [],     "name": "hyrax",     "port":
            8080,     "tags": ["version=1-186-0", "1-186-0"]   } }'
        - name: CONSUL_DC
          value: us-central1-dev
        image: gcr.io/xmatters-eng-mgmt/xm_consul:latest
        imagePullPolicy: IfNotPresent
        lifecycle:
          preStop:
            exec:
              command:
              - /bin/sh
              - -c
              - consul
              - leave
        name: xmatters-eng-mgmt-xmconsul
        ports:
        - containerPort: 8500
          name: ui-port
          protocol: TCP
        - containerPort: 8400
          name: alt-port
          protocol: TCP
        - containerPort: 53
          name: udp-port
          protocol: UDP
        - containerPort: 8443
          name: https-port
          protocol: TCP
        - containerPort: 8080
          name: http-port
          protocol: TCP
        - containerPort: 8301
          name: serfwan
          protocol: TCP
        - containerPort: 8600
          name: consuldns
          protocol: TCP
        - containerPort: 8300
          name: server
          protocol: TCP
        resources:
          limits:
            memory: 256Mi
          requests:
            memory: 128Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
          name: default-token-h7d4h
          readOnly: true

```
    - splunk    
    ```
  - env:
    - name: XM_SPLUNK_DEPLOYMENT_CLIENTNAME
      value: hyrax-dev-uscentral1-1-186-0
    - name: SPLUNK_DEPLOYMENT_SERVER
      value: splunkdeploymentserver.i.xmatters.com:8089
    - name: SPLUNK_START_ARGS
      value: --accept-license --answer-yes
    - name: SPLUNK_USER
      value: splunk
    image: gcr.io/xmatters-eng-mgmt/xm_splunkforwarder:latest
    imagePullPolicy: IfNotPresent
    name: xmatters-eng-mgmt-xmsplunkforwarder
    ports:
    - containerPort: 80
      name: http
      protocol: TCP
    resources:
      limits:
        memory: 512Mi
      requests:
        memory: 265Mi
    securityContext:
      runAsNonRoot: false
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/log/xmatters
      name: xmatters-logs
      readOnly: true
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: default-token-h7d4h
      readOnly: true
```      
    - Service container
        -     resources:
                limits:
                  memory: 4Gi
                requests:
                  memory: 4Gi  
