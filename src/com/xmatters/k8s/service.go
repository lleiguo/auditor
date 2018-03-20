package k8s

type Service struct {
	Items []struct {
		Spec struct {
			ClusterIP       string `json:"clusterIP"`
			SessionAffinity string `json:"sessionAffinity"`
			Type            string `json:"type"`
		} `json:"spec"`
	} `json:"items"`
}


func (service Service) Audit() (bool, string) {
	if len(service.Items) > 0 {
		return service.Items[0].Spec.Type == "ClusterIP", service.Items[0].Spec.Type
	}
	return false, "No Load Balance Found"
}