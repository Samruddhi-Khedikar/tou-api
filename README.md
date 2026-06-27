# TOU Pricing API

A RESTful API for managing and querying **Time-of-Use (TOU) electricity pricing** for individual EV charging stations.

---

## Table of contents

- [Overview](#overview)
- [Stack](#stack)
- [Project structure](#project-structure)
- [Getting started](#getting-started)
- [Environment](#environment)
- [Available commands](#available-commands)
- [API reference](#api-reference)
- [Design decisions](#design-decisions)
- [Development tooling](#development-tooling)

---

## Overview

TOU pricing varies electricity cost by time of day. This service lets you:

- Register EV charging stations with an IANA timezone
- Set time-based pricing schedules per charger
- Query the applicable price at any point in time
- Apply pricing in bulk across multiple chargers

---

## Stack

| Layer | Choice | Reason |
|---|---|---|
| Language | Go 1.22 | Native concurrency, fast compile, strong stdlib |
| Router | `chi` | Lightweight, idiomatic, `net/http` compatible |
| Database | SQLite (`modernc.org/sqlite`) | Pure Go, zero setup for reviewer |
| Validation | `go-playground/validator` | Struct-tag validation |

---

## Project structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── db/
│   │   └── db.go             # SQLite connection + migration runner
│   ├── handler/
│   │   ├── routes.go         # Route registration
│   │   ├── charger_handler.go
│   │   ├── pricing_handler.go
│   │   └── helpers.go        # parseAndValidate, writeJSON, writeError
│   ├── model/
│   │   └── model.go          # Domain types, request/response structs
│   ├── repository/
│   │   ├── charger_repo.go
│   │   ├── pricing_repo.go
│   │   └── errors.go         # Sentinel errors (ErrNotFound)
│   └── service/
│       ├── charger_service.go
│       └── pricing_service.go # Business logic, period validation
├── migrations/
│   └── 001_init.sql          # Schema definition
├── .env.example
├── Makefile
└── openapi.yaml              # OpenAPI 3.0 spec
```

---

## Getting started

**Prerequisites:** Go 1.22+

```bash
git clone https://github.com/Samruddhi-Khedikar/tou-api.git
cd tou-api-v2
cp .env.example .env
make run
# → server listening on :8080
```

No Docker, no external database. SQLite file (`tou.db`) is created automatically on first run.

---

## Environment

Copy `.env.example` to `.env` and adjust if needed:

```env
PORT=8080
DB_PATH=tou.db
```

---

## Available commands

```bash
make run        # run the server
make build      # compile binary to bin/tou-api
make test       # run all tests
make tidy       # go mod tidy
make db-shell   # open sqlite3 shell against tou.db
```

---

## API reference

Full spec: [`tou-api.yaml`](./tou-api.yaml)

Preview options:
- Paste into [Swagger Editor](https://editor.swagger.io)
- VS Code: install **OpenAPI (Swagger) Editor** by 42Crunch → right-click the file → Preview Swagger
