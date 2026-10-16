#!/bin/bash
# Makes the downloaded OpenAPI document generate with oapi-codegen.
#
# oapi-codegen turns a property `const` into a named type <Schema><Property>.
# The spec has boolean properties pinned with `const: true` or `const: false`
# (a `stream` flag on a request schema), and the type named for one of them is
# already taken by a component schema (the streaming variant of that request),
# so generation stops on a duplicate type name. A boolean const carries no type
# information worth keeping, so every `const: true` and `const: false` line is
# dropped and those properties become plain booleans.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR
readonly SOURCE_SPEC="${SPEC_DIR}/openapi.yaml"
readonly OUTPUT_SPEC="${SPEC_DIR}/openapi.patched.yaml"
readonly BOOLEAN_CONST_PATTERN='^[[:space:]]+const:[[:space:]]+(true|false)[[:space:]]*$'

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

main() {
	if [[ ! -f "$SOURCE_SPEC" ]]; then
		log ERROR "source spec missing path=${SOURCE_SPEC}"
		exit 1
	fi

	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.openapi.patched.yaml.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}'" EXIT

	local dropped
	dropped="$(grep -cE "$BOOLEAN_CONST_PATTERN" "$SOURCE_SPEC" || true)" # grep -c exits 1 on zero matches

	grep -vE "$BOOLEAN_CONST_PATTERN" "$SOURCE_SPEC" >"$tmp"

	chmod 0644 "$tmp"
	mv -- "$tmp" "$OUTPUT_SPEC"
	log INFO "patched spec dropped_boolean_consts=${dropped}"
}

main "$@"
