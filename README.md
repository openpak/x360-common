# x360-common

The Go library the OpenPak Xbox 360 services share. One repository per service
([`x360-kdc`](../x360-kdc), [`x360-gateway`](../x360-gateway), [`x360-xuacs`](../x360-xuacs),
[`x360-liveinfo`](../x360-liveinfo)); anything that is one service's business stays in that
service.

| package       | what                                                                                       |
|---------------|--------------------------------------------------------------------------------------------|
| `kdc`         | the Kerberos side: AS/TGS/user AS messages, tickets, the account/profile store, console CA  |
| `accountcore` | client of the OpenPak account core (gRPC): console/user links, permanent handles, presence |
| `config`      | environment configuration of every x360 service (`X360_*`)                                  |
| `der`         | the small DER reader/writer the Kerberos messages use                                       |
| `rc4hmac`     | RC4-HMAC (etype 23) as the firmware uses it                                                 |
| `wire`        | the header names the gateway passes to the service backends                                 |

## How services consume it

Like `nx-common`: each service has `replace github.com/openpak/x360-common => ../x360-common`
and commits `vendor/` (`go mod vendor`), so the image build needs nothing but the service's own
repository. After changing this module run `go mod vendor` in each service and commit the result
there.

## Tests

`TMPDIR=<scratch> X360_KDC_TESTKEYS=<dir with the dev keys> go test ./...`

Licensed AGPL-3.0, see LICENSE.
