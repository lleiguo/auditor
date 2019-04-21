package k8s

import "time"

type Deployment struct {
	Items []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				DeploymentKubernetesIoRevision              string `json:"deployment.kubernetes.io/revision"`
				HootsuiteComDescription                     string `json:"hootsuite.com/description"`
				HootsuiteComGithub                          string `json:"hootsuite.com/github"`
				HootsuiteComMaintainers                     string `json:"hootsuite.com/maintainers"`
				HootsuiteComPagerTeam                       string `json:"hootsuite.com/pager-team"`
				HootsuiteComSensuChecks                     string `json:"hootsuite.com/sensu-checks"`
				HootsuiteComSkeletonType                    string `json:"hootsuite.com/skeleton-type"`
				HootsuiteComSlackChannel                    string `json:"hootsuite.com/slack-channel"`
				HootsuiteComTeam                            string `json:"hootsuite.com/team"`
				KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
			} `json:"annotations"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Generation        int       `json:"generation"`
			Labels            struct {
				App         string `json:"app"`
				ArkRestore  string `json:"ark-restore"`
				ServiceType string `json:"service-type"`
			} `json:"labels"`
			Name            string `json:"name"`
			Namespace       string `json:"namespace"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			UID             string `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			ProgressDeadlineSeconds int `json:"progressDeadlineSeconds"`
			Replicas                int `json:"replicas"`
			RevisionHistoryLimit    int `json:"revisionHistoryLimit"`
			Selector                struct {
				MatchLabels struct {
					App         string `json:"app"`
					ServiceType string `json:"service-type"`
				} `json:"matchLabels"`
			} `json:"selector"`
			Strategy struct {
				Type string `json:"type"`
			} `json:"strategy"`
			Template struct {
				Metadata struct {
					Annotations struct {
						HootsuiteComSensuChecks string `json:"hootsuite.com/sensu-checks"`
						SumologicComInclude     string `json:"sumologic.com/include"`
					} `json:"annotations"`
					CreationTimestamp interface{} `json:"creationTimestamp"`
					Labels            struct {
						App         string `json:"app"`
						ServiceType string `json:"service-type"`
					} `json:"labels"`
				} `json:"metadata"`
				Spec struct {
					Containers []struct {
						Env []struct {
							Name      string `json:"name"`
							ValueFrom struct {
								FieldRef struct {
									APIVersion string `json:"apiVersion"`
									FieldPath  string `json:"fieldPath"`
								} `json:"fieldRef"`
							} `json:"valueFrom,omitempty"`
							Value string `json:"value,omitempty"`
						} `json:"env"`
						Image           string `json:"image"`
						ImagePullPolicy string `json:"imagePullPolicy"`
						LivenessProbe   struct {
							FailureThreshold int `json:"failureThreshold"`
							HTTPGet          struct {
								Path   string `json:"path"`
								Scheme string `json:"scheme"`
							} `json:"httpGet"`
							InitialDelaySeconds int `json:"initialDelaySeconds"`
							PeriodSeconds       int `json:"periodSeconds"`
							SuccessThreshold    int `json:"successThreshold"`
							TimeoutSeconds      int `json:"timeoutSeconds"`
						} `json:"livenessProbe"`
						Name  string `json:"name"`
						Ports []struct {
							ContainerPort int    `json:"containerPort"`
							Name          string `json:"name"`
							Protocol      string `json:"protocol"`
						} `json:"ports"`
						ReadinessProbe struct {
							FailureThreshold int `json:"failureThreshold"`
							HTTPGet          struct {
								Path   string `json:"path"`
								Scheme string `json:"scheme"`
							} `json:"httpGet"`
							PeriodSeconds    int `json:"periodSeconds"`
							SuccessThreshold int `json:"successThreshold"`
							TimeoutSeconds   int `json:"timeoutSeconds"`
						} `json:"readinessProbe"`
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
						TerminationMessagePath   string `json:"terminationMessagePath"`
						TerminationMessagePolicy string `json:"terminationMessagePolicy"`
						VolumeMounts             []struct {
							MountPath string `json:"mountPath"`
							Name      string `json:"name"`
						} `json:"volumeMounts"`
					} `json:"containers"`
					DNSPolicy        string `json:"dnsPolicy"`
					ImagePullSecrets []struct {
						Name string `json:"name"`
					} `json:"imagePullSecrets"`
					RestartPolicy   string `json:"restartPolicy"`
					SchedulerName   string `json:"schedulerName"`
					SecurityContext struct {
					} `json:"securityContext"`
					TerminationGracePeriodSeconds int `json:"terminationGracePeriodSeconds"`
					Volumes                       []struct {
						Name   string `json:"name"`
						Secret struct {
							DefaultMode int    `json:"defaultMode"`
							SecretName  string `json:"secretName"`
						} `json:"secret"`
					} `json:"volumes"`
				} `json:"spec"`
			} `json:"template"`
		} `json:"spec"`
		Status struct {
			AvailableReplicas int `json:"availableReplicas"`
			Conditions        []struct {
				LastTransitionTime time.Time `json:"lastTransitionTime"`
				LastUpdateTime     time.Time `json:"lastUpdateTime"`
				Message            string    `json:"message"`
				Reason             string    `json:"reason"`
				Status             string    `json:"status"`
				Type               string    `json:"type"`
			} `json:"conditions"`
			ObservedGeneration int `json:"observedGeneration"`
			ReadyReplicas      int `json:"readyReplicas"`
			Replicas           int `json:"replicas"`
			UpdatedReplicas    int `json:"updatedReplicas"`
		} `json:"status"`
	} `json:"items"`
}
