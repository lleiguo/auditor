package main

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
	"time"
	"encoding/json"
	"log"
	"bufio"
)

var service string
var allServices = false
var xmServices = []string{"hyrax", "mobileapi", "reapi", "resolution", "scheduler", "soap", "voicexml", "webui", "xmapi", "multinode"}
var checkmark = "https://www.katalon.com/wp-content/themes/katalon/template-parts/page/features/img/supported-icon.png?ver=17.11.07"
var failed = "http://www.vetriias.com/images/Deep_Close.png"

type ServiceDescriptor struct {
	APIVersion string `json:"apiVersion"`
	Items []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata struct {
			Annotations struct {
				DeploymentKubernetesIoDesiredReplicas string `json:"deployment.kubernetes.io/desired-replicas"`
				DeploymentKubernetesIoMaxReplicas     string `json:"deployment.kubernetes.io/max-replicas"`
				DeploymentKubernetesIoRevision        string `json:"deployment.kubernetes.io/revision"`
			} `json:"annotations"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Generation        int       `json:"generation"`
			Labels struct {
				App                   string `json:"app"`
				Cluster               string `json:"cluster"`
				Detail                string `json:"detail"`
				Stack                 string `json:"stack"`
				Version               string `json:"version"`
			} `json:"labels"`
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			OwnerReferences []struct {
				APIVersion         string `json:"apiVersion"`
				BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
				Controller         bool   `json:"controller"`
				Kind               string `json:"kind"`
				Name               string `json:"name"`
				UID                string `json:"uid"`
			} `json:"ownerReferences"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			UID             string `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			Replicas int `json:"replicas"`
			Selector struct {
				MatchLabels struct {
					App                   string `json:"app"`
					Cluster               string `json:"cluster"`
					Detail                string `json:"detail"`
					Stack                 string `json:"stack"`
					Version               string `json:"version"`
				} `json:"matchLabels"`
			} `json:"selector"`
			Template struct {
				Metadata struct {
					Annotations struct {
						PrometheusIoScrape string `json:"prometheus.io/scrape"`
					} `json:"annotations"`
					CreationTimestamp interface{} `json:"creationTimestamp"`
					Labels struct {
						App                   string `json:"app"`
						Cluster               string `json:"cluster"`
						Detail                string `json:"detail"`
						Stack                 string `json:"stack"`
						Version               string `json:"version"`
					} `json:"labels"`
				} `json:"metadata"`
				Spec struct {
					Containers []struct {
						Env []struct {
							Name  string `json:"name"`
							Value string `json:"value"`
						} `json:"env"`
						Image           string `json:"image"`
						ImagePullPolicy string `json:"imagePullPolicy"`
						Name            string `json:"name"`
						Ports []struct {
							ContainerPort int    `json:"containerPort"`
							Name          string `json:"name"`
							Protocol      string `json:"protocol"`
						} `json:"ports"`
						Resources struct {
							Limits struct {
								Memory string `json:"memory"`
							} `json:"limits"`
							Requests struct {
								Memory string `json:"memory"`
							} `json:"requests"`
						} `json:"resources"`
						SecurityContext struct {
							RunAsNonRoot bool `json:"runAsNonRoot"`
						} `json:"securityContext,omitempty"`
						TerminationMessagePath   string `json:"terminationMessagePath"`
						TerminationMessagePolicy string `json:"terminationMessagePolicy"`
						VolumeMounts []struct {
							MountPath string `json:"mountPath"`
							Name      string `json:"name"`
							ReadOnly  bool   `json:"readOnly"`
						} `json:"volumeMounts,omitempty"`
						Args []string `json:"args,omitempty"`
						Lifecycle struct {
							PreStop struct {
								Exec struct {
									Command []string `json:"command"`
								} `json:"exec"`
							} `json:"preStop"`
						} `json:"lifecycle,omitempty"`
					} `json:"containers"`
					DNSPolicy string `json:"dnsPolicy"`
					ImagePullSecrets []struct {
						Name string `json:"name"`
					} `json:"imagePullSecrets"`
					RestartPolicy string `json:"restartPolicy"`
					SchedulerName string `json:"schedulerName"`
					SecurityContext struct {
					} `json:"securityContext"`
					TerminationGracePeriodSeconds int `json:"terminationGracePeriodSeconds"`
					Volumes []struct {
						EmptyDir struct {
						} `json:"emptyDir"`
						Name string `json:"name"`
					} `json:"volumes"`
				} `json:"spec"`
			} `json:"template"`
		} `json:"spec"`
		Status struct {
			AvailableReplicas    int `json:"availableReplicas"`
			FullyLabeledReplicas int `json:"fullyLabeledReplicas"`
			ObservedGeneration   int `json:"observedGeneration"`
			ReadyReplicas        int `json:"readyReplicas"`
			Replicas             int `json:"replicas"`
		} `json:"status"`
	} `json:"items"`
	Kind string `json:"kind"`
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
}

func main() {

	if len(os.Args) > 2 || len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Too many arguments, auditor takes one [service] or [-a || -all], e.g. 'auditor xmapi' or 'auditor -a', actual arguments %d \n", len(os.Args))
		os.Exit(1)
	} else if len(os.Args) == 2 {
		if os.Args[1] == "-a" || os.Args[1] == "-all" {
			allServices = true
		} else {
			service = os.Args[1]
		}
	}

	createFile()
	f, _ := os.OpenFile("./service_configuration.html", os.O_APPEND|os.O_RDWR, 0644)
	writer := bufio.NewWriter(f)
	defer f.Close()

	header := "<!DOCTYPE html><html><head><meta name='viewport' content='width=device-width, initial-scale=0.5'><link rel='stylesheet' href='https://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css'> <script src='https://code.jquery.com/jquery-1.11.3.min.js'></script> <script src='https://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js'></script> </head> <body>"
	table := "<table style=margin: 0px auto; border='1'; align='centre'><tbody><tr align='center'>" +
		"<td style='width: 200px;'><strong>Service</strong></td>" +
		"<td style='width: 57px;'><strong>Kind (Deployment)</strong></td>" +
		"<td style='width: 57px;'><strong>Replica Count (Dev: 3; TST: 3</strong></td>" +
		"<td style='width: 57px;'><strong>Annotation (DeploymentKubernetesIoDesiredReplicas:3 DeploymentKubernetesIoMaxReplicas:4 )</strong></td>" +
		"<td style='width: 250px;'><strong>Labels</strong></td>" +
		"<td style='width: 57px;'><strong>DNS Policy (ClusterFirst)</strong></td>" +
		"<td style='width: 57px;'><strong>Volumes</strong></td>" +
		"<td style='width: 57px;'><strong>Termination Grace Period (30s)</strong></td>" +
		"<td style='width: 57px;'><strong>Splunk Forwarder</strong></td>" +
		"<td style='width: 57px;'><strong>Consul</strong></td>" +
		"<td style='width: 57px;'><strong>Service Container</strong></td></tr>"
	fmt.Fprintln(writer, header, table)
	writer.Flush()

	if allServices {
		for _, service := range xmServices {
			getServiceDescription(service, writer, f)
		}
	} else {
		getServiceDescription(service, writer, f)
	}
	footer := "</tbody></table></body></html>"
	fmt.Fprintln(writer, footer)
	writer.Flush()
}

func getServiceDescription(service string, writer *bufio.Writer, f *os.File ) {
	cmd := exec.Command("kubectl", "get", "rs", "-o", "json", "-n", service)
	output, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)
	parseServiceDescription(output, writer, f, service)
}

func parseServiceDescription(serviceDescription []byte, writer *bufio.Writer, f *os.File, service string) {

	var serviceDescriptor ServiceDescriptor
	err := json.Unmarshal(serviceDescription, &serviceDescriptor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, item := range serviceDescriptor.Items {
		if !strings.Contains(item.Metadata.Name, "monitoring") {

			fmt.Fprintln(writer, fmt.Sprintf("<tr align='center'><td>%s</td>", item.Metadata.Name))
			if len(item.Metadata.OwnerReferences) == 1 {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", item.Metadata.OwnerReferences[0].Kind, checkmark))
			}else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", "Replica Set", failed))
			}

			if item.Spec.Replicas == 3 {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%d' src=%s width='32' height='32'></td>", item.Spec.Replicas, checkmark))
			}else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%d' src=%s width='32' height='32'></td>", item.Spec.Replicas, failed))
			}

			if item.Metadata.Annotations.DeploymentKubernetesIoDesiredReplicas == "3" && item.Metadata.Annotations.DeploymentKubernetesIoMaxReplicas == "4" {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", fmt.Sprintf("%+v", item.Metadata.Annotations), checkmark))
			} else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", fmt.Sprintf("%+v", item.Metadata.Annotations), failed))
			}

			fmt.Fprintln(writer, fmt.Sprintf("<td>%+v</td>", item.Metadata.Labels))

			if item.Spec.Template.Spec.DNSPolicy == "ClusterFirst" {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", item.Spec.Template.Spec.DNSPolicy, checkmark))
			}else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", item.Spec.Template.Spec.DNSPolicy, failed))
			}

			var xmLogsVolumeExists = false
			for _, volume := range item.Spec.Template.Spec.Volumes {
				if volume.Name == "xmatters-logs" {
					xmLogsVolumeExists = true
				}
			}

			if xmLogsVolumeExists {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='xmatters-logs volume exist' src=%s width='32' height='32'></td>", checkmark))
			} else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", fmt.Sprintf("%+v", item.Spec.Template.Spec.Volumes), failed))
			}

			if item.Spec.Template.Spec.TerminationGracePeriodSeconds == 30 {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%d' src=%s width='32' height='32'></td>", item.Spec.Template.Spec.TerminationGracePeriodSeconds, checkmark))
			}else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%d width='32' height='32'></td>", fmt.Sprintf("%+v", item.Spec.Template.Spec.TerminationGracePeriodSeconds), failed))
			}

			for _, container := range item.Spec.Template.Spec.Containers {
				if container.Name == "xmatters-eng-mgmt-xmsplunkforwarder" {
					fmt.Fprintln(writer, fmt.Sprintf("<td>%+v</td>", container))
				}
			}
			for _, container := range item.Spec.Template.Spec.Containers {
				if container.Name == "xmatters-eng-mgmt-xmconsul" {
					fmt.Fprintln(writer, fmt.Sprintf("<td>%+v</td>", container))
				}
			}
			for _, container := range item.Spec.Template.Spec.Containers {
				if container.Name == "xmatters-eng-mgmt-" +service {
					fmt.Fprintln(writer, fmt.Sprintf("<td>%+v</td>", container))
				}
			}
		}
		fmt.Fprintln(writer, fmt.Sprintf("</tr>"))
		writer.Flush()
	}

}

func createFile() {
	// detect if file exists
	fullPath := "./service_configuration.html"
	var _, err = os.Stat(fullPath)

	// create file if not exists
	if !os.IsNotExist(err) {
		e := os.Remove(fullPath)
		if isError(e) {
			return
		}

		fmt.Println("==> done deleting file")
	}
	var file, e = os.Create(fullPath)
	if isError(e) {
		return
	}
	defer file.Close()
	fmt.Println("==> done creating file", fullPath)
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}