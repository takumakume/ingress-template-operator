# ingress-template-operator

Generate an Ingress with a golang template.

Deploy `IngressTemplate`
  ```
  apiVersion: ingress-template.takumakume.github.io/v1alpha1
  kind: IngressTemplate
  metadata:
    name: example
    namespace: hoge
  spec:
    ingressAnnotations:
      cert-manager.io/cluster-issuer: example-com-issuer
      key: "This namespace is {{ .Metadata.Namespace }}"
    ingressLabels:
      key: "This namespace is {{ .Metadata.Namespace }}"
    ingressSpecTemplate:
      tls:
      - hosts:
          - "www-{{ .Metadata.Namespace }}.example.com"
        secretName: example-com-tls
      rules:
      - host: "www-{{ .Metadata.Namespace }}.example.com"
        http:
          paths:
          - backend:
              service:
                name: example
                port:
                  number: 80
            path: /
            pathType: Prefix
  ```

Generate `Ingress`
  ```
  apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    name: example
    namespace: hoge
    annotations:
      cert-manager.io/cluster-issuer: example-com-issuer
      key: "This namespace is hoge"
    labels:
      key: "This namespace is hoge"
  spec:
    tls:
    - hosts:
        - "www-hoge.example.com"
      secretName: example-com-tls
    rules:
    - host: "www-hoge.example.com"
      http:
        paths:
        - backend:
            service:
              name: example
              port:
                number: 80
          path: /
          pathType: Prefix
  ```

The variable `.` used in the template is the IngressTemplate resource itself.

For example, if you need the namespace where the IngressTemplate is deployed, you can access it like `.Metadata.Namespace`.
