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

package controllers

import (
	"reflect"
	"testing"

	"github.com/takumakume/ingress-template-operator/api/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ingressTemplateToIngress(t *testing.T) {
	type args struct {
		ingresstemplate *v1alpha1.IngressTemplate
	}
	tests := []struct {
		name    string
		args    args
		want    *networkingv1.Ingress
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				ingresstemplate: &v1alpha1.IngressTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ing",
						Namespace: "ns",
					},
					Spec: v1alpha1.IngressTemplateSpec{
						IngressAnnotations: map[string]string{
							"key1": "value1-{{ .namespace }}",
						},
						IngressLabels: map[string]string{
							"key2": "value2-{{ .namespace }}",
						},
						IngressSpecTemplate: networkingv1.IngressSpec{
							TLS: []networkingv1.IngressTLS{
								{
									Hosts: []string{
										"{{ .namespace }}.example.com",
									},
								},
							},
							Rules: []networkingv1.IngressRule{
								{
									Host: "{{ .namespace }}.example.com",
								},
							},
						},
					},
				},
			},
			want: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "ing",
					Namespace:    "ns",
					Annotations: map[string]string{
						"key1": "value1-ns",
					},
					Labels: map[string]string{
						"key2": "value2-ns",
					},
				},
				Spec: networkingv1.IngressSpec{
					TLS: []networkingv1.IngressTLS{
						{
							Hosts: []string{
								"ns.example.com",
							},
						},
					},
					Rules: []networkingv1.IngressRule{
						{
							Host: "ns.example.com",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ingressTemplateToIngress(tt.args.ingresstemplate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ingressTemplateToIngress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ingressTemplateToIngress() = %v, want %v", got, tt.want)
			}
		})
	}
}
