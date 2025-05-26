# README.md

## Overview

The SIM-Auth-Token-Broker is a Go microservice that:

1. Routes by SIM MCC+MNC prefix (via `prefix_map.yaml`)  
2. Proxies OAuth2/OIDC `/token` calls to upstream Telco mocks  
3. Validates Telco JWTs (JWKS cache) and issues its own JWTs

## Prerequisites

- Go 1.24.3+  
- Docker
- `make`, `bash`

## Setup

```bash
cp .env.example .env
# Edit .env for ENV, PORT, SECRET_… values
```

## Running locally

```bash
# to run all services concurrently 
make dev-all
```

## Usage

```bash
 curl -X POST http://localhost:8080/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=yourAuthorizationCode" \
  -d "phone=+972541234567" \
  -d "redirect_uri=https://your.client/callback" \
  -d "code_verifier=yourCodeVerifier"
```
