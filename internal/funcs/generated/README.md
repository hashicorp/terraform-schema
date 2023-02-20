# Generated Terraform function signatures

This package can generate function signature files for Terraform >= 1.4 automatically.

It is intended to run whenever HashiCorp releases a new Terraform version. If the Terraform version contains updated function signatures, it generates a new file for that version. When no changes are detected, one should commit the version bump in `gen/gen.go`.

It may only work with full releases and will likely fail for pre-releases or versions containing metadata.

## Running

Update the `terraformVersion` in `gen/gen.go` with the version of the new Terraform release.

Run `go generate ./internal/funcs/generated`. This command will:

1. Install the specified Terraform version.
1. Run `terraform metadata functions -json` to obtain the function signatures.
1. Compare a hash of the JSON string to `functionSignatureHash` in `gen/gen.go`.
   1. If there are no changes, we will stop here.
1. Create a new Go file with the function signatures of that Terraform version.
1. Regenerate the function signature selection in `functions.go`.

If everything looks solid, update the `functionSignatureHash` with the one from the output, and commit all changes.
