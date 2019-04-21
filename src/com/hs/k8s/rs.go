package k8s

type ReplicaSet struct {
	Items []struct {
		Kind string `json:"kind"`
		Metadata struct {
			Annotations struct {
				DeploymentKubernetesIoDesiredReplicas string `json:"deployment.kubernetes.io/desired-replicas"`
				DeploymentKubernetesIoMaxReplicas     string `json:"deployment.kubernetes.io/max-replicas"`
			} `json:"annotations"`
			Labels struct {
				App     string `json:"app"`
				Cluster string `json:"cluster"`
				Detail  string `json:"detail"`
				Stack   string `json:"stack"`
			} `json:"labels"`
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			OwnerReferences []struct {
				Kind string `json:"kind"`
				Name string `json:"name"`
			} `json:"ownerReferences"`
		} `json:"metadata"`
		Spec struct {
			Replicas int `json:"replicas"`
			Template struct {
				Metadata struct {
					Annotations struct {
						PrometheusIoScrape string `json:"prometheus.io/scrape"`
					} `json:"annotations"`
					Labels struct {
						App     string `json:"app"`
						Cluster string `json:"cluster"`
						Detail  string `json:"detail"`
						Stack   string `json:"stack"`
					} `json:"labels"`
				} `json:"metadata"`
				Spec struct {
					Containers                    []ContainerSpec
					DNSPolicy                     string `json:"dnsPolicy"`
					RestartPolicy                 string `json:"restartPolicy"`
					TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
					Volumes []struct {
						EmptyDir struct {
						} `json:"emptyDir"`
						Name string `json:"name"`
					} `json:"volumes"`
				} `json:"spec"`
			} `json:"template"`
		} `json:"spec"`
	} `json:"items"`
	Kind string `json:"kind"`
}