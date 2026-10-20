package utils

type PendingRessource struct {
	Charts      string `json:"charts"`
	Branch      string `json:"branch"`
	ProjectName string `json:"project_name"`
	Uid         string `json:"uid"`
}

type RessourceSpec struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
	// CreationTimestamp  `json:"name"`
	Namespace string `json:"namespace"`
}
