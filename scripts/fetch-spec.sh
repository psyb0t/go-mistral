#!/bin/bash
# Downloads the Mistral AI OpenAPI document into api/. The spec is not
# committed; gen.go fetches it as the first step of `go generate`.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

readonly SPEC_URL="https://raw.githubusercontent.com/mistralai/platform-docs-public/main/openapi.yaml"
readonly DOWNLOAD_TIMEOUT_SECONDS=60
readonly SPEC_OPENAPI="openapi.yaml"
SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

usage() {
	printf 'Usage: %s <%s>\n' "$(basename "$0")" "$SPEC_OPENAPI" >&2
}

# Prints info.version from a YAML OpenAPI document: the first `version:` line
# indented by two spaces inside the top-level `info:` block.
spec_version() {
	awk '
		/^info:/ { in_info = 1; next }
		/^[^ ]/ { in_info = 0 }
		in_info && /^  version:/ { sub(/^  version:[ ]*/, ""); gsub(/["\047]/, ""); print; exit }
	' "$1"
}

main() {
	local spec="${1:-}"
	if [[ "$spec" != "$SPEC_OPENAPI" ]]; then
		log ERROR "unknown spec name=${spec}"
		usage
		exit 2
	fi

	mkdir -p "$SPEC_DIR"
	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.${spec}.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}'" EXIT

	log INFO "downloading spec=${spec}"
	wget -q \
		-T "$DOWNLOAD_TIMEOUT_SECONDS" \
		-O "$tmp" \
		"$SPEC_URL"

	if ! grep -Eq '^openapi: *["\047]?3\.' "$tmp"; then
		log ERROR "downloaded file is not an OpenAPI 3 document spec=${spec}"
		exit 1
	fi

	chmod 0644 "$tmp"
	mv -- "$tmp" "${SPEC_DIR}/${spec}"
	log INFO "saved spec=${spec} version=$(spec_version "${SPEC_DIR}/${spec}")"
}

main "$@"
