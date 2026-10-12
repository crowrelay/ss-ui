# ss-ui

**An advanced web panel for the [Sudoku](https://github.com/SUDOKU-ASCII/sudoku) sing-box core**

ss-ui is a fork of [s-ui](https://github.com/alireza0/s-ui) that is built against the
[`sudoku-sing-box`](https://github.com/SUDOKU-ASCII/sudoku-sing-box) core (sing-box 1.14.2
+ the Sudoku protocol) instead of the official sing-box core. The core is linked into the
`ss-ui` binary as a library, so one executable is the whole product.

📖 **Documentation:** [English](docs/en/README.md) · [فارسی](docs/fa/README.md)

> **Disclaimer:** This project is for personal learning and communication. Do not use it
> for illegal purposes.

## Quick overview

| Features | Enabled |
|----------|:-------:|
| **Sudoku protocol (inbound + outbound)** | ✅ |
| Multi-Protocol | ✅ |
| Multi-Language | ✅ |
| Multi-Client / Inbound | ✅ |
| Advanced Traffic Routing | ✅ |
| Client, traffic & system status | ✅ |
| Subscription (link / json / clash + info) | ✅ |
| Dark / Light theme | ✅ |
| REST API | ✅ |

## Supported platforms

| Platform | Architecture | Status |
|----------|--------------|--------|
| Linux | amd64, arm64, armv7, armv6, armv5, 386, s390x | ✅ Supported |
| Windows | amd64, 386, arm64 | ✅ Supported |
| macOS | amd64, arm64 | 🚧 Experimental |

## Install

Linux / macOS:

```sh
bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

Alpine Linux (OpenRC): `apk add bash` first, then run the same command.

Docker:

```sh
mkdir ss-ui && cd ss-ui
wget -q https://raw.githubusercontent.com/crowrelay/ss-ui/main/docker-compose.yml
docker compose up -d
```

Windows: download `ss-ui-windows-<arch>.zip` from the releases page and run
`install-windows.bat` as Administrator.

Full instructions (per OS, from source, and the Sudoku field reference) are in the
[documentation](docs/en/README.md).

## Default installation information

- Panel port: **2095**  · Panel path: `/app/`
- Subscription port: **2096**  · Subscription path: `/sub/`
- User / password: **admin** / **admin** (change it immediately)

## Using Sudoku

1. Create an inbound of type **Sudoku**.
2. Generate a key pair — the panel's refresh button or `./ss-ui generate sudoku-keypair`.
3. Use the **same** key on both sides (shared-key mode), or the **public** key inbound /
   **private** key client (split-key mode).
4. Assign the inbound to a client and share the subscription.

See the [English](docs/en/README.md#4-first-run-and-creating-a-sudoku-inbound) or
[Persian](docs/fa/README.md) docs for every field and for the two-layer fallback setup.

## Automatic updates

`.github/workflows/sync-upstream.yml` runs daily: it bumps the `sudoku-sing-box` pin,
merges the latest `alireza0/s-ui` (best-effort, `-X ours`), builds, and — if the build
passes — tags a release that `release.yml` publishes for every platform.

## Repository layout

- The Go module path is still `github.com/alireza0/s-ui` (upstream), so automatic merges
  from s-ui stay clean. User-facing names, the binary, service, install paths and release
  assets are **ss-ui**.
- The `frontend/` directory is a git submodule of a fork of `s-ui-frontend` carrying the
  Sudoku UI. It is marked `merge=ours` in `.gitattributes`.

## Languages

English · Farsi · Vietnamese · Chinese (Simplified) · Chinese (Traditional) · Russian

## License

GPLv3, as upstream. See [LICENSE](LICENSE).
