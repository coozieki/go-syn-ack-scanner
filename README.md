# About

This is a training project whose purpose is to scan all ports of a given IP address and return the result (a list of open ports) to the console

# Usage

Run from the root of the project
```bash
go run cmd/app/main.go ip max_threads
```
### Params:
- `ip` (*required*) should be replaced with ip of the server you want to scan

  **Example**: 127.0.0.1
- `max_threads` (*optional, default: 8*) determines how many go routines would be used simultaniously for scanning

  **Example**: 10