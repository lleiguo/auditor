package main

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"log"
	"bufio"
	"strconv"
)

var service string
var allServices = false
var xmServices = []string{"billing", "customerconfig", "dbjobsequencer", "hyrax", "mobileapi", "multinode", "reapi", "resolution", "scheduler", "soap", "voicexml", "webui", "xerus", "xmapi"}
var passed = "https://www.katalon.com/wp-content/themes/katalon/template-parts/page/features/img/supported-icon.png?ver=17.11.07"
var failed = "http://www.vetriias.com/images/Deep_Close.png"

type Container struct {
	Name string `json:"name"`
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

type Service struct {
	Items []struct {
		Spec struct {
			ClusterIP       string `json:"clusterIP"`
			SessionAffinity string `json:"sessionAffinity"`
			Type            string `json:"type"`
		} `json:"spec"`
	} `json:"items"`
}

func (container Container) audit() bool {
	commonSettings := len(container.Resources.Limits.Memory) > 0 && len(container.Resources.Requests.Memory) > 0 && container.SecurityContext.RunAsNonRoot == false && len(container.Lifecycle.PreStop.Exec.Command) == 0
	switch container.Name {
	case "xmatters-eng-mgmt-xmsplunkforwarder":
		return commonSettings && container.Resources.Limits.Memory == "512Mi" && container.Resources.Requests.Memory == "256Mi"
	case "xmatters-eng-mgmt-xmconsul":
		return commonSettings && container.Resources.Limits.Memory == "256Mi" && container.Resources.Requests.Memory == "128Mi"
	case "xmatters-eng-mgmt-billing":
		return commonSettings && container.Resources.Limits.Memory == "" && container.Resources.Requests.Memory == ""
	case "xmatters-eng-mgmt-customerconfig":
		return commonSettings && container.Resources.Limits.Memory == "1Gi" && container.Resources.Requests.Memory == "512Mi"
	case "xmatters-eng-mgmt-dbjobsequencer":
		return commonSettings && container.Resources.Limits.Memory == "16Gi" && container.Resources.Requests.Memory == "12Gi"
	case "xmatters-eng-mgmt-hyrax":
		return commonSettings && container.Resources.Limits.Memory == "4Gi" && container.Resources.Requests.Memory == "2Gi"
	case "xmatters-eng-mgmt-mobileapi":
		return commonSettings && container.Resources.Limits.Memory == "1Gi" && container.Resources.Requests.Memory == "1Gi"
	case "xmatters-eng-mgmt-multinode":
		return commonSettings && container.Resources.Limits.Memory == "2Gi" && container.Resources.Requests.Memory == "2Gi"
	case "xmatters-eng-mgmt-reapi":
		return commonSettings && container.Resources.Limits.Memory == "6656Mi" && container.Resources.Requests.Memory == "3328Mi"
	case "xmatters-eng-mgmt-resolution":
		return commonSettings && container.Resources.Limits.Memory == "2Gi" && container.Resources.Requests.Memory == "1Gi"
	case "xmatters-eng-mgmt-scheduler":
		return commonSettings && container.Resources.Limits.Memory == "3Gi" && container.Resources.Requests.Memory == "2Gi"
	case "xmatters-eng-mgmt-soap":
		return commonSettings && container.Resources.Limits.Memory == "6656Mi" && container.Resources.Requests.Memory == "3328Mi"
	case "xmatters-eng-mgmt-voicexml":
		return commonSettings && container.Resources.Limits.Memory == "6656Mi" && container.Resources.Requests.Memory == "3328Mi"
	case "xmatters-eng-mgmt-webui":
		return commonSettings && container.Resources.Limits.Memory == "6656Mi" && container.Resources.Requests.Memory == "3328Mi"
	case "xmatters-eng-mgmt-xerus":
		return commonSettings && container.Resources.Limits.Memory == "" && container.Resources.Requests.Memory == ""
	case "xmatters-eng-mgmt-xmapi":
		return commonSettings && container.Resources.Limits.Memory == "3Gi" && container.Resources.Requests.Memory == "2Gi"
	}
	return false
}

func (service Service) audit() (bool, string) {
	if len(service.Items) > 0 {
		return service.Items[0].Spec.Type == "ClusterIP", service.Items[0].Spec.Type
	}
	return false, "No Load Balance Found"
}

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
					Containers                    []Container
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
		"<td style='width: 200px;'><strong>Service RS</strong></td>" +
		"<td style='width: 57px;'><strong>Kind (Replica Set)</strong></td>" +
		"<td style='width: 57px;'><strong>Load Balancer (ClusterIP)</strong></td>" +
		"<td style='width: 100px;'><strong>Replica Count (Dev: 3; TST: 3</strong></td>" +
		"<td style='width: 300px;'><strong>Labels</strong></td>" +
		"<td style='width: 57px;'><strong>DNS Policy (ClusterFirst)</strong></td>" +
		"<td style='width: 57px;'><strong>Volumes (xmatters-logs)</strong></td>" +
		"<td style='width: 57px;'><strong>Termination Grace Period (30s)</strong></td>" +
		"<td style='width: 57px;'><strong>Splunk Forwarder</strong></td>" +
		"<td style='width: 57px;'><strong>Consul</strong></td>" +
		"<td style='width: 57px;'><strong>Service Container</strong></td></tr>"
	fmt.Fprintln(writer, header, table)
	writer.Flush()

	if allServices {
		for _, service := range xmServices {
			getServiceDescription(service, writer)
		}
	} else {
		getServiceDescription(service, writer)
	}
	footer := "</tbody></table></body></html>"
	fmt.Fprintln(writer, footer)
	writer.Flush()
}

func getServiceDescription(service string, writer *bufio.Writer) {
	cmd := exec.Command("kubectl", "get", "rs", "-o", "json", "-n", service)
	rs, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)

	cmd = exec.Command("kubectl", "get", "service", "-o", "json", "-n", service)
	svc, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)

	parseServiceDescription(rs, svc, writer, service)
}

func parseServiceDescription(rs []byte, svc []byte, writer *bufio.Writer, service string) {

	var serviceDescriptor ReplicaSet
	err := json.Unmarshal(rs, &serviceDescriptor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var k8sService Service
	err = json.Unmarshal(svc, &k8sService)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	lb, lbType := k8sService.audit()

	for _, item := range serviceDescriptor.Items {
		if !strings.Contains(item.Metadata.Name, "monitoring") {

			fmt.Fprintln(writer, fmt.Sprintf("<tr align='center'><td><strong>%s</strong></td>", item.Metadata.Name))

			kind := "Replica Set"
			if len(item.Metadata.OwnerReferences) > 0 {
				kind = item.Metadata.OwnerReferences[0].Kind
			}
			writeTD(len(item.Metadata.OwnerReferences) == 0, writer, kind)
			writeTD(lb, writer, lbType)

			writeTD(item.Spec.Replicas == 3, writer, strconv.Itoa(item.Spec.Replicas))

			labels := strings.Replace(fmt.Sprintf("%+v", item.Metadata.Labels), " ", "<BR>", -1)
			labels = strings.Replace(labels, "{", "", -1)
			labels = strings.Replace(labels, "}", "", -1)
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", labels))

			writeTD(item.Spec.Template.Spec.DNSPolicy == "ClusterFirst", writer, fmt.Sprintf("%+v", item.Spec.Template.Spec.DNSPolicy))

			var xmLogsVolumeExists = false
			for _, volume := range item.Spec.Template.Spec.Volumes {
				if volume.Name == "xmatters-logs" {
					xmLogsVolumeExists = true
				}
			}

			writeTD(xmLogsVolumeExists, writer, fmt.Sprintf("%+v", item.Spec.Template.Spec.Volumes))

			writeTD(item.Spec.Template.Spec.TerminationGracePeriodSeconds == 30, writer, strconv.Itoa(item.Spec.Template.Spec.TerminationGracePeriodSeconds))

			var splunkContainer, consulContainer, serviceContainer Container

			for _, container := range item.Spec.Template.Spec.Containers {
				switch container.Name {
				case "xmatters-eng-mgmt-xmsplunkforwarder":
					splunkContainer = container
				case "xmatters-eng-mgmt-xmconsul":
					consulContainer = container
				default:
					if container.Name == "xmatters-eng-mgmt-"+service {
						serviceContainer = container
					}
				}
			}
			writeTD(splunkContainer.audit(), writer, fmt.Sprintf("%+v", splunkContainer))
			writeTD(consulContainer.audit(), writer, fmt.Sprintf("%+v", consulContainer))
			writeTD(serviceContainer.audit(), writer, fmt.Sprintf("%+v", serviceContainer))
			fmt.Fprintln(writer, fmt.Sprintf("</tr>"))
			writer.Flush()
		}
	}
}

func writeTD(pass bool, writer *bufio.Writer, title string) {
	if pass {
		fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", title, passed))
	} else {
		fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", title, failed))
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

	return err != nil
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}
