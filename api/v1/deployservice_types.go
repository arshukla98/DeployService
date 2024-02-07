package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// DeployService is a specification for a DeployService resource
type DeployService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeployServiceSpec   `json:"spec"`
	Status DeployServiceStatus `json:"status,omitempty"`
}

// DeployServiceSpec is the spec for a DeployService resource
type DeployServiceSpec struct {
	DepTemp DeploymentTemplate `json:"deploymentTemplate"`
	ServiceTemp ServiceTemplate `json:"serviceTemplate"`
}

type DeploymentTemplate struct{
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas int32 `json:"replicas"`
	ImageName string `json:"imageName"`
}

type ServiceTemplate struct{
	Name string `json:"name"`
	Type string `json:"type"`
	ServicePort int32 `json:"servicePort"`
}

// DeployServiceStatus is the status for a DeployService resource
type DeployServiceStatus struct {
	DepCreated bool `json:"DepCreated"`
	DepImage string `json:"DepImage"`
	DepReplicasCount int32 `json:"DepReplicasCount"`
	SvcCreated bool `json:"SvcCreated"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
// DeployServiceList is a list of DeployService resources
type DeployServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []DeployService `json:"items"`
}
