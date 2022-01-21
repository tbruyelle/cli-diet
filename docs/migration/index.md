---
order: 1
title: v0.19.2
parent:
  title: Migration
  order: 3
description: For chains that were scaffolded with Starport versions lower than v0.19.2, changes are required to use Starport v0.19.2. 
---

# Upgrading a blockchain to use Starport v0.19.2

Starport v0.19.2 comes with IBC v2.0.2. 

With Starport v0.19.2, the contents of the deprecated Starport Modules `tendermint/spm` repo are moved to the official Starport repo which introduces breaking changes.

To migrate your chain that was scaffolded with Starport versions lower than v0.19.2: 

1. IBC upgrade: Apply the changes that are introduced in PR [#1975](https://github.com/tendermint/starport/pull/1975/files) to your chain.
   
2. In your chain's `go.mod` file, remove `tendermint/spm` and add the v0.19.2 version of `tendermint/starport`. If your chain uses these packages, change the import paths as shown: 

    - `spm/ibckeeper` moved to `pkg/cosmosibckeeper`
    - `spm/cosmoscmd` moved to `pkg/cosmoscmd` 
    - `spm/openapiconsole` moved to `pkg/openapiconsole`
    - `testutil/sample` moved to `cosmostestutil/sample`

