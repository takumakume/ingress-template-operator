package render

import (
	"reflect"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOptions_ToMap(t *testing.T) {
	type fields struct {
		Namespace string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "default",
			fields: fields{
				Namespace: "hoge",
			},
			want: map[string]string{
				"namespace": "hoge",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Options{
				Namespace: tt.fields.Namespace,
			}
			if got := opt.ToMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_renderer_render(t *testing.T) {
	type fields struct {
		data map[string]string
	}
	type args struct {
		tmpl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			args: args{
				tmpl: "{{ .key1 }}-{{ .key2 }}",
			},
			want: "value1-value2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &renderer{
				data: tt.fields.data,
			}
			got, err := r.render(tt.args.tmpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderer.render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("renderer.render() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRender(t *testing.T) {
	type args struct {
		ing *networkingv1.Ingress
		opt Options
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
				ing: &networkingv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-ingress",
						Annotations: map[string]string{
							"annotation/key": "value-{{ .namespace }}",
						},
						Labels: map[string]string{
							"key": "value-{{ .namespace }}",
						},
					},
					Spec: networkingv1.IngressSpec{
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
				opt: Options{
					Namespace: "hoge",
				},
			},
			want: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-ingress",
					Annotations: map[string]string{
						"annotation/key": "value-hoge",
					},
					Labels: map[string]string{
						"key": "value-hoge",
					},
				},
				Spec: networkingv1.IngressSpec{
					TLS: []networkingv1.IngressTLS{
						{
							Hosts: []string{
								"hoge.example.com",
							},
						},
					},
					Rules: []networkingv1.IngressRule{
						{
							Host: "hoge.example.com",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.args.ing, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
