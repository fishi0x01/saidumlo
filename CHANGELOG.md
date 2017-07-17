## 0.2.0 (Jul 17, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Yaml config schema changed: Secret mappings are now separated from the `vaults` definition, which allows to use the same mapping against different vault backends

IMPROVEMENTS:

* You can now specify the `lease_ttl` of secrets inside a group (handy when using consul-template)
* Properly print a version string with `sdl version`
