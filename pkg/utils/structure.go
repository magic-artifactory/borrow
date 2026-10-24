package utils

type SpecPendingResource struct {
	Charts      string `json:"charts_url"`
	ProjectName string `json:"project_name"`
	Uuid        string `json:"uuid"`
}

type ResourceSpec struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
	// CreationTimestamp  `json:"name"`
	Namespace string `json:"namespace"`
}
