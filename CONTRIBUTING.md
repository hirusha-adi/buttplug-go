# Contributing

Thanks for your interest in contributing to [buttplug-go](https://github.com/hirusha-adi/buttplug-go)!

## Code of Conduct

Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before participating.

## Getting Started

1. Fork [hirusha-adi/buttplug-go](https://github.com/hirusha-adi/buttplug-go) on GitHub.
2. Clone your fork and create a feature branch.
3. Make your changes and add tests where appropriate.
4. Run the test suite:

   ```bash
   go test ./...
   ```

5. Open a pull request against `main` on [hirusha-adi/buttplug-go](https://github.com/hirusha-adi/buttplug-go).

## Project Goals

This library aims to stay a faithful Go translation of the official Python client:

- [buttplug-py](https://github.com/buttplugio/buttplug-py) — upstream reference implementation
- [Buttplug protocol docs](https://docs.buttplug.io/docs/spec) — protocol specification

When adding features, prefer matching the Python API and behavior unless there is a strong Go-idiomatic reason not to (for example, using `context.Context` instead of `async`/`await`).

## Bug Reports & Feature Requests

Open an issue on GitHub:

**https://github.com/hirusha-adi/buttplug-go/issues**

Please include:

- Go version (`go version`)
- Operating system
- Buttplug server in use (e.g. [Intiface Central](https://intiface.com/central/))
- Steps to reproduce
- Expected vs actual behavior

## Pull Requests

- Keep changes focused and reasonably sized.
- Match existing code style and naming.
- Update documentation and examples when behavior changes.
- Ensure `go test ./...` passes before submitting.

## Examples

New examples belong in [`examples/`](https://github.com/hirusha-adi/buttplug-go/tree/main/examples) as standalone `main` packages:

```bash
go run ./examples/your_example
```

## Questions

For protocol or server questions, see the upstream Buttplug community resources:

- [Buttplug Developer Guide](https://docs.buttplug.io)
- [Buttplug Discord](https://discord.buttplug.io)

For issues specific to this Go client, use [GitHub Issues](https://github.com/hirusha-adi/buttplug-go/issues).
