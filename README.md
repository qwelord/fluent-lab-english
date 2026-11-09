# Fluent Lab

A writing and speaking practice studio built with a browser interface, Go HTTP server, and Rust analysis engine.

## Quick demo

Open `index.html` directly in a browser. The studio, prompt library, live feedback, speech playback, focus mode, and draft comparison work without setup. The browser does not save responses.

## Full stack run

Requires Go 1.22+ and Rust/Cargo.

```sh
cd analyzer && cargo build --release && cd ..
go run ./cmd/server
```

Open `http://localhost:8080`. The Go API passes draft text to the Rust CLI and returns validated JSON. The browser uses this endpoint when available and otherwise keeps its built-in analysis.

## Features

- Six discussion prompts with a topic filter and focused writing angles.
- Live word, character, sentence, variety, reading-time, and sentence-pace indicators.
- Rule-based revision checklist and precision-word suggestions.
- Browser speech playback, session timer, focus mode, copy, and side-by-side draft comparison.
- Go API with body-size and execution-time limits, content-type validation, and security headers.
- Rust analyzer with deterministic output and unit tests.

The analysis offers writing cues. It does not claim to grade proficiency or correct grammar.

## Verify

```sh
node --check app.js
go test ./...
cd analyzer && cargo test
```

GitHub Actions runs these checks on every push and pull request.
