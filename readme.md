# ![Starport](./assets/starport.jpg)

Starport is the all-in-one platform to build, launch and maintain any crypto application on a sovereign and secured blockchain. It is a developer-friendly interface to the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), the world's most widely-used blockchain application framework. Starport generates boilerplate code for you, so you can focus on writing business logic.

* [**Build a blockchain with Starport in a web-based IDE** (stable)](https://gitpod.io/#https://github.com/tendermint/starport/tree/master) or use [nightly version](https://gitpod.io/#https://github.com/tendermint/starport/)
* [Check out the latest features in v0.18](https://medium.com/tendermint/starport-v0-18-cosmos-sdk-updates-and-scaffolding-enhancements-5ea5654bcd0c)

## Quick start

Open Starport [in your browser](https://gitpod.io/#https://github.com/tendermint/starport/tree/master), or [install it](https://docs.starport.network/guide/install.html). Create and start a blockchain:

```bash
starport scaffold chain github.com/cosmonaut/mars

cd mars

starport chain serve
```

## Documentation

To learn how to use Starport, check out the [Starport Documentation](https://docs.starport.network). To learn more about how to build blockchain apps with Starport, see the [developer guide](https://docs.starport.network/guide/). To install Starport locally on GNU/Linux or macOS, follow [these steps](https://docs.starport.network/guide/install.html).

To learn more about building a JavaScript frontend for your Cosmos SDK blockchain, see [`tendermint/vue`](https://github.com/tendermint/vue).

## Questions

For questions and support, join the official [Starport Discord server](https://discord.gg/7fwqwc3afK). The issue list in this repo is exclusively for bug reports and feature requests.

## Cosmos SDK Compatibility

Blockchains created with Starport use the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) framework. To ensure the best possible experience, use the version of Starport that corresponds to the version of Cosmos SDK that you blockchain is built with. Unless noted otherwise, a row refers to a minor version and all associated patch versions.

| Starport | Cosmos SDK | Notes                                            |
| -------- | ---------- | ------------------------------------------------ |
| v0.18    | v0.44      |                                                  |
| v0.17    | v0.42      | `starport chain serve` works with v0.44.x chains |

To upgrade your blockchain to the newer version of Cosmos SDK refer to the [migration guide](https://docs.starport.network/migration/).

## Contributing

We welcome contributions from everyone. The `develop` branch contains the development version of the code. You can branch of from `develop` and create a pull request, or maintain your own fork and submit a cross-repository pull request. If you're not sure where to start check out [contributing.md](contributing.md) for our guidelines & policies for how we develop Starport. Thank you to all those who have contributed to Starport!

## Stay in touch

Starport is a free and open source product maintained by [Tendermint](https://tendermint.com). Follow us to get the latest updates!

- [Twitter](https://twitter.com/starportHQ)
- [Blog](https://medium.com/tendermint)
- [Jobs](https://tendermint.com/careers)
