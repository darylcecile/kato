# Kato Language Design

> Core idea: **YAML-like readability, TOML-like predictability, JSON-like cheap parsing, with first-class metadata for agents.**

## Design Goals

Kato should be:

- **Human friendly:** braces for structure, but low visual noise.
- **Easy to parse:** no indentation semantics, no implicit typing surprises, no complex anchors/references.
- **Cheap to deserialize:** maps directly to structs/objects, single-pass parseable, predictable token stream.
- **Agent friendly:** comments, schema hints, stable structure, explicit values, no invisible meaning from whitespace.
- **Extensible:** directives and namespaces without breaking old parsers.

## Core Syntax

### Objects

```kato
server {
  port: 3000
  host: "localhost"
}
```

Equivalent JSON:

```json
{
  "server": {
    "port": 3000,
    "host": "localhost"
  }
}
```

### Arrays

```kato
regions: ["eu-west-1", "us-east-1"]

plugins [
  "auth"
  "payments"
  "search"
]
```

Both comma and newline-separated arrays are allowed. Trailing commas are allowed.

### Scalars

```kato
name: "Daryl"
age: 25
enabled: true
mode: production
empty: null
```

Bare words are strings unless they are reserved literals:

- `true`
- `false`
- `null`
- `on`
- `off`

So this:

```kato
env: production
```

means `{ "env": "production" }`. No YAML-style "yes means true" nonsense.

## Comments

```kato
# line comment

server {
  port: 3000  # inline comment
}
```

Triple-slash comments are agent-visible documentation:

```kato
/// Agent-visible documentation
/// This explains what the setting is for.
timeout: 30s
```

Regular comments (`#`) are ignored. Triple-slash comments (`///`) may be preserved in an AST/doc model for tooling and agents.

## Units as Typed Literals

```kato
timeout: 30s
memory: 512MiB
rateLimit: 100/min
```

These deserialize as structured values, not vague strings:

```json
{
  "timeout": { "$unit": "s", "value": 30 }
}
```

Strict mode could require schemas to define valid units.

## Environment and References

Keep interpolation explicit and function-like:

```kato
database {
  url: env("DATABASE_URL")
  replicaUrl: env("REPLICA_DATABASE_URL", fallback: null)
}
```

No magic `${...}` by default.

References are also explicit:

```kato
shared {
  region: "eu-west-1"
}

deploy {
  region: ref("shared.region")
}
```

## Directives

Top-level `@` directives are for metadata and extensions:

```kato
@version 1
@schema "https://example.com/schemas/app.kato"
@profile production
```

Unknown directives can be ignored or rejected depending on parser mode.

## Namespaces for Extension

```kato
@use github.actions as gha

gha.workflow {
  name: "Deploy"
  trigger: "push"
}
```

This gives standards a clean extension model without polluting the base language.

## Agent-Friendly Metadata

A special optional block:

```kato
@agent {
  purpose: "Configure the production deployment"
  owner: "platform-team"
  safeToEdit: ["features", "server.timeout"]
  doNotEdit: ["database.url"]
}
```

Not needed by runtimes, but useful for AI tooling, editors, automation, and code review bots.

## Includes

Make includes explicit, static, and boring:

```kato
@include "./shared.kato"
@include "./production.secrets.kato" as secrets
```

Includes are resolved before validation, but after lexical parsing. No dynamic includes unless a runtime deliberately enables it.

## Profiles / Overlays

Instead of multiple files full of duplication:

```kato
server {
  port: 3000
  workers: 2
}

@profile production {
  server {
    workers: 8
  }
}
```

Merged result for `production`:

```kato
server {
  port: 3000
  workers: 8
}
```

Merge rules:

- **Objects:** deep merge
- **Scalars:** replace
- **Arrays:** replace by default

Optional array operations:

```kato
plugins += ["analytics"]
plugins -= ["debug-toolbar"]
```

## Schema (also written in Kato)

```kato
@schemaVersion 1

type AppConfig {
  app: App
  server: Server
  database: Database
}

type Server {
  host: string = "localhost"
  port: int range(1, 65535)
  timeout: duration = 30s
}

type Database {
  url: secret<string>
  pool?: Pool
}
```

## The Big Rule

The format should have **one obvious AST**.

No indentation meaning. No implicit dates. No surprise booleans. No anchors. No hidden merge semantics. No magic interpolation.

That is what keeps it cheap, standardisable, and agent-friendly.
