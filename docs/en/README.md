# ss-ui — S-UI panel with the Sudoku sing-box core

**English** | [فارسی](../fa/README.md)

ss-ui is a fork of [s-ui](https://github.com/alireza0/s-ui) — a web panel and user
manager for the sing-box proxy core — that is compiled against a fork of sing-box
carrying the **Sudoku** protocol, instead of the official core.

- Panel: this repository.
- Core: [`SUDOKU-ASCII/sudoku-sing-box`](https://github.com/SUDOKU-ASCII/sudoku-sing-box) (sing-box 1.14.2 + Sudoku).
- Protocol: [`SUDOKU-ASCII/sudoku`](https://github.com/SUDOKU-ASCII/sudoku).

The sing-box core is **not** a separate binary here: it is linked into the `ss-ui`
executable as a library. Swapping the core is a `replace` directive in `go.mod`, so the
panel and the core always ship together as one binary.

---

## Table of contents

1. [What is the Sudoku protocol](#1-what-is-the-sudoku-protocol)
2. [What is ss-ui](#2-what-is-ss-ui)
3. [Installation](#3-installation)
4. [First run and creating a Sudoku inbound](#4-first-run-and-creating-a-sudoku-inbound)
5. [Connecting a client](#5-connecting-a-client)
6. [Updating and uninstalling](#6-updating-and-uninstalling)
7. [Automatic updates](#7-automatic-updates)
8. [FAQ and troubleshooting](#8-faq-and-troubleshooting)

---

## 1. What is the Sudoku protocol

Sudoku is a traffic-obfuscation protocol that abandons randomness. Instead of adding
random noise, it maps a data stream into **4×4 Sudoku puzzles**. Every byte (256
possibilities) maps onto one of the 288 non-isomorphic 4×4 Sudoku variants, and each
puzzle can be encoded by several valid clue sets — so the same data yields many
different wire encodings, and the traffic can be shaped to look like plaintext or like
low-entropy data.

Highlights:

- **Arbitrary padding** at any position, without negotiating a padding length first.
- **Plaintext / low-entropy modes:** `prefer_ascii` (≈98% printable ASCII) or
  `prefer_entropy` (≈65% ASCII, low Hamming weight) — below the range the GFW Report
  associates with blocking.
- **User-defined fingerprint:** a `custom_table` / `custom_tables` layout lets you pick
  what the bytes look like. Mixing layouts between users is recommended.
- **HTTPMask:** the stream can be shaped to look like HTTP(S)/WebSocket traffic.
- **Defensive fallback:** when a handshake looks wrong (scanner, wrong key, malformed
  data), the server transparently forwards the connection to a decoy address
  (`fallback_address`) — e.g. an Nginx page, or **another Sudoku server** for a two-layer
  chained setup.
- **Encryption:** optional AEAD underneath the obfuscation — `chacha20-poly1305`
  (default) or `aes-128-gcm`, with per-connection keys and anti-replay.
- **UDP over TCP (UoT):** UDP is carried inside TCP; there is no raw UDP listener.

Caveats:

- TCP-native; UDP goes through UoT.
- The bare core only offers SOCKS5/HTTP unless you use the official GUI clients.
- Official client support is limited (Mihomo, Sudodroid, sudoku4x4, Anywhere) — see
  [section 5](#5-connecting-a-client).

### Keys: shared vs split

Sudoku authenticates with a key.

- **Shared-key mode:** client and server use the **same string** as `key`. Simplest.
- **Split-key mode:** the server uses a **public key**, the client a matching **private
  key**. Generate a pair with ss-ui (see below) and give only the private key to clients.

---

## 2. What is ss-ui

ss-ui is the s-ui panel with the Sudoku core baked in. Everything s-ui offers is here —
multi-client management, subscriptions, traffic stats, routing, TLS — plus:

- The Sudoku protocol in the inbound and outbound editors.
- A **key-pair generator** for Sudoku.
- Subscription output for Sudoku clients.
- An [automatic upstream-sync workflow](#7-automatic-updates).

Default endpoints after installation:

| What | Value |
|------|-------|
| Panel port | 2095 |
| Panel path | `/app/` |
| Subscription port | 2096 |
| Subscription path | `/sub/` |
| Default credentials | `admin` / `admin` |

> Change the default credentials immediately after the first login.

---

## 3. Installation

### Linux / macOS

```sh
bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

If your fork lives under a different owner/repository, override it:

```sh
SS_UI_REPO=owner/ss-ui bash <(curl -Ls https://raw.githubusercontent.com/owner/ss-ui/main/install.sh)
```

The installer defaults to English; set `SUI_LANG=fa` for Persian:

```sh
SUI_LANG=fa bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

### Alpine Linux (OpenRC)

Alpine has no `bash` by default:

```sh
apk add bash
bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

Manage the service with OpenRC:

```sh
rc-service ss-ui start|stop|restart
rc-update add ss-ui default
```

### Docker

Using Compose:

```sh
mkdir ss-ui && cd ss-ui
wget -q https://raw.githubusercontent.com/crowrelay/ss-ui/main/docker-compose.yml
docker compose up -d
```

Or `docker run`:

```sh
mkdir ss-ui && cd ss-ui
docker run -itd \
    -p 2095:2095 -p 2096:2096 -p 443:443 -p 80:80 \
    -v $PWD/db/:/app/db/ \
    -v $PWD/cert/:/root/cert/ \
    --name ss-ui --restart=unless-stopped \
    ghcr.io/crowrelay/ss-ui:latest
```

### Windows

1. Download the latest Windows package (`ss-ui-windows-amd64.zip`) from the releases page.
2. Extract it.
3. Run `install-windows.bat` **as Administrator** and follow the wizard.
4. Open the panel at `http://localhost:2095/app`.

### Build from source

Requirements: Go (see `go.mod` for the version), Node.js 26, and a C toolchain for CGO.

```sh
git clone https://github.com/crowrelay/ss-ui
cd ss-ui
git submodule update --init --recursive

# one-shot: builds the frontend and the backend into ./ss-ui
./build.sh
```

Or step by step:

```sh
(cd frontend && npm ci && npm run build)
rm -fr web/html/* && cp -R frontend/dist/* web/html/
. ./build-tags.sh
go build -ldflags "$(ldflags_for dev)" -tags "$(tags_for dev)" -o ss-ui main.go
./ss-ui
```

> After changing the core pin in `go.mod`, run `go mod tidy` once so `go.sum` picks up
> the fork before building.

---

## 4. First run and creating a Sudoku inbound

Start the service (`systemctl enable ss-ui --now`, `rc-service ss-ui start`, or run
`./ss-ui`) and log in at `http://<host>:2095/app` with `admin` / `admin`.

### Generate a key pair

Either from the command line:

```sh
./ss-ui generate sudoku-keypair
# PrivateKey: ...
# PublicKey:  ...
```

or from the panel: **Inbounds → Add → type `Sudoku` → the refresh icon next to “Key”.**

- **Shared-key mode:** paste the same string into the inbound's **Key** and the client's
  **Key**.
- **Split-key mode:** put the **PublicKey** in the inbound and the **PrivateKey** in the
  client (`out_json`), so a leak of the server config does not leak the client key.

### Inbound fields

| Field | Meaning |
|-------|---------|
| `key` | Shared secret, or the public key in split-key mode. |
| `aead` | `chacha20-poly1305` (default), `aes-128-gcm`, or `none` (testing only). |
| `ascii` | `prefer_ascii` or `prefer_entropy`. |
| `padding_min` / `padding_max` | Padding bounds, 0–100. |
| `custom_table` | A layout string, e.g. `xpxvvpvv` (2 `x`, 2 `p`, 4 `v`). |
| `custom_tables` | Several layouts; mixing across users is recommended. |
| `enable_pure_downlink` | `true` (default) = pure Sudoku downlink; `false` = bandwidth-optimised downlink (needs `aead != none`). |
| `multiplex` | `off` or `auto`. |
| `httpmask` | Shape the stream as HTTP(S)/WebSocket: `mode` (`auto`/`legacy`/`ws`), `path_root`, and on clients `host`, `tls`. |
| `handshake_timeout` | Handshake timeout in seconds. |
| `suspicious_action` | `fallback` (default) or `silent` for probing connections. |
| `fallback_address` | Where a suspicious connection is forwarded (decoy page, or an inner Sudoku server). |

A minimal inbound:

```json
{
  "type": "sudoku",
  "tag": "sudoku-in",
  "listen": "0.0.0.0",
  "listen_port": 8080,
  "key": "YOUR_PUBLIC_OR_SHARED_KEY",
  "aead": "chacha20-poly1305",
  "ascii": "prefer_entropy",
  "padding_min": 5,
  "padding_max": 15,
  "enable_pure_downlink": true,
  "handshake_timeout": 5,
  "fallback_address": "127.0.0.1:80",
  "suspicious_action": "fallback"
}
```

### Two-layer (chained) setup

Use an outer Sudoku server as the public entry point (disguise, anti-probing). Give it a
**different key** from the inner server and set `suspicious_action: "fallback"` with
`fallback_address` pointing at the inner Sudoku server. Normal clients fail the outer
handshake and are forwarded to the inner server, which performs the real handshake. The
client dials the outer address but uses the **inner server's key**.

---

## 5. Connecting a client

Sudoku is supported by:

| Platform | Client |
|----------|--------|
| Android | [CMFA](https://github.com/MetaCubeX/ClashMetaForAndroid) (Mihomo), [Sudodroid](https://github.com/SUDOKU-ASCII/sudoku-android) |
| iOS | [Anywhere Proxy](https://apps.apple.com/us/app/anywhere-proxy/id6758235178) |
| Desktop | [sudoku4x4](https://github.com/SUDOKU-ASCII/sudoku-desktop) |
| Router | [OpenClash](https://github.com/vernesong/OpenClash) |
| Server/tools | [easy-install](https://github.com/SUDOKU-ASCII/easy-install) |

In ss-ui, add a **client**, assign it the Sudoku inbound, and share the subscription. The
Sudoku outbound the client receives is built from the inbound's **Client `out_json`**
block — open the inbound, switch to the **Client** tab, and set the client `key` there
for split-key mode.

Client outbound shape:

```json
{
  "type": "sudoku",
  "tag": "sudoku-out",
  "server": "example.com",
  "server_port": 443,
  "key": "YOUR_PRIVATE_OR_SHARED_KEY",
  "aead": "chacha20-poly1305",
  "ascii": "prefer_entropy",
  "padding_min": 5,
  "padding_max": 15,
  "enable_pure_downlink": true,
  "httpmask": { "disable": false, "mode": "auto", "tls": true, "host": "example.com", "path_root": "aabbcc" }
}
```

---

## 6. Updating and uninstalling

Update to the latest release:

```sh
bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

Install a specific version:

```sh
VERSION=v1.6.4-sudoku.1 bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh) $VERSION
```

Uninstall (systemd):

```sh
systemctl disable ss-ui --now
rm -f /etc/systemd/system/ss-ui.service
systemctl daemon-reload
rm -fr /usr/local/ss-ui /usr/bin/ss-ui
```

Alpine (OpenRC):

```sh
rc-service ss-ui stop
rc-update del ss-ui default
rm -f /etc/init.d/ss-ui
rm -fr /usr/local/ss-ui /usr/bin/ss-ui
```

---

## 7. Automatic updates

`.github/workflows/sync-upstream.yml` runs daily (and on demand):

1. Resolves the latest `sudoku-sing-box` tag and rewrites the `replace` line in `go.mod`.
2. Fetches `alireza0/s-ui` and merges it best-effort (`-X ours`, so this fork's customised
   files and the frontend submodule pointer win).
3. Runs `go mod tidy`, a build check and the core/util tests.
4. Commits, pushes to `main`, and pushes a version tag — which triggers
   `release.yml` to build every platform and publish the release.

If the upstream merge conflicts, it is aborted (only the core pin is released) and an
issue labelled `sync-conflict` is opened for a human. `.gitattributes` marks `frontend`
as `merge=ours`; the frontend is refreshed from its own fork.

---

## 8. FAQ and troubleshooting

**The core will not start / “missing key”.**
The Sudoku inbound needs a non-empty `key`. Generate one and paste it in.

**The client connects but nothing loads.**
Check that client and server keys match (shared-key mode), or that the client's private
key pairs with the server's public key (split-key mode), and that `aead`, `ascii` and
`custom_table` are consistent.

**`prefer_ascii` vs `prefer_entropy`?**
`prefer_ascii` looks like plaintext (≈98% ASCII); `prefer_entropy` is lower-entropy and
less obviously printable. If unsure, mix layouts across users.

**Fallback never fires.**
Confirm `suspicious_action: "fallback"` and a reachable `fallback_address`.

**Where do the Sudoku docs live?**
Upstream sing-box docs do not document Sudoku; the Doc link in the panel points at the
[sudoku-sing-box](https://github.com/SUDOKU-ASCII/sudoku-sing-box) repository.
