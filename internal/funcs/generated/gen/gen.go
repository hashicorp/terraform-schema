// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/go-version"
	hcinstall "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

var (
	terraformVersion = version.Must(version.NewVersion("1.5.0-beta1"))
)

const (
	functionSignatureHash = "3edcd73cb8643903dde229b04dfc36d10d8dd6679e7804ea37001be38650d950"
)

func main() {
	ctx := context.Background()

	functions, err := signaturesFromTerraform(ctx)
	if err != nil {
		log.Fatal(err)
	}
	newSignatureHash, err := signatureHash(functions)
	if err != nil {
		log.Fatal(err)
	}
	if newSignatureHash == functionSignatureHash {
		log.Println("function signatures haven't changed, nothing to do here")
		return
	}
	log.Printf("generating new signatures for %q\n", newSignatureHash)
	functionsFile := fmt.Sprintf("%s.go", terraformVersion.Core().String())
	err = writeFunctions(functionsFile, functions)
	if err != nil {
		log.Fatal(err)
	}

	versions, err := generatedVersions(".")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found versions %s, regenerating includes", versions)
	err = writeFunctionVersions("functions.go", versions)
	if err != nil {
		log.Fatal(err)
	}
}

// signaturesFromTerraform gets the function signatures for the specified
// Terraform version.
func signaturesFromTerraform(ctx context.Context) (*tfjson.MetadataFunctions, error) {
	// find or install Terraform
	log.Println("ensuring terraform is installed")
	installDir, err := os.MkdirTemp("", "hcinstall")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(installDir)
	i := hcinstall.NewInstaller()
	execPath, err := i.Ensure(ctx, []src.Source{
		&releases.ExactVersion{
			Product:    product.Terraform,
			InstallDir: installDir,
			Version:    terraformVersion,
		},
	})
	if err != nil {
		return nil, err
	}
	defer i.Remove(ctx)

	// log version
	tf, err := tfexec.NewTerraform(installDir, execPath)
	if err != nil {
		return nil, err
	}
	coreVersion, _, err := tf.Version(ctx, true)
	if err != nil {
		return nil, err
	}
	log.Printf("using Terraform %s (%s)", coreVersion, execPath)

	return tf.MetadataFunctions(ctx)
}

// signatureHash calculates the SHA256 checksum for the given function
// signatures.
func signatureHash(functions *tfjson.MetadataFunctions) (string, error) {
	var rawJson bytes.Buffer
	err := json.NewEncoder(&rawJson).Encode(functions)
	if err != nil {
		return "", fmt.Errorf("failed to encode functions: %w", err)
	}

	hash := sha256.Sum256(rawJson.Bytes())
	return fmt.Sprintf("%x", hash), nil
}

// writeFunctions generates the Go code for the function signatures,
// formats it and writes it to a new file.
func writeFunctions(filename string, functions *tfjson.MetadataFunctions) error {
	outputTpl := `// Code generated by "gen"; DO NOT EDIT.
package funcs

import (
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)
func {{ .FunctionName }}() map[string]schema.FunctionSignature {
	return map[string]schema.FunctionSignature{
{{- range $funcName, $func := .Signatures }}
		{{- $varpar := $func.VariadicParameter }}
		{{- $params := $func.Parameters }}
		"{{ $funcName }}": {
			{{- with $params }}
			Params: []function.Parameter{
				{{- range $key, $value := $params }}
				{{- $desc := $value.Description }}
				{
					Name:        "{{ $value.Name }}",
					{{- with $desc }}
					Description: "{{ escapeDescription . }}",{{- end }}
					Type:        {{ $value.Type.GoString }},
				},
				{{- end }}
			},{{- end }}
			{{- with $varpar }}
			VarParam:    &function.Parameter{
				Name:             "{{ .Name }}",
				Description:      "{{ .Description }}",
				Type:             {{ .Type.GoString }},
			},{{- end }}
			ReturnType: {{ $func.ReturnType.GoString }},
			Description: "{{ escapeDescription $func.Description }}",
		},	
{{- end }}
	}
}
`

	tpl, err := template.New("output").Funcs(template.FuncMap{
		"escapeDescription": escapeDescription,
	}).Parse(outputTpl)
	if err != nil {
		return err
	}

	type data struct {
		FunctionName string
		Signatures   map[string]*tfjson.FunctionSignature
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data{
		FunctionName: fmt.Sprintf("v%s_Functions", escapeVersion(terraformVersion.Core().String())),
		Signatures:   functions.Signatures,
	})
	if err != nil {
		return err
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	_, err = f.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func escapeDescription(description string) string {
	s := strings.ReplaceAll(description, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func escapeVersion(version string) string {
	return strings.ReplaceAll(version, ".", "_")
}

// generatedVersions returns a sorted slice of all available function
// signature files within a path. It will skip all directories and
// filenames that don't parse as a version string.
func generatedVersions(path string) ([]*version.Version, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	versions := []*version.Version{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		v, err := version.NewVersion(filename)
		if err != nil {
			// Skip any filenames that don't parse as version string
			continue
		}
		versions = append(versions, v)
	}

	sort.Sort(sort.Reverse(version.Collection(versions)))

	return versions, nil
}

// writeFunctionVersions generates Go code for selecting the correct generated
// function signature for a specific version.
func writeFunctionVersions(filename string, versions []*version.Version) error {
	outputTpl := `// Code generated by "gen"; DO NOT EDIT.
package funcs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

var (
{{- range $version := .Versions }}
	v{{ escapeVersion $version.String }} = version.Must(version.NewVersion("{{ $version.String }}"))
{{- end }}
)

func Functions(v *version.Version) map[string]schema.FunctionSignature {
{{- range $version := .Versions }}
	if v.GreaterThanOrEqual(v{{ escapeVersion $version.String }}) {
		return v{{ escapeVersion $version.String }}_Functions()
	}
{{- end }}

	return v1_4_0_Functions()
}	
`

	tpl, err := template.New("output").Funcs(template.FuncMap{
		"escapeDescription": escapeDescription,
		"escapeVersion":     escapeVersion,
	}).Parse(outputTpl)
	if err != nil {
		return err
	}

	type data struct {
		Versions []*version.Version
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data{
		Versions: versions,
	})
	if err != nil {
		return err
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	_, err = f.Write(p)
	if err != nil {
		return err
	}

	return nil
}
