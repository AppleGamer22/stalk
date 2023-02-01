# `stalk`
[![Go Reference](https://pkg.go.dev/badge/github.com/AppleGamer22/stalk.svg)](https://pkg.go.dev/github.com/AppleGamer22/stalk) [![Test](https://github.com/AppleGamer22/stalk/actions/workflows/test.yml/badge.svg)](https://github.com/AppleGamer22/stalk/actions/workflows/test.yml) [![CodeQL](https://github.com/AppleGamer22/stalk/actions/workflows/codeql.yml/badge.svg)](https://github.com/AppleGamer22/stalk/actions/workflows/codeql.yml) [![Release](https://github.com/AppleGamer22/stalk/actions/workflows/release.yml/badge.svg)](https://github.com/AppleGamer22/stalk/actions/workflows/release.yml) [![Update Documentation](https://github.com/AppleGamer22/stalk/actions/workflows/tag.yml/badge.svg)](https://github.com/AppleGamer22/stalk/actions/workflows/tag.yml)

## Description
`stalk` is a cross-platform CLI utility for file-watching.

## Why This Name?
This name is simply a stupid ~~pun~~, therefore **I do not condone and do not promote stalking** (excluding stalking fictional individuals for the purposes of a [CTF](https://en.wikipedia.org/wiki/Capture_the_flag_(cybersecurity)) challenge).

## Installation
### Arch Linux Distributions
* [`yay`](https://github.com/Jguer/yay):
```bash
yay -S stalk-bin
```
* [`paru`](https://github.com/morganamilo/paru):
```bash
paru -S stalk-bin
```

### macOS
* [Homebrew Tap](https://github.com/AppleGamer22/homebrew-tap):
```bash
brew install AppleGamer22/tap/stalk
```

### Windows (working progress)
* [`winget`](https://github.com/microsoft/winget-cli):
```bash
winget install AppleGamer22.stalk
```

### Other
* `go`:
	* Does not ship with:
		* a manual page
		* pre-built shell completion scripts
```
go install github.com/AppleGamer22/stalk
```

## Functionality
### `watch` Sub-command
The arguments have to be existing files, that the user running the command can read.

#### `-c`/`--command` Flag
This flag specifies the command to be stopped and  re-run after each file system event.

#### `-v`/`--verbose` Flag
This flag enables more detailed logs to be printed regarding file system events.

### `wait` Sub-command
The arguments have to be existing files, that the user running the command can read. This command is meant to be integrated inside other scripts, which is why **it only waits for a single file system event** before terminating.


### `version` Sub-command
#### `-v`/`--verbose` Flag
* If this flag is provided, the following details are printed to the screen:
	1. semantic version number
	2. commit hash
	3. Go compiler version
	4. processor architecture & operating system
* Otherwise, only the semantic version number is printed.

## Dependencies
The only dependencies `stalk` requires are the OS-level dependencies of the [`github.com/fsnotify/fsnotify`](https://github.com/fsnotify/fsnotify) library.

## Common Contributor Routines
### Testing
Running the following command will run `go test` on the cmd and session sub-modules:
```bash
make test
```
### Building From Source
#### Development
* Using the following `make` command will save a `stalk` binary with the last version tag and the latest git commit hash:
```bash
make debug
```

#### Release
* Using the following [GoReleaser](https://github.com/goreleaser/goreleaser) command with a version `git` tag and a clean `git` state:
```bash
goreleaser build --clean
```
* All release artificats will stored in the `dist` child directory in the codebase's root directory:
	* compressed package archives with:
		* a `stalk` binary
		* manual page
		* shell completion scripts
	* checksums
	* change log

## Copyright
`stalk` is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3, or (at your option) any later version.

`stalk` is distributed in the hope that it will be useful, but **WITHOUT ANY WARRANTY**; without even the implied warranty of **MERCHANTABILITY** or **FITNESS FOR A PARTICULAR PURPOSE**.  See the GNU General Public License for more details.