## 0.3.0 (Oct 01, 2017)

IMPROVEMENTS:

* Allow to specify chmod settings of local copy
* Allow wildcard `*` in paths to describe subtree
* Support binary data through base64 encoding

## 0.2.2 (Jul 26, 2017)

BUGFIXES:

* Create missing local nested directories on read (previously only single level created..)

## 0.2.1 (Jul 26, 2017)

IMPROVEMENTS:

* On read sdl creates the local directory specified in the .secrets.yml if it doesn't exist

## 0.2.0 (Jul 17, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Yaml config schema changed: Secret mappings are now separated from the `vaults` definition, which allows to use the same mapping against different vault backends

IMPROVEMENTS:

* You can now specify the `lease_ttl` of secrets inside a group (handy when using consul-template)
* Properly print a version string with `sdl version`
