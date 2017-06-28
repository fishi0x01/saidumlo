### Status
[![Master branch build status](https://api.travis-ci.org/fishi0x01/saidumlo.svg?branch=master)](https://travis-ci.org/fishi0x01/saidumlo.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/saidumlo)](https://goreportcard.com/report/github.com/fishi0x01/saidumlo)
[![Code Climate](https://codeclimate.com/github/fishi0x01/saidumlo/badges/gpa.svg)](https://codeclimate.com/github/fishi0x01/saidumlo)
[![Issue Count](https://codeclimate.com/github/fishi0x01/saidumlo/badges/issue_count.svg)](https://codeclimate.com/github/fishi0x01/saidumlo)

# SaiDumLo

SaiDumLo aims to be a client site secret management tool primarily designed for local development. 

Currently, SaiDumLo only interacts as a wrapper for HashiCorp's [vault](https://www.vaultproject.io/) client. 
It lets you easily define and manage different secret groups like `qa` or `prod` in a single yaml config file. 

Example **.secrets.yml:**
```
---
vault_address: "http://127.0.0.1:8200"
vault_bin: "my/path/to/vault"

groups:
  qa:
    secrets:
    - local: "local/path/to/qa-foo"
      vault: "secret/qa/qa-foo"
    - local: "local/path/to/qa-bar"
      vault: "secret/qa/qa-bar"

  prod:
    secrets:
    - local: "local/path/to/prod-foo"
      vault: "secret/prod/prod-foo"
```

SaiDumLo handles reads/writes of your secret groups by using the vault client. 

### Build and Test

```
make deps
make build
make verify
```

Tested with vault `0.7.0` on Ubuntu Xenial.

### TODO

* versioning / rollback capabilities

