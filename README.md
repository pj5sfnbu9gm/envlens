# envlens

A CLI tool to audit and diff environment variable configurations across multiple deployment targets.

## Installation

```bash
go install github.com/yourname/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envlens.git && cd envlens && go build -o envlens .
```

## Usage

Compare environment configurations between two deployment targets:

```bash
envlens diff production.env staging.env
```

Audit a single environment file for missing or undefined variables:

```bash
envlens audit --config .envlens.yaml production.env
```

Show all keys present in one target but missing in another:

```bash
envlens diff --missing-only production.env staging.env
```

### Example Output

```
[MISSING]  DATABASE_URL        found in production, not in staging
[MISMATCH] LOG_LEVEL           production=error  staging=debug
[OK]       PORT                matches across all targets
```

## Configuration

`envlens` can be configured via a `.envlens.yaml` file in your project root:

```yaml
targets:
  - name: production
    file: ./envs/production.env
  - name: staging
    file: ./envs/staging.env
ignore:
  - AWS_SESSION_TOKEN
  - CI_BUILD_ID
```

## License

MIT © [yourname](https://github.com/yourname)