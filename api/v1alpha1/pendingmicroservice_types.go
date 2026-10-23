/*
Copyright 2026 magic-artifactory.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PendingMicroserviceSpec defines the desired state of PendingMicroservice.
type PendingMicroserviceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Charts string `json:"charts_url"`
	Uuid   string `json:"uuid"`
}

// PendingMicroserviceStatus defines the observed state of PendingMicroservice.
type PendingMicroserviceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories:brw;all;brw-pending;pending-ms // TODO: FIX CATEGORIES DOESN'T WORK
// +kubebuilder:printcolumn:name="UUID",type="string",JSONPath=".spec.uuid",description="UUID"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// PendingMicroservice is the Schema for the pendingmicroservices API.
type PendingMicroservice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PendingMicroserviceSpec   `json:"spec,omitempty"`
	Status PendingMicroserviceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PendingMicroserviceList contains a list of PendingMicroservice.
type PendingMicroserviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PendingMicroservice `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PendingMicroservice{}, &PendingMicroserviceList{})
}
