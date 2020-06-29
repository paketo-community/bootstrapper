#!/usr/bin/env bash
set -eu
set -o pipefail

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BootstrapDir="$(cd "${PROGDIR}/../bootstrapper" && pwd)"

function main() {
    pushd "${BootstrapDir}" > /dev/null
        if go test ./... -v -count=1; then
            printf "** GO Test Succeeded **"
        else
            printf "** GO Test Failed **"
        fi
    popd > /dev/null
}

main "${@:-}"
