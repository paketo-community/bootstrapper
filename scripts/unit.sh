#!/usr/bin/env bash

set -e
set -u
set -o pipefail

function main() {
  if go test . -v -count=1; then
      echo "** GO Test Succeeded **"
  else
      echo "** GO Test Failed **"
      exit 1
  fi
}

main "${@:-}"
