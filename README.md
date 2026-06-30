# pqc-proxy (Hybrid Post-Quantum TCP Tunnel)

A lightweight infrastructure proxy server (Client/Server architecture) engineered to secure legacy TCP application traffic against interception (sniffing) and retrospective decryption utilizing quantum cryptanalysis methodologies (e.g., Shor's algorithm).

## Why Post-Quantum?

Traditional cryptography (RSA, ECDH) relies on the hardness of integer factorization and discrete logarithms. Large-scale quantum computers will render these obsolete. **pqc-proxy** provides an immediate security layer for legacy applications by implementing a hybrid approach, combining classical security with lattice-based cryptography today, ensuring "harvest now, decrypt later" protection.

## Architecture

The project implements a **Crypto-Agility** paradigm through a hybrid key exchange mechanism.

* **Classical Layer:** Diffie-Hellman over Curve25519 (`X25519`).
* **Post-Quantum Layer:** `ML-KEM-768` (NIST FIPS 203 standard).
* **KDF:** `HKDF-SHA256` (concatenated shared secrets).
* **Transport Layer (AEAD):** `ChaCha20-Poly1305` with per-frame monotonicity.
* **Performance:** Zero-Allocation Pipeline using `sync.Pool`.

## Project Structure

```text
├── .github/workflows/   # CI/CD pipelines (GitHub Actions)
├── cmd/
│   └── pqc-proxy/       # Application entry point
├── deployments/         # Docker Compose and Prometheus manifests
├── internal/
│   ├── config/          # Configuration parsing
│   ├── crypto/          # Cryptographic core engine
│   ├── network/         # Client, Server, Chaos injection & Pipe tests
│   └── ...
├── scripts/             # Automation scripts
└── web/                 # Monitoring dashboard static files

```

## Compilation

Compile from the root directory:

**Windows:** `go build -o pqc-proxy.exe ./cmd/pqc-proxy`

**Linux/macOS:** `go build -o pqc-proxy ./cmd/pqc-proxy`

## Verification & Testing

The project features a comprehensive test suite covering cryptographic primitives and network pipeline stability.

```bash
# Validate cryptographic integrity and network pipe logic
go test -v ./internal/...

```

## Quick Start Topology

```text
[Client App] -> (Local:3000) -> [PQC Client] -> (Tunnel:9090) -> [PQC Server] -> (Target:8000) -> [Backend App]

```

1. **Start Backend:** `python -m http.server 8000`
2. **Start Server:** `./pqc-proxy -mode server -listen :9090 -target 127.0.0.1:8000`
3. **Start Client:** `./pqc-proxy -mode client -listen :3000 -target 127.0.0.1:9090`
4. **Access:** `curl http://127.0.0.1:3000`

## Status & Roadmap

* [x] Hybrid Key Exchange (X25519 + ML-KEM-768)
* [x] AEAD Transport Encryption (ChaCha20-Poly1305)
* [x] CI/CD Pipeline (GitHub Actions)
* [x] Chaos Testing Framework
* [ ] UDP Support
* [ ] Certificate-based Authentication

## License

Distributed under the MIT License.
