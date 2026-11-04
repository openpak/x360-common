# Deploy: x360-common

## How services consume it

Like `nx-common`: each service has `replace github.com/openpak/x360-common => ../x360-common`
and commits `vendor/` (`go mod vendor`), so the image build needs nothing but the service's own
repository. After changing this module run `go mod vendor` in each service and commit the result
there.

## Tests

`TMPDIR=<scratch> X360_KDC_TESTKEYS=<dir with the dev keys> go test ./...`

## Account erasure (card #120)

`erase` serves `DELETE /internal/accounts/{account_id}/data` (bearer `ACCOUNT_ERASE_KEY`, listener `ERASE_ADDR`, default `:9099`, container network only). It is not served while the key is unset. `erase.Run` logs the single `account erased` line. Services pass an `erase.Func` that deletes the account's rows; `xuacsapi.Client.AccountXUIDs` lists the XUIDs of an account for the services whose data is keyed by XUID.
