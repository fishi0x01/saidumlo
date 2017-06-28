#!/bin/bash
set -e
set -o pipefail
set -o nounset

VAULT_VERSION="0.7.0"
VAULT_DOWNLOAD_BASE_URL="https://releases.hashicorp.com/vault/${VAULT_VERSION}"
VAULT_ZIP="vault_${VAULT_VERSION}_linux_amd64.zip"

################
# Download Vault
#
if [[ ! -x "./bin/vault" ]]; then
  mkdir -p ./bin
  echo "Downloading vault ${VAULT_VERSION}..."
  curl -Ls -o vault-checksums "${VAULT_DOWNLOAD_BASE_URL}/vault_${VAULT_VERSION}_SHA256SUMS"
  curl -Ls -O "${VAULT_DOWNLOAD_BASE_URL}/${VAULT_ZIP}"
  echo "Verifying checksum"
  grep ${VAULT_ZIP} vault-checksums | shasum -a 256 --check
  echo "Unzip vault ${VAULT_VERSION}"
  unzip ${VAULT_ZIP}
  chmod +x vault
  mv vault ./bin
  rm vault*
fi

########################
# Start vault dev server
#
bin/vault server -dev > /dev/null 2>&1 &
echo -n "$!" >> ./test/vault_pid

# wait for vault to boot
sleep 5
