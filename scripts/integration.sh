#!/usr/bin/env bash
set -eu
set -o pipefail

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly IntegrationDir="$(cd "${PROGDIR}/../integration" && pwd)"

function main() {
    pushd "${IntegrationDir}" > /dev/null
        if go test ./... -v -count=1; then
            printf "** GO Test Succeeded **"
        else
            printf "** GO Test Failed **"
            exit 1
        fi
    popd > /dev/null
}

main "${@:-}"
