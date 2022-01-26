# envflagset

[![License: APACHE2](https://img.shields.io/github/license/mikluko/envflagset.svg)](LICENSE)

A simple golang tool to set flag via environment variables inspired by [Go: Best Practices for Production Environments](http://peter.bourgon.org/go-in-production/#configuration)

# Features

- Set flag via environment variables.
- Auto mapping environment variables to flag. (e.g. `DATABASE_PORT` to `-database-port`)
- Customizable env - flag mapping support.
- Min length (default is 3) support in order to avoid parsing short flag.
- Show environment variable key in usage (-h).

# Basic Usage

### __Just keep it SIMPLE and SIMPLE and SIMPLE!__

Use `envflagset.Parse()` instead of `flag.Parse()`.

See `example` folder for complete examples.
