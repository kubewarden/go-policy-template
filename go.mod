module github.com/kubewarden/go-policy-template

go 1.20

replace github.com/go-openapi/strfmt => github.com/kubewarden/strfmt v0.1.2

require (
	github.com/francoispqt/onelog v0.0.0-20190306043706-8c2bb31b10a4
	github.com/kubewarden/k8s-objects v1.26.0-kw2
	github.com/kubewarden/policy-sdk-go v0.4.1
	github.com/mailru/easyjson v0.7.7
	github.com/wapc/wapc-guest-tinygo v0.3.3
)

require (
	github.com/francoispqt/gojay v0.0.0-20181220093123-f2cc13a668ca // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/josharian/intern v1.0.0 // indirect
)
