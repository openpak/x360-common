# Changelog

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
