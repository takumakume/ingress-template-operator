package render

import (
	"bytes"
	"text/template"

	networkingv1 "k8s.io/api/networking/v1"
)

type Options struct {
	Namespace string
}

func (opt *Options) ToMap() map[string]string {
	return map[string]string{
		"namespace": opt.Namespace,
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
	}

	return ing, nil
}

type renderer struct {
	data map[string]string
}

func newRenderer(data map[string]string) *renderer {
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
