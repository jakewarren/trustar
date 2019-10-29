# trustar
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/trustar/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/trustar)](https://goreportcard.com/report/github.com/jakewarren/trustar)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)

> a cli swiss army knife for working with Trustar

## Install

```
go get github.com/jakewarren/trustar
```

## Usage

### Installation
The program currently reads the API tokens from the `TRUSTAR_CLIENT_ID` and `TRUSTAR_CLIENT_SECRET` environment variables.

### Commands
| Command      | Subcommand   | Description                                  | Notes                                                                 |
|--------------|--------------|----------------------------------------------|-----------------------------------------------------------------------|
| autocomplete |              | Generates bash completion scripts            |                                                                       |
| help         |              | Help about any command                       |                                                                       |
| indicator    | find-reports | Find all correlated reports for an indicator | As of now, does not support wildcard searches                         |
| indicator    | search       | Search indicators                            | Not really that useful since you don't get much data back.            |
| reports      | search       | Search reports                               | Does support wildcard searches for indicators.                        |
| reports      | open         | Open the specified report(s) in your browser |                                                                       |
| list         |              | Lists enclaves                               |                                                                       |
| quota        |              | Print API request quota information          |                                                                       |
| token        |              | Print access token                           | Helper function to print your API token, useful for working with curl |
| whitelist    | add          | Add items to the whitelist                   |                                                                       |
| whitelist    | delete       | Delete items from the whitelist              |                                                                       |
| whitelist    | list         | List items in the whitelist                  |                                                                       |

## Roadmap

Current planned features are listed in [TODO.md](TODO.md). I may or may not get around to them.

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2019 Jake Warren

[changelog]: https://github.com/jakewarren/trustar/blob/master/CHANGELOG.md
