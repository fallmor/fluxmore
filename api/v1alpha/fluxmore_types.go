/*
Copyright 2024.

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

package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FluxMoreSpec defines the desired state of FluxMore.
type FluxMoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// helmreleasename is the name of the HelmRealease to suspend
	HelmReleaseName string `json:"helmreleasename,omitempty"`
	// resourcescheck is the resource that need to exist before the HelmRelease spec suspend is set to false
	ResourcesCheck string `json:"resourcescheck,omitempty"`
	// +kubebuilder:validation:Enum=secret;configmap;deploy
	// when deploy, it will check if the ReplicaNumber is equel to the running pod
	ResourcesKind string `json:"resourceskind,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	ReplicaNumber *int32 `json:"replicanumber,omitempty"`

	Namespace string `json:"namespace,omitempty"`
}

// FluxMoreStatus defines the observed state of FluxMore.
type FluxMoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ResourcesCheckFound bool `json:"resourceCheckfound,omitempty"`

	HelmReleasePatched bool `json:"helmReleasePatched,omitempty"`

	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// FluxMore is the Schema for the fluxmores API.
type FluxMore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluxMoreSpec   `json:"spec,omitempty"`
	Status FluxMoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FluxMoreList contains a list of FluxMore.
type FluxMoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FluxMore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FluxMore{}, &FluxMoreList{})
}
