package utils

type PendingResource struct {
	Charts      string `json:"charts_url"`
	ProjectName string `json:"project_name"`
	Uid         string `json:"uid"`
}

type ResourceSpec struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
	// CreationTimestamp  `json:"name"`
	Namespace string `json:"namespace"`
}
