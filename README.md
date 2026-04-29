# Oshudh Somadhan

Oshudh Somadhan is a Bangladesh-focused medicine information and pharmacy backend project built with Go.

## Phase 0 Features

- Go project setup
- Docker Compose setup
- PostgreSQL database
- Existing medicine database import
- Environment configuration
- Health endpoint
- Development DB ping endpoint
- Basic tests

## Tech Stack

- Go
- Gin
- PostgreSQL
- pgx
- Docker Compose
- Redis
- Testify

## Run Locally

```bash
make docker-up
make db-import
make run