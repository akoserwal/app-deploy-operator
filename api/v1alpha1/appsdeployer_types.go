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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type App struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Image     string `json:"image,omitempty"`
	Port      int32  `json:"port,omitempty"`
}

// AppsdeployerSpec defines the desired state of Appsdeployer
type AppsdeployerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Appsdeployer. Edit appsdeployer_types.go to remove/update
	Apps []App `json:"apps,omitempty"`
}

// AppsdeployerStatus defines the observed state of Appsdeployer
type AppsdeployerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Appsdeployer is the Schema for the appsdeployers API
type Appsdeployer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppsdeployerSpec   `json:"spec,omitempty"`
	Status AppsdeployerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AppsdeployerList contains a list of Appsdeployer
type AppsdeployerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Appsdeployer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Appsdeployer{}, &AppsdeployerList{})
}
