package render

import (
	"reflect"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_renderer_render(t *testing.T) {
	type fields struct {
		data map[string]interface{}
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
				data: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			args: args{
				tmpl: "{{ .key1 }}-{{ .key2 }}",
			},
			want: "value1-value2",
		},
		{
			name: "struct",
			fields: fields{
				data: map[string]interface{}{
					"key1": metav1.ObjectMeta{
						Namespace: "hoge",
					},
				},
			},
			args: args{
				tmpl: "{{ .key1.Namespace }}",
			},
			want: "hoge",
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
						Name: "test",
						Annotations: map[string]string{
							"annotation/key": "value-{{ .Metadata.Namespace }}",
						},
						Labels: map[string]string{
							"key": "value-{{ .Metadata.Namespace }}",
						},
					},
					Spec: networkingv1.IngressSpec{
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
				opt: Options{
					Metadata: metav1.ObjectMeta{
						Namespace: "hoge",
					},
				},
			},
			want: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
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

func TestOptions_ToMap(t *testing.T) {
	type fields struct {
		Metadata metav1.ObjectMeta
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "default",
			fields: fields{
				Metadata: metav1.ObjectMeta{
					Namespace: "hoge",
				},
			},
			want: map[string]interface{}{
				"Metadata": metav1.ObjectMeta{
					Namespace: "hoge",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Options{
				Metadata: tt.fields.Metadata,
			}
			if got := opt.ToMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
