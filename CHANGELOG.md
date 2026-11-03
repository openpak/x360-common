# Changelog

- 2026-11-03: v0.6.0 kdc: the account store can live in Postgres (`kdc.OpenPGStore`, schema `x360_kdc`, document `accounts`):
  one-time import of `accounts.json` (file left in place), `Reload` for read-only users (the gateway), `Close`. The file store
  stays for tests. pgstate.Open uses a schema that already exists without issuing CREATE SCHEMA, so a role without CREATE
  on the database can use a schema made for it.

- 2026-11-03: v0.5.1 kdc.Server: every datagram line (answered, rejected, recorded) carries `ms`, the handling time.

- 2026-11-03: v0.5.0 httplog: one stderr line per request, `http <method> <path> -> <status> <duration>ms`, and
  `httplog.Command` for UDP/TCP commands. Never the query, headers, cookies, tokens or bodies; `/healthz` skipped;
  the status writer keeps Flusher, Hijacker and ResponseController working.

- 2026-11-03: v0.4.0 pgstate: Postgres document store (one schema per service, pool of 2, refuses an empty URL) so
  service state lives in the database, not in files.

- 2026-11-02: v0.3.0 accountcore: ListIncomingRequests, ListOutgoingRequests, Block and Unblock clients.

- 2026-11-01: v0.2.0 xuacsapi: client for x360-xuacs's internal identity API (identity, online XUID, XUID to
  account, console profiles, WLID CA), used by the services split out of xuacs.

- 2026-10-20: read account-core birthdays for native age policy, requiring an active matching account.

- 2026-10-20: Social.ListFriends, ListBlocks and discovery-enforcement clients.

- 2026-10-19: Media.GetAvatar client for account avatars.

- 2026-10-19: compare decimal REST path XUIDs with hexadecimal token claims, retaining hexadecimal
  diagnostic paths and checking the full uint64 range.

- 2026-10-18: split out of x360-kdc v0.8.14 (history kept): `kdc`, `accountcore`, `config`, `der`,
  `rc4hmac`, plus the new `wire` package (the account/users header names, formerly constants of the
  gateway package). Import path `github.com/openpak/x360-common/<package>`.
- 2026-10-18: `kdc` D7 record: the LE u32 at 0x98 now carries the TGS request's title id (it was zero). The console's
  logon engine stores it as the title of the services in that ticket; its token service (XSTS, REST data recovery)
  stays pending until the service-12 title equals the running title.
- 2026-10-18: new `xstoken` package: the token x360-xsts issues and the REST services verify (Claims, Signer,
  FromAuthorization), moved out of x360-xsts.
