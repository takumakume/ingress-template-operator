package render

import (
	"bytes"
	"text/template"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Metadata metav1.ObjectMeta
}

func (opt *Options) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Metadata": opt.Metadata,
	}
}

func Render(ing *networkingv1.Ingress, opt Options) (*networkingv1.Ingress, error) {
	r := newRenderer(opt.ToMap())

	for k, v := range ing.Labels {
		if ret, err := r.render(v); err == nil {
			ing.Labels[k] = ret
		} else {
			return nil, err
		}
	}

	for k, v := range ing.Annotations {
		if ret, err := r.render(v); err == nil {
			ing.Annotations[k] = ret
		} else {
			return nil, err
		}
	}

	for i, tls := range ing.Spec.TLS {
		secretName, err := r.render(tls.SecretName)
		if err != nil {
			return nil, err
		}
		ing.Spec.TLS[i].SecretName = secretName

		for ii, host := range tls.Hosts {
			if ret, err := r.render(host); err == nil {
				ing.Spec.TLS[i].Hosts[ii] = ret
			} else {
				return nil, err
			}
		}
	}

	for i, rule := range ing.Spec.Rules {
		if ret, err := r.render(rule.Host); err == nil {
			ing.Spec.Rules[i].Host = ret
		} else {
			return nil, err
		}

		if rule.HTTP != nil && len(rule.HTTP.Paths) > 0 {
			for ii, path := range rule.HTTP.Paths {
				if path.Backend.Resource != nil && path.Backend.Resource.Name != "" {
					if ret, err := r.render(path.Backend.Resource.Name); err == nil {
						ing.Spec.Rules[i].HTTP.Paths[ii].Backend.Resource.Name = ret
					} else {
						return nil, err
					}
				}
				if path.Backend.Service != nil && path.Backend.Service.Name != "" {
					if ret, err := r.render(path.Backend.Service.Name); err == nil {
						ing.Spec.Rules[i].HTTP.Paths[ii].Backend.Service.Name = ret
					} else {
						return nil, err
					}
				}
			}
		}

	}

	return ing, nil
}

type renderer struct {
	data map[string]interface{}
}

func newRenderer(data map[string]interface{}) *renderer {
	return &renderer{
		data: data,
	}
}

func (r *renderer) render(tmpl string) (string, error) {
	tpl, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, r.data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
