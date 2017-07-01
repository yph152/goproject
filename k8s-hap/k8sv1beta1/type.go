package k8sv1beta1

type K8sScaleSpec struct {
	Replicas             int32   `json:"replicas,omitempty"`
	MinReplicas          int32   `json:"minReplicas,omitempty"`
	MaxReplicas          int32   `json:"maxReplicas,omitempty"`
	Name                 string  `json:"name,omitempty"`
	Namespace            string  `json:"namespace,omitempty"`
	TargetPercentage     float64 `json:"targetPercentage,omitempty"`
	MinTargetPercentage  float64 `json:"minTargetPercentage,omitempty"`
	ForCPUUtilization    bool    `json:"forcpuutilization,omitempty"`
	ForMemoryUtilization bool    `json:"formemoryutilization,omitempty"`
	CurrentPercentage    float64 `json:"currentPercentage,omitempty"`
	CurrentReplicas      int32   `json:"currentReplicas,omitempty"`
}
