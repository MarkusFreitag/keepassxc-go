# keepassxc-go

This repository contains a library as well as a basic CLI tool to interact with KeepassXC using the provided unix socket.

## Installation

To install it, you can download the binary or one of the packages (deb, rpm) from the [releases](https://github.com/MarkusFreitag/keepassxc-go/releases/latest).

## Usage

### CLI tool
The CLI tool currently is quite limited as it only provides a way
to search for an url.
```
$ ./ci-build/keepassxc-go --help
interact with keepassxc via unix-socket

Usage:
  keepassxc-go [command]

Available Commands:
  get-logins  query info for the specified url
  help        Help about any command

Flags:
  -h, --help             help for keepassxc-go
  -p, --profile string   Only necessary if keystore contains multiple profiles

Use "keepassxc-go [command] --help" for more information about a command.
```
```
$ ./ci-build/keepassxc-go get-logins --help
query info for the specified url

Usage:
  keepassxc-go get-logins URL [flags]

Flags:
      --all         show all matches otherwise only the first will be printed
  -h, --help        help for get-logins
      --plaintext   print out the password - BE CAREFUL

Global Flags:
  -p, --profile string   Only necessary if keystore contains multiple profiles
```
