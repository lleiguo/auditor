package main

import (
	"bufio"
	"com/hs/k8s"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var namespace = "default"
var allNamespaces = true
var passed = "https://www.katalon.com/wp-content/themes/katalon/template-parts/page/features/img/supported-icon.png?ver=17.11.07"
var failed = "http://www.vetriias.com/images/Deep_Close.png"

func main() {

	if len(os.Args) > 4 {
		fmt.Fprintf(os.Stderr, "Too many arguments, auditor takes one namespace or [-a || -all] for all namespaces , e.g. 'auditor' or 'auditor -n default' or 'auditor -a , actual arguments %d \n", len(os.Args))
		os.Exit(1)
	} else if len(os.Args) == 4 {
		if os.Args[1] == "-a" || os.Args[1] == "-all" {
			allNamespaces = true
			namespace = "--all-namespaces"
		} else if os.Args[1] == "-n" {
			namespace = os.Args[1]
		}
		if os.Args[2] == "-n" {
			namespace = os.Args[2]
		} else if os.Args[2] == "-a" || os.Args[2] == "-all" {
			allNamespaces = true
			namespace = "--all-namespaces"
		}
	} else if len(os.Args) == 2 {
		if os.Args[1] == "-a" || os.Args[1] == "-all" {
			allNamespaces = true
			namespace = "--all-namespaces"
		} else if os.Args[1] == "-n" {
			namespace = os.Args[1]
		}
	}

	createFile()
	f, _ := os.OpenFile("./"+namespace+"_service_configuration.html", os.O_APPEND|os.O_RDWR, 0644)
	writer := bufio.NewWriter(f)
	defer f.Close()

	header := "<!DOCTYPE html><html><head><meta name='viewport' content='width=device-width, initial-scale=0.5'><link rel='stylesheet' href='https://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css'> <script src='https://code.jquery.com/jquery-1.11.3.min.js'></script><script src='https://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js'></script><script src='https://www.kryogenix.org/code/browser/sorttable/sorttable.js'></script> </head> <body>"
	table := "<table class='sortable table table-bordered'; data-resizable-columns-id='demo-table-v2'><thread><tbody><tr align='center'>" +
		"<th data-resizable-column-id='service'><strong>Deployed Service</strong></th>" +
		"<th data-resizable-column-id='description'><strong>Description</strong></th>" +
		"<th data-resizable-column-id='team'><strong>Team</strong></th>" +
		"<th data-resizable-column-id='pager team'><strong>Pager Team</strong></th>" +
		"<th data-resizable-column-id='skeleton'><strong>Skeleton Type</strong></th>" +
		"<th data-resizable-column-id='slack'><strong>Slack</strong></th>" +
		"<th data-resizable-column-id='github'><strong>GitHub</strong></th>" +
		"<th data-resizable-column-id='maintainer'><strong>Maintainer(s)</strong></th>" +
		"<th data-resizable-column-id='sensu'><strong>Sensu Checks</strong></th>" +
		"<th data-resizable-column-id='service type'><strong>Service Type</strong></th>" +
		"<th data-resizable-column-id='labels'><strong>Labels</strong></th>" +
		"<th data-resizable-column-id='replica'><strong>Replica</strong></th>" +
		"<th data-resizable-column-id='sumologic'><strong>Sumologic</strong></th>" +
		"<th data-resizable-column-id='resource limits'><strong>Resource Limits</strong></th>" +
		"<th data-resizable-column-id='resource requests'><strong>Resource Requests</strong></th>" +
		"<th data-resizable-column-id='liveness'><strong>Liveness</strong></th>" +
		"<th data-resizable-column-id='readiness'><strong>Readiness</strong></th></tr></thread>"
	fmt.Fprintln(writer, header, table)
	writer.Flush()

	getServiceDescription("default", writer)
	footer := "</tbody></table></body></html>"
	fmt.Fprintln(writer, footer)
	writer.Flush()
}

func getServiceDescription(namespace string, writer *bufio.Writer) {
	cmd := exec.Command("kubectl", "get", "deployment", "-o", "json", "-n", namespace)
	deploy, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)

	cmd = exec.Command("kubectl", "get", "service", "-o", "json", "-n", namespace)
	svc, err := cmd.CombinedOutput()
	printCommand(cmd)
	printError(err)

	parseServiceDescription(deploy, svc, writer, namespace)
}

func parseServiceDescription(deploy []byte, svc []byte, writer *bufio.Writer, service string) {

	var serviceDescriptor k8s.Deployment
	err := json.Unmarshal(deploy, &serviceDescriptor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	k8sService := k8s.Service{}
	err = json.Unmarshal(svc, &k8sService)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// lb, lbType := k8sService.Audit()

	for _, item := range serviceDescriptor.Items {
		if !strings.Contains(item.Metadata.Name, "monitoring") {

			fmt.Fprintln(writer, fmt.Sprintf("<tr align='center'><td><strong>%s</strong></td>", item.Metadata.Name))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Metadata.Annotations.HootsuiteComDescription))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Metadata.Annotations.HootsuiteComTeam))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Metadata.Annotations.HootsuiteComPagerTeam))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Metadata.Annotations.HootsuiteComSkeletonType))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left><a href=https://hootsuite.slack.com/messages/%s>%s</a></td>", item.Metadata.Annotations.HootsuiteComSlackChannel, item.Metadata.Annotations.HootsuiteComSlackChannel))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left><a href=%s>%s</a></td>", item.Metadata.Annotations.HootsuiteComGithub, item.Metadata.Annotations.HootsuiteComGithub))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Metadata.Annotations.HootsuiteComMaintainers))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Metadata.Annotations.HootsuiteComSensuChecks))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Spec.Selector.MatchLabels.ServiceType))

			labels := strings.Replace(fmt.Sprintf("%+v", item.Metadata.Labels), " ", "<BR>", -1)
			labels = strings.Replace(labels, "{", "", -1)
			labels = strings.Replace(labels, "}", "", -1)
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", labels))
			writeTD(item.Status.AvailableReplicas >= 1 && item.Status.ReadyReplicas >= 1, writer, fmt.Sprintf("%s - [%+v]", strconv.Itoa(item.Status.AvailableReplicas), strconv.Itoa(item.Status.ReadyReplicas)))
			fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%s</td>", item.Spec.Template.Metadata.Annotations.SumologicComInclude))

			// var serviceContainer

			// for _, containerSpec := range item.Spec.Template.Spec.Containers {
			// 	switch containerSpec.Name {
			// 	default:
			// 		serviceContainer = k8s.XMService{k8s.Container{containerSpec}}
			// 	}
			// }
			// fmt.Fprintln(writer, fmt.Sprintf("<td align=left>%+v</td>", serviceContainer))

			// var result bool
			// var reason string
			// result, reason = serviceContainer.Audit()
			// writeTD(result, writer, fmt.Sprintf("%s - [%+v]", reason, serviceContainer))
			// fmt.Fprintln(writer, fmt.Sprintf("</tr>"))
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
	fullPath := "./" + namespace + "_service_configuration.html"
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
