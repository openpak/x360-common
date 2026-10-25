# Xbox 360 Common Library

OpenPak's Xbox 360 Common Library.

| | |
| --- | --- |
| Network | Xbox 360 |
| Language | Go module `github.com/openpak/x360-common` |
| Used by | [`x360-kdc`](../x360-kdc), [`x360-gateway`](../x360-gateway), [`x360-xuacs`](../x360-xuacs), [`x360-liveinfo`](../x360-liveinfo) |
| Licence | AGPL-3.0-only |

## Packages

| Package | What |
| --- | --- |
| `kdc` | the Kerberos side: AS/TGS/user AS messages, tickets, the account/profile store, console CA |
| `accountcore` | client of the OpenPak account core (gRPC): console/user links, permanent handles, presence |
| `config` | environment configuration of every x360 service (`X360_*`) |
| `der` | the small DER reader/writer the Kerberos messages use |
| `rc4hmac` | RC4-HMAC (etype 23) as the firmware uses it |
| `wire` | the header names the gateway passes to the service backends |
| `xstoken` | the token x360-xsts issues and the REST services verify (Claims, Signer, FromAuthorization) |

## Use

```go
import "github.com/openpak/x360-common/config"
```

Each service uses `replace github.com/openpak/x360-common => ../x360-common` and commits `vendor/`; details in [deploy](docs/deploy.md).

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)
- [Permanent account handles](docs/handles.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
