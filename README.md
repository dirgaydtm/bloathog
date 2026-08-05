# Bloathog

<img align="right" width="150" height="150" title="Bloathog logo" src="./docs/logo.svg">

A real-time resource monitor for Node.js, Bun, and Deno dev servers.

* **Deep.** Finds every hidden process running under the hood.
* **Fast.** Written in pure Go. It monitors your heaviest dev servers with zero overhead.
* **Interactive.** A slick terminal interface that just works.

[![Release](https://img.shields.io/github/v/release/dirgaydtm/bloathog?style=flat-square)](https://github.com/dirgaydtm/bloathog/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/dirgaydtm/bloathog?style=flat-square)](https://goreportcard.com/report/github.com/dirgaydtm/bloathog)
[![License](https://img.shields.io/github/license/dirgaydtm/bloathog?style=flat-square)](https://github.com/dirgaydtm/bloathog/blob/main/LICENSE)
[![Issues](https://img.shields.io/github/issues/dirgaydtm/bloathog?style=flat-square)](https://github.com/dirgaydtm/bloathog/issues)


## Installation

### with NPM

```bash
npm i -D bloathog
# or
npx bloathog
# or install globally
npm install -g bloathog
```

### with Go

```bash
go install github.com/dirgaa/bloathog/cmd/bloathog@latest
```

### with cURL

```bash
curl -sSL https://raw.githubusercontent.com/dirgaa/bloathog/main/install.sh | sh
```

## Usage

If you are in a JavaScript/TypeScript project directory, simply run:

```bash
bloathog

# or manually
bloathog bun run dev
bloathog pnpm run start
bloathog deno task dev
```

## Why Bloathog?

Ever noticed how tools like Next.js, Webpack, or Vite love to secretly spawn a bunch of background workers?

Standard monitoring tools like `top` or `htop` only show memory usage for a single PID. Trying to figure out the actual memory cost of your dev server is a nightmare. Bloathog solves this by natively tracking the entire process tree down to the last leaf, giving you the real, combined numbers. No more hidden memory.

## Contributing

Want to help make Bloathog better? Check out [CONTRIBUTING.md](CONTRIBUTING.md) to see how to set up your environment and shoot a PR.

Got a question or a cool idea? Just [open an issue](https://github.com/dirgaydtm/bloathog/issues) or hit me up at [dirgayuditama6@gmail.com](mailto:dirgayuditama6@gmail.com).
