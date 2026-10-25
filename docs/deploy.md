# Deploy: x360-common

## How services consume it

Like `nx-common`: each service has `replace github.com/openpak/x360-common => ../x360-common`
and commits `vendor/` (`go mod vendor`), so the image build needs nothing but the service's own
repository. After changing this module run `go mod vendor` in each service and commit the result
there.

## Tests

`TMPDIR=<scratch> X360_KDC_TESTKEYS=<dir with the dev keys> go test ./...`
