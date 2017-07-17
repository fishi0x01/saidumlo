### Status
[![Master branch build status](https://api.travis-ci.org/fishi0x01/saidumlo.svg?branch=master)](https://travis-ci.org/fishi0x01/saidumlo.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/saidumlo)](https://goreportcard.com/report/github.com/fishi0x01/saidumlo)
[![Code Climate](https://codeclimate.com/github/fishi0x01/saidumlo/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/saidumlo)
[![Issue Count](https://codeclimate.com/github/fishi0x01/saidumlo/badges/issue_count.svg)](https://codeclimate.com/github/fishi0x01/saidumlo)

# SaiDumLo

SaiDumLo aims to be a client site secret management tool primarily designed for local development. 

Currently, SaiDumLo only interacts as a wrapper for HashiCorp's [vault](https://www.vaultproject.io/) client. 
Vault is awesome, but lacks an easy configurable config file to synch your local ops repo with the vault secrets. 
I always find myself writing and maintaining different `Makefile` commands for different secrets of different stages (qa/staging/live..).
SaiDumLo lets you easily define and manage different secret groups like `qa` or `prod` in a single yaml config file. 

Example **.secrets.yml:**
```
---
vaults:
  vaultA:
    default: true
    address: "http://127.0.0.1:8200"
    bin: "my/path/to/vault"
    auth:
      method: "github"
      credential_file: "my/path/to/credentials"

  vaultB:
    address: "https://vault.b.int.company.local:8200"
    bin: "my/path/to/vault"
    auth:
      method: "github"
      credential_file: "my/path/to/credentials"

secrets:
  qa:
    lease_ttl: "1h"
    mappings:
    - local: "local/path/to/qa-foo"
      vault: "secret/qa/qa-foo"
    - local: "local/path/to/qa-bar"
      vault: "secret/qa/qa-bar"

  prod:
    mappings:
    - local: "local/path/to/prod-foo"
      vault: "secret/prod/prod-foo"

```

SaiDumLo handles reads/writes of your secret groups by using the vault client. 
Using `sdl read qa` synchronizes your local `qa` secrets with the current ones from the default vault (`vaultA`). 
`sdl -b vaultB write prod` writes your local `prod` secrets to `vaultB`. 

Before reading/writing SaiDumLo authenticates with the vault by using the specified method. 
In the example `.secrets.yml` the `github` method is used, which requires a github auth token from your account. 
The auth credentials file must contain key/value pairs of the necessary parameters, e.g., for github:

**github.credentials.auth:**
```
token=<my-github-token>
```

For the `userpass` mechanism it should be:

**userpass.credentials.auth:**
```
username=<my-user>
password=<my-password>
```

Consult the vault [auth documentation](https://www.vaultproject.io/docs/auth/index.html) to see which parameters need to be specified in the credentials file for your auth method. 

**NOTE: Do not forget to add the auth credential file to your .gitignore!**

### Build and Test

```
make verify
```

Tested with vault `0.7.0` on Ubuntu Xenial.

