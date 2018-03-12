package main

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"log"
	"bufio"
	//"reflect"
)

var service string
var allServices = false
var xmServices = []string{"billing", "customerconfig", "dbjobsequencer", "hyrax", "mobileapi", "multinode", "reapi", "resolution", "scheduler", "soap", "voicexml", "webui", "xerus", "xmapi",}
var checkmark = "https://www.katalon.com/wp-content/themes/katalon/template-parts/page/features/img/supported-icon.png?ver=17.11.07"
var failed = "http://www.vetriias.com/images/Deep_Close.png"

type Container struct {
	Name string `json:"name"`
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
	Lifecycle struct {
		PreStop struct {
			Exec struct {
				Command []string `json:"command"`
			} `json:"exec"`
		} `json:"preStop"`
	} `json:"lifecycle,omitempty"`
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

var splunkForwarder = "{Name:xmatters-eng-mgmt-xmsplunkforwarder Resources:{Limits:{Memory:512Mi} Requests:{Memory:256Mi}} SecurityContext:{RunAsNonRoot:false} Lifecycle:{PreStop:{Exec:{Command:[]}}}}"

var consul = "{Name:xmatters-eng-mgmt-xmconsul Resources:{Limits:{Memory:256Mi} Requests:{Memory:128Mi}} SecurityContext:{RunAsNonRoot:false} Lifecycle:{PreStop:{Exec:{Command:[]}}}}"

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
	output, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)
	parseServiceDescription(output, writer, service)
}

func parseServiceDescription(serviceDescription []byte, writer *bufio.Writer, service string) {

	var serviceDescriptor ReplicaSet
	err := json.Unmarshal(serviceDescription, &serviceDescriptor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, item := range serviceDescriptor.Items {
		if !strings.Contains(item.Metadata.Name, "monitoring") {

			fmt.Fprintln(writer, fmt.Sprintf("<tr align='center'><td><strong>%s</strong></td>", item.Metadata.Name))
			if len(item.Metadata.OwnerReferences) == 1 {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", item.Metadata.OwnerReferences[0].Kind, failed))
			} else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", "Replica Set", checkmark))
			}

			if item.Spec.Replicas == 3 {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%d' src=%s width='32' height='32'></td>", item.Spec.Replicas, checkmark))
			} else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%d' src=%s width='32' height='32'></td>", item.Spec.Replicas, failed))
			}

			labels := strings.Replace(fmt.Sprintf("%+v", item.Metadata.Labels), " ", "<BR>", -1)
			labels = strings.Replace(labels, "{", "", -1)
			labels = strings.Replace(labels, "}", "", -1)
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", labels))

			if item.Spec.Template.Spec.DNSPolicy == "ClusterFirst" {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%s width='32' height='32'></td>", item.Spec.Template.Spec.DNSPolicy, checkmark))
			} else {
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
			} else {
				fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%s' src=%d width='32' height='32'></td>", fmt.Sprintf("%+v", item.Spec.Template.Spec.TerminationGracePeriodSeconds), failed))
			}

			for _, container := range item.Spec.Template.Spec.Containers {
				if container.Name == "xmatters-eng-mgmt-xmsplunkforwarder" {
					if fmt.Sprintf("%+v", container) == splunkForwarder {
						fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%+v' src=%s width='32' height='32'></td>", fmt.Sprintf("%+v", container), checkmark))
					} else {
						fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%+v' src=%s width='32' height='32'></td>", fmt.Sprintf("%+v", container), failed))
					}
				}
			}
			for _, container := range item.Spec.Template.Spec.Containers {
				if container.Name == "xmatters-eng-mgmt-xmconsul" {
					if fmt.Sprintf("%+v", container) == consul {
						fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%+v' src=%s width='32' height='32'></td>", fmt.Sprintf("%+v", container), checkmark))
					} else {
						fmt.Fprintln(writer, fmt.Sprintf("<td><img border='0' title='%+v' src=%s width='32' height='32'></td>", fmt.Sprintf("%+v", container), failed))
					}
				}
			}
			for _, container := range item.Spec.Template.Spec.Containers {
				if container.Name == "xmatters-eng-mgmt-"+service {
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
