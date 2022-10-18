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
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ingresstemplatev1alpha1 "github.com/takumakume/ingress-template-operator/api/v1alpha1"
)

func Test_ingressTemplateToIngress(t *testing.T) {
	type args struct {
		ingresstemplate *ingresstemplatev1alpha1.IngressTemplate
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
				ingresstemplate: &ingresstemplatev1alpha1.IngressTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "ns",
					},
					Spec: ingresstemplatev1alpha1.IngressTemplateSpec{
						IngressAnnotations: map[string]string{
							"key1": "value1-{{ .Metadata.Namespace }}",
						},
						IngressLabels: map[string]string{
							"key2": "value2-{{ .Metadata.Namespace }}",
						},
						IngressSpecTemplate: networkingv1.IngressSpec{
							TLS: []networkingv1.IngressTLS{
								{
									Hosts: []string{
										"{{ .Metadata.Namespace }}.example.com",
									},
								},
							},
							Rules: []networkingv1.IngressRule{
								{
									Host: "{{ .Metadata.Namespace }}.example.com",
								},
							},
						},
					},
				},
			},
			want: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "ns",
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

var _ = Describe("IngressTemplate controller", func() {
	BeforeEach(func() {
		err := k8sClient.DeleteAllOf(ctx, &ingresstemplatev1alpha1.IngressTemplate{}, client.InNamespace("test"))
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &networkingv1.Ingress{}, client.InNamespace("test"))
		Expect(err).NotTo(HaveOccurred())
		time.Sleep(100 * time.Millisecond)
	})

	It("aa", func() {
		ingresstemplate := &ingresstemplatev1alpha1.IngressTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample",
				Namespace: "test",
			},
			Spec: ingresstemplatev1alpha1.IngressTemplateSpec{
				IngressAnnotations: map[string]string{
					"key1": "value1-{{ .Metadata.Namespace }}",
				},
				IngressLabels: map[string]string{
					"key2": "value2-{{ .Metadata.Namespace }}",
				},
				IngressSpecTemplate: networkingv1.IngressSpec{
					TLS: []networkingv1.IngressTLS{
						{
							Hosts: []string{
								"{{ .Metadata.Namespace }}.example.com",
							},
						},
					},
					Rules: []networkingv1.IngressRule{
						{
							Host: "{{ .Metadata.Namespace }}.example.com",
						},
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, ingresstemplate)).Should(Succeed())

		Eventually(func() (v1.ConditionStatus, error) {
			o := &ingresstemplatev1alpha1.IngressTemplate{}

			err := k8sClient.Get(ctx, client.ObjectKey{Namespace: "test", Name: "sample"}, o)
			if err != nil {
				return "", err
			}
			return o.Status.Ready, nil
		}, 20, 1).Should(Equal(v1.ConditionTrue))

		created := &networkingv1.Ingress{}
		Expect(k8sClient.Get(ctx, client.ObjectKey{Namespace: "test", Name: "sample"}, created)).Should(Succeed())
		Expect(created.ObjectMeta.Name).Should(Equal("sample"))
		Expect(created.ObjectMeta.OwnerReferences[0].APIVersion).Should(Equal("ingress-template.takumakume.github.io/v1alpha1"))
		Expect(created.ObjectMeta.OwnerReferences[0].Kind).Should(Equal("IngressTemplate"))
		Expect(created.ObjectMeta.OwnerReferences[0].Name).Should(Equal("sample"))
		Expect(created.ObjectMeta.OwnerReferences[0].UID).Should(Equal(ingresstemplate.GetUID()))
		Expect(created.ObjectMeta.Annotations["key1"]).Should(Equal("value1-test"))
		Expect(created.ObjectMeta.Labels["key2"]).Should(Equal("value2-test"))
		Expect(created.Spec.TLS[0].Hosts[0]).Should(Equal("test.example.com"))
		Expect(created.Spec.Rules[0].Host).Should(Equal("test.example.com"))

		Expect(k8sClient.Get(ctx, client.ObjectKey{Namespace: "test", Name: "sample"}, ingresstemplate))
		ingresstemplate.Spec.IngressSpecTemplate.TLS[0].Hosts[0] = "{{ .Metadata.Namespace }}.hoge.com"
		ingresstemplate.Spec.IngressSpecTemplate.Rules[0].Host = "{{ .Metadata.Namespace }}.hoge.com"
		Expect(k8sClient.Update(ctx, ingresstemplate)).Should(Succeed())
		Eventually(func() error {
			o := &ingresstemplatev1alpha1.IngressTemplate{}
			err := k8sClient.Get(ctx, client.ObjectKey{Namespace: "test", Name: "sample"}, o)
			if err != nil {
				return err
			}
			if o.Spec.IngressSpecTemplate.TLS[0].Hosts[0] != "{{ .Metadata.Namespace }}.hoge.com" {
				return fmt.Errorf("IngressTemplate.Spec.IngressSpecTemplate.TLS[0].Hosts[0] has not been updated: %s", o.Spec.IngressSpecTemplate.TLS[0].Hosts[0])
			}
			if o.Spec.IngressSpecTemplate.Rules[0].Host != "{{ .Metadata.Namespace }}.hoge.com" {
				return fmt.Errorf("IngressTemplate.Spec.IngressSpecTemplate.Rules[0].Host has not been updated: %s", o.Spec.IngressSpecTemplate.Rules[0].Host)
			}
			return nil
		}, 20, 1).Should(Succeed())
		Eventually(func() error {
			o := &networkingv1.Ingress{}
			err := k8sClient.Get(ctx, client.ObjectKey{Namespace: "test", Name: "sample"}, o)
			if err != nil {
				return err
			}
			if o.Spec.TLS[0].Hosts[0] != "test.hoge.com" {
				return fmt.Errorf("Ingress.Spec.TLS[0].Hosts[0] has not been updated: %s", o.Spec.TLS[0].Hosts[0])
			}
			if o.Spec.Rules[0].Host != "test.haoge.com" {
				return fmt.Errorf("Ingress.Spec.Rules[0].Host has not been updated: %s", o.Spec.Rules[0].Host)
			}
			return nil
		}, 20, 1).Should(Succeed())
	})
})
