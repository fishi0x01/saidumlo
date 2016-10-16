### Status
[![Master branch build status](https://api.travis-ci.org/fishi0x01/saidumlo.svg?branch=master)](https://travis-ci.org/fishi0x01/saidumlo.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishi0x01/saidumlo)](https://goreportcard.com/report/github.com/fishi0x01/saidumlo)

# SaiDumLo

SaiDumLo is a secret management tool primarily designed for local development. 
Working locally with different secret groups, each group having a different encryption/decryption mechanisms, can easily become cumbersome. 
For instance, when writing ansible roles you might have some secrets encrypted with `ansible-vault` and others, such as key files encrypted with `openssl`. 
Further, you might have complex permission rules (not everyone might be allowed to decrypt all the secrets from the repo) in which case multiple secret keys are necessary which adds even more complexity. 

SaiDumLo lets you define your secret groups and encryption mechanisms inside a single `secrets.yml` file. 
Basically, the `secrets.yml` file is like a Makefile for encrypting and decrypting your secrets. 

## Build

```
go get gopkg.in/yaml.v2
go get gopkg.in/alecthomas/kingpin.v2 
go get github.com/fatih/color
go build -o sdl src/*
```


