# Kato

A proof-of-concept notation language — an alternative to YAML and JSON — designed to be both human and agent friendly.

The name comes from the Mauritian Creole *"kato ver"* (green box).

## Why?

- **JSON** is easy to parse but painful for humans.
- **YAML** is nice until it becomes spooky: indentation semantics, implicit booleans, anchors, surprising typing.
- **TOML** is readable but gets awkward for nested structures.

Kato's niche:

```
JSON's predictability
+ TOML's explicitness
+ YAML's readability
+ first-class tooling/agent metadata
```

## Quick Look

```kato
@version 1

server {
  host: "0.0.0.0"
  port: 3000
  timeout: 30s
}

database {
  url: env("DATABASE_URL")
  pool {
    min: 2
    max: 10
  }
}

features {
  search: on
  bookings: off
}
```

## The Big Rule

**One obvious AST.** No indentation meaning. No implicit dates. No surprise booleans. No anchors. No hidden merge semantics. No magic interpolation.

## What's in this repo

- [`spec/`](./spec/) — Language design and grammar
- [`examples/`](./examples/) — Sample `.kato` files
- `parsers/` — Reference implementations (coming)
  - `parsers/go/` — Go parser
  - `parsers/ts/` — TypeScript parser

## File Extension

`.kato`

## Status

🚧 Early development — spec is evolving.

## License

TBD
