---
order: 3
description: Steps to install Starport on your local computer.
---

# Install Starport

You can run Starport in a web-based Gitpod IDE or you can install Starport on your local computer.

## Prerequisites

Local Starport installation requires the follow software be installed and running:

- [Golang >=1.16](https://golang.org/)

- [Node.js >=12.19.0](https://nodejs.org/)


## Installing Starport with cURL

```
curl https://get.starport.network/starport! | bash
```

The latest `starport` binary is downloaded from the Github repo and installed in `/usr/local/bin`.

To install previous versions of the precompiled `starport` binary or customize the installation process, see [Starport installer docs](https://github.com/allinbits/starport-installer) on GitHub.

### Write permission

Starport installation requires write permission to the `/usr/local/bin/` directory. If the installation fails because you do not have write permission to `/usr/local/bin/`, run the following command:

```
curl https://get.starport.network/starport | bash
```

Then run this command to move the `starport` executable to `/usr/local/bin/`:

```
sudo mv starport /usr/local/bin/
```

## Installing Starport on macOS with Homebrew

```
brew install tendermint/tap/starport
```

## Build from source

```
git clone https://github.com/tendermint/starport --depth=1
make -C starport install
```

## Summary

- To setup a local development environment, install Starport locally on your computer.
- Install Starport by fetching the binary using cURL, Homebrew, or by building from source.
- The latest version is installed by default. You can install previous versions of the precompiled `starport` binary.
