package main

import (
	"os/exec"
	"os"
	"fmt"
	"encoding/json"
	"log"
	"bufio"
	"strconv"
	"com/hs/k8s"
)

var service string
var allServices = true
var passed = "https://www.katalon.com/wp-content/themes/katalon/template-parts/page/features/img/supported-icon.png?ver=17.11.07"
var failed = "http://www.vetriias.com/images/Deep_Close.png"
var region = "dev"

func main() {

	if len(os.Args) > 4 || len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Too many arguments, auditor takes one [service] or [-a || -all] and/or [-r || -region] , e.g. 'auditor xmapi' or 'auditor -a' or 'auditor -a -r active, actual arguments %d \n", len(os.Args))
		os.Exit(1)
	} else if len(os.Args) == 4 {
		if os.Args[1] == "-a" || os.Args[1] == "-all" {
			allServices = true
		} else {
			service = os.Args[1]
		}
		if os.Args[2] == "-r" || os.Args[2] == "-region" {
			region = os.Args[3]
		}
	}else if len(os.Args) == 2 {
		if os.Args[1] == "-a" || os.Args[1] == "-all" {
			allServices = true
		} else {
			service = os.Args[1]
		}
	}

	createFile()
	f, _ := os.OpenFile("./" + region + "_service_configuration.html", os.O_APPEND|os.O_RDWR, 0644)
	writer := bufio.NewWriter(f)
	defer f.Close()

	header := "<!DOCTYPE html><html><head><meta name='viewport' content='width=device-width, initial-scale=0.5'><link rel='stylesheet' href='https://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css'> <script src='https://code.jquery.com/jquery-1.11.3.min.js'></script> <script src='https://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js'></script> </head> <body>"
	table := "<table style=margin: 0px auto; border='1'; align='centre'><tbody><tr align='center'><td colspan=11><strong>%s</strong></td></tr><tr align='center'>" +
		"<td style='width: 200px;'><strong>Service RS</strong></td>" +
		"<td style='width: 57px;'><strong>Kind (Replica Set)</strong></td>" +
		"<td style='width: 57px;'><strong>Load Balancer (ClusterIP)</strong></td>" +
		"<td style='width: 100px;'><strong>Replica Count (Dev/Passive: 1; TST: 3</strong></td>" +
		"<td style='width: 300px;'><strong>Labels</strong></td>" +
		"<td style='width: 57px;'><strong>DNS Policy (ClusterFirst)</strong></td>" +
		"<td style='width: 57px;'><strong>Volumes (xmatters-logs)</strong></td>" +
		"<td style='width: 57px;'><strong>Termination Grace Period (30s)</strong></td>" +
		"<td style='width: 57px;'><strong>Splunk Forwarder</strong></td>" +
		"<td style='width: 57px;'><strong>Consul</strong></td>" +
		"<td style='width: 57px;'><strong>Service Container</strong></td></tr>"
	fmt.Fprintln(writer, header, table)
	writer.Flush()

	getServiceDescription("default", writer)
	footer := "</tbody></table></body></html>"
	fmt.Fprintln(writer, footer)
	writer.Flush()
}

func getServiceDescription(namespace string, writer *bufio.Writer) {
	cmd := exec.Command("kubectl", "get", "rs", "-o", "json", "-n", namespace)
	rs, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)

	cmd = exec.Command("kubectl", "get", "service", "-o", "json", "-n", namespace)
	svc, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)

	parseServiceDescription(rs, svc, writer, namespace)
}

func parseServiceDescription(rs []byte, svc []byte, writer *bufio.Writer, service string) {

	var serviceDescriptor k8s.ReplicaSet
	err := json.Unmarshal(rs, &serviceDescriptor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	k8sService := k8s.Service{}
	err = json.Unmarshal(svc, &k8sService)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	lb, lbType := k8sService.Audit()

	for _, item := range serviceDescriptor.Items {
		if !strings.Contains(item.Metadata.Name, "monitoring") {

			fmt.Fprintln(writer, fmt.Sprintf("<tr align='center'><td><strong>%s</strong></td>", item.Metadata.Name))

			kind := "Replica Set"
			if len(item.Metadata.OwnerReferences) > 0 {
				kind = item.Metadata.OwnerReferences[0].Kind
			}
			writeTD(len(item.Metadata.OwnerReferences) == 0, writer, kind)
			writeTD(lb, writer, lbType)

			var replica = 1
			if region == "active" {
				//If passive or dev, replica will be 1
				replica = 3
			}
			writeTD(item.Spec.Replicas == replica, writer, strconv.Itoa(item.Spec.Replicas))

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

			var serviceContainer k8s.XMService

			for _, containerSpec := range item.Spec.Template.Spec.Containers {
				switch containerSpec.Name {
				default:
					serviceContainer = k8s.XMService{k8s.Container{containerSpec}}
				}
			}
			result, reason = serviceContainer.Audit()
			writeTD(result, writer, fmt.Sprintf("%s - [%+v]", reason, serviceContainer))
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
	fullPath := "./" + region + "_service_configuration.html"
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