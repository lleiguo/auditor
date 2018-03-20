package k8s

import "strings"

type Container struct {
	Spec ContainerSpec
}

func (container Container) Audit() (bool, string) {
	if container.Spec.ImagePullPolicy != "IfNotPresent" {
		return false, "ImagePullPolicy is not set to `IfNotPresent`"
	}
	if strings.Contains(container.Spec.Image, "latest") {
		return false, "Image using latest tag"
	}
	if len(container.Spec.Resources.Limits.Memory) == 0 && len(container.Spec.Resources.Requests.Memory) == 0 {
		return false, "Missing memory setting"
	}
	if container.Spec.SecurityContext.RunAsNonRoot == true {
		return false, "Container running as non-root"
	}
	if len(container.Spec.Lifecycle.PreStop.Exec.Command) > 0 {
		return false, "Container using PreStop hooks"
	}
	return false, "Unknown reason"
}

type ContainerSpec struct {
	Name            string   `json:"name"`
	Image           string   `json:"image"`
	ImagePullPolicy string   `json:"imagePullPolicy"`
	Args            []string `json:"args,omitempty"`
	Env []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"env"`
	Resources struct {
		Limits struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"limits"`
		Requests struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"requests"`
	} `json:"resources"`
	SecurityContext struct {
		RunAsNonRoot bool `json:"runAsNonRoot"`
	} `json:"securityContext,omitempty"`
	Lifecycle struct {
		PreStop struct {
			Exec struct {
				Command []string `json:"command"`
			} `json:"exec"`
		} `json:"preStop"`
	} `json:"lifecycle,omitempty"`
}

type SplunkForwarder struct {
	Container Container
}

type Consul struct {
	Container Container
}

type XMService struct {
	Container Container
}

func (xmService XMService) clone(spec ContainerSpec) {
	xmService.Container.Spec = spec
}

func (splunkForwarder SplunkForwarder) Audit() (bool, string) {
	commonSettings, reason := splunkForwarder.Container.Audit()

	for _, env := range splunkForwarder.Container.Spec.Env {
		switch  env.Name {
		case "SPLUNK_DEPLOYMENT_SERVER":
			if env.Value != "splunkdeploymentserver.i.xmatters.com:8089" {
				return false, "SPLUNK_DEPLOYMENT_SERVER != splunkdeploymentserver.i.xmatters.com:8089"
			}
		case "SPLUNK_START_ARGS":
			if env.Value != "--accept-license --answer-yes" {
				return false, "SPLUNK_START_ARGS != --accept-license --answer-yes"
			}
		case "SPLUNK_USER":
			if env.Value != "root" {
				return false, "SPLUNK_USER != root"
			}
		case "SPLUNK_FORWARD_SERVER_1":
			if env.Value != "10.26.101.79:9997" {
				return false, "SPLUNK_FORWARD_SERVER_1 != 10.26.101.79:9997"
			}
		case "SPLUNK_FORWARD_SERVER_2":
			if env.Value != "10.26.101.80:9997" {
				return false, "SPLUNK_FORWARD_SERVER_1 != 10.26.101.80:9997"
			}
		}
	}
	return commonSettings && splunkForwarder.Container.Spec.Resources.Limits.Memory == "512Mi" && splunkForwarder.Container.Spec.Resources.Requests.Memory == "256Mi", reason
}

func (consul Consul) Audit() (bool, string) {
	commonSettings, reason := consul.Container.Audit()

	//Need to check the args
	//args := []string{agent -advertise=$(POD_IP) -bind=0.0.0.0 -retry-join=$(CONSUL_CLUSTER) -datacenter=$(CONSUL_DC) -disable-host-node-id}
	for _, env := range consul.Container.Spec.Env {

		switch env.Name {
		case "CONSUL_LOCAL_CONFIG":
			//return "{   "service": {     "checks": [       {         "interval": "10s",         "http": "http://localhost:8888/ping",         "timeout": "1s"       }     ";
		case "port":
			if env.Value != "8888" {
				return false, "port != 8888"
			}
		case "tags":
			//["version=5-5-202", "5-5-202"]
		case "CONSUL_DC":
			//us-central1-tst
		case "CONSUL_CLUSTER":
			if env.Value != "consul.service.consul" {
				return false, "CONSUL_CLUSTER != consul.service.consul"
			}
		}
	}

	return commonSettings && consul.Container.Spec.Resources.Limits.Memory == "256Mi" && consul.Container.Spec.Resources.Requests.Memory == "128Mi", reason
}

func (xmService XMService) Audit() (bool, string) {
	commonSettings, reason := xmService.Container.Audit()
	switch xmService.Container.Spec.Name {
	case "xmatters-eng-mgmt-billing":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "" && xmService.Container.Spec.Resources.Requests.Memory == "", reason
	case "xmatters-eng-mgmt-customerconfig":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "1Gi" && xmService.Container.Spec.Resources.Requests.Memory == "512Mi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-dbjobsequencer":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "16Gi" && xmService.Container.Spec.Resources.Requests.Memory == "12Gi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-hyrax":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "4Gi" && xmService.Container.Spec.Resources.Requests.Memory == "2Gi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-mobileapi":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "1Gi" && xmService.Container.Spec.Resources.Requests.Memory == "1Gi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-multinode":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "2Gi" && xmService.Container.Spec.Resources.Requests.Memory == "2Gi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-reapi":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "6656Mi" && xmService.Container.Spec.Resources.Requests.Memory == "3328Mi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-resolution":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "2Gi" && xmService.Container.Spec.Resources.Requests.Memory == "1Gi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-scheduler":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "3Gi" && xmService.Container.Spec.Resources.Requests.Memory == "2Gi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-soap":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "6656Mi" && xmService.Container.Spec.Resources.Requests.Memory == "3328Mi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-voicexml":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "6656Mi" && xmService.Container.Spec.Resources.Requests.Memory == "3328Mi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-webui":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "6656Mi" && xmService.Container.Spec.Resources.Requests.Memory == "3328Mi", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-xerus":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "" && xmService.Container.Spec.Resources.Requests.Memory == "", "CPU/MEM Setting"
	case "xmatters-eng-mgmt-xmapi":
		return commonSettings && xmService.Container.Spec.Resources.Limits.Memory == "3Gi" && xmService.Container.Spec.Resources.Requests.Memory == "2Gi", "CPU/MEM Setting"
	}
	return false, "Unknown reason"

}
