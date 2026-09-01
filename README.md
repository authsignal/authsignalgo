<img width="1070" alt="Authsignal" src="https://raw.githubusercontent.com/authsignal/authsignalgo/main/.github/images/authsignal.png">

# Authsignal Go SDK

## Installation

```
go get github.com/authsignal/authsignalgo
```

## Retry policy

Requests use a 3-second connect timeout, 10-second total timeout, and retry twice by default with exponential backoff and jitter. Transient network failures, `429`, and `5xx` responses are retried for `GET`, `HEAD`, and `OPTIONS`; writes are retried only when they carry an idempotency key. Set `client.Retries = 0` to disable retries or replace `client.Client` to customize HTTP timeouts.

## Documentation

Check out our [official documentation](https://docs.authsignal.com/api-reference/server-api/overview) to get up and running quickly.
