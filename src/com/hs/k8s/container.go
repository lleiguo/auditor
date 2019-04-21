package k8s

type Container struct {
	Spec ContainerSpec
}

var commonSettings = true

func (container Container) Audit() (bool, string) {
	if container.Spec.ImagePullPolicy != "IfNotPresent" {
		return false, "ImagePullPolicy is not set to `IfNotPresent`"
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
	if container.Spec.LivenessProbe.FailureThreshold == 0 {
		return false, "Container missing liveness check"
	}
	if container.Spec.ReadinessProbe.FailureThreshold == 0 {
		return false, "Container missing readiness check"
	}
	return true, "Non-compliant CPU/MEM settings"
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
	LivenessProbe struct {
		FailureThreshold int `json:"failureThreshold"`
	} `json:"livenessprobe"`
	ReadinessProbe struct {
		FailureThreshold int `json:"failureThreshold"`
	}`json:"readinessprobe"`
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

type XMService struct {
	Container Container
}

func (xmService XMService) clone(spec ContainerSpec) {
	xmService.Container.Spec = spec
}

func (xmService XMService) Audit() (bool, string) {
	commonSettings, reason := xmService.Container.Audit()
	return commonSettings, reason + "; Memory Limit: " + xmService.Container.Spec.Resources.Limits.Memory + "; Memory Requested: " + xmService.Container.Spec.Resources.Requests.Memory
}
