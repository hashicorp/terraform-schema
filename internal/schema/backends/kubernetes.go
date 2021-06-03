package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func kubernetesBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.15.0/backend/remote-state/kubernetes/backend.go
	// Docs:
	// https://github.com/hashicorp/terraform/blob/v0.15.0/website/docs/language/settings/backends/kubernetes.html.md
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/kubernetes.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Kubernetes secret"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"secret_suffix": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("Suffix used when creating the secret. The secret will be named in the format: `tfstate-{workspace}-{secret_suffix}`."),
			},
			"labels": {
				Expr: schema.ExprConstraints{
					schema.MapExpr{Elem: schema.LiteralTypeOnly(cty.String)},
				},
				IsOptional:  true,
				Description: lang.Markdown("Map of additional labels to be applied to the secret."),
			},
			"namespace": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Namespace to store the secret in."),
			},
			"in_cluster_config": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Used to authenticate to the cluster from inside a pod."),
			},
			"load_config_file": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Load local kubeconfig."),
			},
			"host": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The hostname (in form of URI) of Kubernetes master."),
			},
			"username": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The username to use for HTTP basic authentication when accessing the Kubernetes master endpoint."),
			},
			"password": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The password to use for HTTP basic authentication when accessing the Kubernetes master endpoint."),
			},
			"insecure": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Whether server should be accessed without verifying the TLS certificate."),
			},
			"client_certificate": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("PEM-encoded client certificate for TLS authentication."),
			},
			"client_key": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("PEM-encoded client certificate key for TLS authentication."),
			},
			"cluster_ca_certificate": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("PEM-encoded root certificates bundle for TLS authentication."),
			},
			"config_path": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Path to the kube config file, defaults to ~/.kube/config"),
			},
			"config_context": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Context to choose from the config file."),
			},
			"config_context_auth_info": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Authentication info context of the kube config (name of the kubeconfig user, `--user` flag in `kubectl`)"),
			},
			"config_context_cluster": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Cluster context of the kube config (name of the kubeconfig cluster, `--cluster` flag in `kubectl`)"),
			},
			"token": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Token to authentifcate a service account."),
			},
			"exec": {
				Expr: schema.ExprConstraints{
					schema.ObjectExpr{
						Attributes: schema.ObjectExprAttributes{
							"api_version": {
								Expr:        schema.LiteralTypeOnly(cty.String),
								IsRequired:  true,
								Description: lang.Markdown("API version to use when decoding the ExecCredentials resource, e.g. `client.authentication.k8s.io/v1beta1`."),
							},
							"command": {
								Expr:        schema.LiteralTypeOnly(cty.String),
								IsRequired:  true,
								Description: lang.Markdown("Command to execute"),
							},
							"env": {
								Expr: schema.ExprConstraints{
									schema.MapExpr{Elem: schema.LiteralTypeOnly(cty.String)},
								},
								IsOptional:  true,
								Description: lang.Markdown("List of arguments to pass when executing the plugin."),
							},
							"args": {
								Expr: schema.ExprConstraints{
									schema.ListExpr{Elem: schema.LiteralTypeOnly(cty.String)},
								},
								IsOptional:  true,
								Description: lang.Markdown("Map of environment variables to set when executing the plugin."),
							},
						},
					},
				},
				IsOptional: true,
				Description: lang.Markdown("Configuration for an [exec-based credential plugin]" +
					"(https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins)," +
					" e.g. call an external command to receive user credentials."),
			},
		},
	}

	return bodySchema
}
