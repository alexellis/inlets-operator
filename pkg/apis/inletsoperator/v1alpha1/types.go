/*
Copyright 2017 The Kubernetes Authors.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:printcolumn:name="Service",type=string,JSONPath=`.spec.serviceName`
// +kubebuilder:printcolumn:name="Tunnel",type=string,JSONPath=`.spec.client_deployment.name`
// +kubebuilder:printcolumn:name="HostStatus",type=string,JSONPath=`.status.hostStatus`
// +kubebuilder:printcolumn:name="HostIP",type=string,JSONPath=`.status.hostIP`
// +kubebuilder:printcolumn:name="HostID",type=string,JSONPath=`.status.hostId`
// Tunnel is a specification for a Tunnel resource
type Tunnel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TunnelSpec   `json:"spec"`
	Status TunnelStatus `json:"status"`
}

// TunnelSpec is the spec for a Tunnel resource
type TunnelSpec struct {
	ServiceName string `json:"serviceName"`

	ClientDeploymentRef *metav1.ObjectMeta `json:"client_deployment"`
	AuthToken           string             `json:"auth_token"`
}

// TunnelStatus is the status for a Tunnel resource
type TunnelStatus struct {
	HostStatus string `json:"hostStatus"`
	HostIP     string `json:"hostIP"`
	HostID     string `json:"hostId"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TunnelList is a list of Tunnel resources
type TunnelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Tunnel `json:"items"`
}
