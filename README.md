# Bank System

A small **distributed systems** demo: a **Java** UDP server that hosts in-memory bank accounts, and a **Go** command-line client. The design highlights **unreliable UDP**, configurable **invocation semantics** (at-least-once vs at-most-once), **retries with exponential backoff** on the client, **simulated packet loss**, and **push callbacks** for registered clients when the server processes updates.

## Repository layout

| Path | Role |
|------|------|
| `server/` | Java 21. Entry point: `system.Server` |
| `client/` | Go module (`go.mod`). UDP client with menus and RPC helpers |

## Features

**Banking operations (via client menu)**

- Create account (username, password, currency, initial deposit)
- Close account
- Deposit and withdraw (withdraw is sent as a **negative** deposit on the wire)
- View all balances for an account (multi-currency)
- Transfer between accounts

**Cross-cutting behavior**

- **Multi-currency** balances per account: SGD, EUR, USD, JPY (`CurrencyType` on the server)
- **Authentication** on sensitive operations: username, account ID, and password must match the stored account
- **Monitor / callbacks**: the client can register for a duration (seconds). While active, the server sends UDP **callback** messages to that client when it finishes handling requests (so observers see activity). Callback lines are prefixed for the Go listener (`8:CALLBACK`…)


## Invocation semantics

These are meant to pair with a course-style discussion of RPC over unreliable channels.

**Server (chosen at startup)**  
On boot, `Server` asks for:

1. **At-Least-Once** — each datagram is handled once; duplicates may run the operation again.
2. **At-Most-Once** — the first field of the payload is treated as a **request id** (UUID string). Duplicate ids return the **cached** reply without re-executing the handler (`AtMostOnce`).

**Client (chosen before main menu)**  
The Go client selects:

- **Default** — single send, optional simulated loss, fixed reply wait (`TIMEOUT_MS`).
- **At-Least-Once** — retries `defaultInvocation` with **exponential backoff** (same request bytes each time).
- **At-Most-Once** — prepends `len(requestID):requestID` to the request; retries use the **same** id so the server can deduplicate.

For consistent at-most-once behavior, run the server in **At-Most-Once** mode and the client in **At-Most-Once** mode.

**Packet loss simulation**

- **Server** (`Server.java`): after building the reply, it may **drop** the outbound packet randomly (`packetLossRate`, default `0.1`).
- **Client** (`rpc.go`): in Default/At-Least-Once/At-Most-Once paths, `defaultInvocation` may **skip sending** the request to simulate loss (`PACKET_LOSS_PROBABILITY`, default `0.1`).

## Prerequisites

- **Java 21** (for `server/`)
- **Go** toolchain matching `client/go.mod` (module declares `go 1.25.0`; use the same or newer Go version your environment supports)

## Run the server

From the `java` directory:

```bash
javac system/*.java
java system.Server
```

On startup you will be prompted to select **At-Least-Once** (`1`) or **At-Most-Once** (`2`). The server then listens on UDP **port 2222**.


## Run the client

From the `client` directory:

```bash
cd client
go run .
```

You will first pick client **invocation** mode, then the **main menu** for banking operations. Ensure the server is already running on `localhost:2222` (or change `SERVER_IP` / `SERVER_PORT` in `main.go`).

## Seeded demo accounts

`AccountHandler` preloads three accounts (password for all is **`123`**):

| Account ID | Username | Initial currency / balance |
|-----------|----------|----------------------------|
| 1000 | tom | SGD 50 |
| 1001 | dick | USD 100 |
| 1002 | harry | EUR 290 |

New accounts receive the next integer ID from the server’s counter after these seeds.

## Configuration (quick reference)

| Location | Constant | Meaning |
|----------|----------|---------|
| `client/main.go` | `SERVER_IP`, `SERVER_PORT` | Server address |
| `client/main.go` | `TIMEOUT_MS` | Reply wait in default invocation |
| `client/main.go` | `PACKET_LOSS_PROBABILITY` | Client-side simulated send loss |
| `client/rpc.go` | `retries` (via callers) | At-least-/at-most-once retry cap (10 in `sendRequestReceiveReply`) |
| `server/.../Server.java` | `PORT`, `BUFFER_SIZE` | Listen port and receive buffer |
| `server/.../Server.java` | `packetLossRate` | Server-side simulated reply drop rate |

## Protocol sketch

Commands are length-prefixed tokens; the **first** logical field is the command name (e.g. `CREATEACCOUNT`, `CLOSEACCOUNT`, `DEPOSIT`, `VIEW`, `TRANSFER`, `MONITOR`). The server responds with length-prefixed fields as well (see `RequestHandler` return strings for success/error tokens like `CREATEACCOUNTSUCCESS`, `FAIL`, etc.).

At-most-once requests from the client add a **leading** `len(id):id` block before the usual command fields; the server strips the id in `AtMostOnce` before dispatching to `RequestHandler` / `MonitorHandler`.

# Contributors
| Name              | GitHub Username                          |
| ----------------- | ---------------------------------------- |
| Chen Xing Wei     | [dozetype](https://github.com/dozetype)  |
| Phua Wei Jie      | [Alpths](https://github.com/Alpths)      |
