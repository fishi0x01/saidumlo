#!/bin/bash
set -e
set -o pipefail
set -o nounset

#########
# Cleanup
#
# TODO: ensure that this always runs
kill $(cat ./test/vault_pid)
rm ./test/vault_pid
