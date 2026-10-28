#!/bin/bash
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

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

SARIF_OUT="${SARIF_OUT:-sec.sarif}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

log INFO "scanning the current worktree for secrets"
gitleaks_dir_rc=0
gitleaks dir --no-banner --no-color --redact \
	--gitleaks-ignore-path /work/.gitleaksignore \
	--report-format sarif \
	--report-path "$tmp/gitleaks-dir.sarif" \
	/work || gitleaks_dir_rc=$?

log INFO "scanning Git history for secrets"
gitleaks_git_rc=0
gitleaks git --no-banner --no-color --redact \
	--gitleaks-ignore-path /work/.gitleaksignore \
	--report-format sarif \
	--report-path "$tmp/gitleaks-git.sarif" \
	/work || gitleaks_git_rc=$?

log INFO "running govulncheck"
# SARIF mode exits 0 even with findings, so it only writes the report; the plain
# text run below is what gates, exiting non-zero on a reachable vulnerability.
govulncheck -format sarif ./... >"$tmp/govulncheck.sarif"
gv_rc=0
govulncheck ./... || gv_rc=$?

log INFO "running semgrep"
# --error gates on a finding; --sarif writes the report; --metrics=off keeps the
# scan from phoning home. One run does both.
sg_rc=0
semgrep_targets=(
	.github
	Dockerfile.dev
	Makefile
	scripts
)
# client.gen.go is left out of the target list: it is oapi-codegen output, so a
# finding there can't be fixed in place, and the only rule it trips is unsafe-
# deserialization-interface on JSON decoding into the interface{} fields the
# spec leaves untyped (free-form objects), which is plain encoding/json, not a
# gadget-chain decoder. --exclude does not apply to a file named on the command
# line, and at 60k lines the file also makes several rules hit the per-file
# timeout, so whether those findings appeared depended on runner speed.
for go_file in ./*.go; do
	[[ "$go_file" == *.gen.go ]] || semgrep_targets+=("$go_file")
done
semgrep scan --config p/golang --config p/security-audit \
	--no-git-ignore --sarif --output "$tmp/semgrep.sarif" \
	--metrics=off --error "${semgrep_targets[@]}" || sg_rc=$?

log INFO "merging SARIF"
# Code scanning rejects a file holding two runs from the same tool, so the
# gitleaks directory and history runs are folded into one run per tool.
jq -s '{version:"2.1.0","$schema":"https://json.schemastore.org/sarif-2.1.0.json",runs:([.[].runs // []] | add | group_by(.tool.driver.name) | map(.[0] + {results: (map(.results // []) | add)}))}' \
	"$tmp/gitleaks-dir.sarif" \
	"$tmp/gitleaks-git.sarif" \
	"$tmp/govulncheck.sarif" \
	"$tmp/semgrep.sarif" >"$SARIF_OUT"

if [ "$gitleaks_dir_rc" -ne 0 ] || [ "$gitleaks_git_rc" -ne 0 ] ||
	[ "$gv_rc" -ne 0 ] || [ "$sg_rc" -ne 0 ]; then
	log ERROR "security findings, gitleaks_dir=${gitleaks_dir_rc}, gitleaks_git=${gitleaks_git_rc}, govulncheck=${gv_rc}, semgrep=${sg_rc}, report=${SARIF_OUT}"
	exit 1
fi
log INFO "security scan complete with no findings"
