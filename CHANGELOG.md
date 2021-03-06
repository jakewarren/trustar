# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.3] - 2019-11-12
### Fixed
- fixed an issue where OAuth tokens were being requested when the autocomplete command was executed

### Changed
- update deps

## [0.1.2] - 2019-11-01
### Fixed
- fixed an issue where a report assigned to multiple enclaves would be printed multiple times

### Changed
- changed environment variable names for clarity

## [0.1.1] - 2019-11-01
### Added
- resolves enclave ids to enclave names when searching reports

### Changed
- added a verbose flag to get additional debug information
- added a log flag to write request/response data from TruStar's API

## [0.1.0] - 2019-10-29
- Initial Release

[unreleased]: https://github.com/jakewarren/trustar/compare/v0.1.3...HEAD
[0.1.3]: https://github.com/jakewarren/trustar/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/jakewarren/trustar/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/jakewarren/trustar/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/jakewarren/trustar/releases/tag/v0.1.0