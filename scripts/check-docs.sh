#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
required=(
    README.md CHANGELOG.md COMPATIBILITY.md CONTRIBUTING.md LICENSE
    SECURITY.md SUPPORT.md docs/README.md docs/adoption.md docs/api.md
    docs/architecture.md docs/compatibility.md docs/faq.md
    docs/operations.md docs/troubleshooting.md example_test.go
)

cd "${root}"
for path in "${required[@]}"; do
    test -s "${path}" || {
        printf 'missing required documentation: %s\n' "${path}" >&2
        exit 1
    }
done

while IFS=: read -r source match; do
    link="$(sed -E 's/.*\(([^)]+)\).*/\1/' <<<"${match}")"
    case "${link}" in
        ""|http://*|https://*|mailto:*) continue ;;
    esac
    path="${link%%#*}"
    fragment=""
    if [[ "${link}" == *#* ]]; then
        fragment="${link#*#}"
    fi
    target="$(dirname "${source}")/${path}"
    test -e "${target}" || {
        printf 'broken local documentation link: %s -> %s\n' "${source}" "${link}" >&2
        exit 1
    }
    if [[ -n "${fragment}" && "${target}" == *.md ]]; then
        found=false
        while IFS= read -r heading; do
            anchor="$(printf '%s' "${heading}" |
                tr '[:upper:]' '[:lower:]' |
                tr -cd '[:alnum:] _-' |
                sed -E 's/[[:space:]]+/-/g')"
            if [[ "${anchor}" == "${fragment}" ]]; then
                found=true
                break
            fi
        done < <(sed -nE 's/^#{1,6}[[:space:]]+//p' "${target}")
        [[ "${found}" == true ]] || {
            printf 'broken local documentation anchor: %s -> %s\n' \
                "${source}" "${link}" >&2
            exit 1
        }
    fi
done < <(rg -o --with-filename '\[[^]]+\]\([^)]+\)' \
    README.md COMPATIBILITY.md CONTRIBUTING.md SECURITY.md SUPPORT.md docs)

packages="$(go list ./...)"
while IFS= read -r package; do
    go doc "${package}" >/dev/null
done <<<"${packages}"

go test ./... -run '^Example' -count=1

printf 'required correlation documentation and examples are valid\n'
