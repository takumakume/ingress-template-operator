/*
Copyright 2022.

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
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IngressTemplateSpec defines the desired state of IngressTemplate
type IngressTemplateSpec struct {
	// IngressSpec Template for Ingress.Spec
	// +kubebuilder:validation:Required
	IngressSpecTemplate networkingv1.IngressSpec `json:"ingressSpecTemplate"`

	// IngressName This name is generated in Ingress
	// +optional
	IngressName string `json:"ingressName,omitempty"`

	// Annotations This annotation is generated in Ingress
	// +optional
	IngressAnnotations map[string]string `json:"ingressAnnotations,omitempty"`

	// Labels This labels is generated in Ingress
	// +optional
	IngressLabels map[string]string `json:"ingressLabels,omitempty"`
}

// IngressTemplateStatus defines the observed state of IngressTemplate
type IngressTemplateStatus struct {
	// Ready Ingress generation status
	Ready corev1.ConditionStatus `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IngressTemplate is the Schema for the ingresstemplates API
type IngressTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngressTemplateSpec   `json:"spec,omitempty"`
	Status IngressTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IngressTemplateList contains a list of IngressTemplate
type IngressTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngressTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IngressTemplate{}, &IngressTemplateList{})
}
