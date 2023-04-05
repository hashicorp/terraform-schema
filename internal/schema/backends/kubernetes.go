// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Suffix used when creating the secret. The secret will be named in the format: `tfstate-{workspace}-{secret_suffix}`."),
			},
			"labels": {
				Constraint: schema.Map{
					Elem: schema.LiteralType{Type: cty.String},
				},
				IsOptional:  true,
				Description: lang.Markdown("Map of additional labels to be applied to the secret."),
			},
			"namespace": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Namespace to store the secret in."),
			},
			"in_cluster_config": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Used to authenticate to the cluster from inside a pod."),
			},
			"load_config_file": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Load local kubeconfig."),
			},
			"host": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The hostname (in form of URI) of Kubernetes master."),
			},
			"username": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The username to use for HTTP basic authentication when accessing the Kubernetes master endpoint."),
			},
			"password": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The password to use for HTTP basic authentication when accessing the Kubernetes master endpoint."),
			},
			"insecure": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Whether server should be accessed without verifying the TLS certificate."),
			},
			"client_certificate": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("PEM-encoded client certificate for TLS authentication."),
			},
			"client_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("PEM-encoded client certificate key for TLS authentication."),
			},
			"cluster_ca_certificate": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("PEM-encoded root certificates bundle for TLS authentication."),
			},
			"config_path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Path to the kube config file, defaults to ~/.kube/config"),
			},
			"config_context": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Context to choose from the config file."),
			},
			"config_context_auth_info": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Authentication info context of the kube config (name of the kubeconfig user, `--user` flag in `kubectl`)"),
			},
			"config_context_cluster": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Cluster context of the kube config (name of the kubeconfig cluster, `--cluster` flag in `kubectl`)"),
			},
			"token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Token to authentifcate a service account."),
			},
			"exec": {
				Constraint: schema.Object{
					Attributes: schema.ObjectAttributes{
						"api_version": {
							Constraint:  schema.LiteralType{Type: cty.String},
							IsRequired:  true,
							Description: lang.Markdown("API version to use when decoding the ExecCredentials resource, e.g. `client.authentication.k8s.io/v1beta1`."),
						},
						"command": {
							Constraint:  schema.LiteralType{Type: cty.String},
							IsRequired:  true,
							Description: lang.Markdown("Command to execute"),
						},
						"env": {
							Constraint: schema.Map{
								Elem: schema.LiteralType{Type: cty.String},
							},
							IsOptional:  true,
							Description: lang.Markdown("List of arguments to pass when executing the plugin."),
						},
						"args": {
							Constraint: schema.List{
								Elem: schema.LiteralType{Type: cty.String},
							},
							IsOptional:  true,
							Description: lang.Markdown("Map of environment variables to set when executing the plugin."),
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
