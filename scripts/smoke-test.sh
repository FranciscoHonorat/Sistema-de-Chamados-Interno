#!/usr/bin/env bash
# Exercises a running Sys-Called through its public entry point: probes, the
# SPA, authentication (with the refresh cookie), the ticket lifecycle and the
# integration between the modules (support agents reach the tickets module
# through the outbox and the event bus).
#
#   scripts/smoke-test.sh [base-url]        (default http://localhost:8000)
#
# With SMOKE_USERNAME / SMOKE_PASSWORD it logs in as that account instead of
# the demo administrator, which is how it runs against environments without
# the demo users. Such an account must be an administrator.
#
# SMOKE_READONLY=1 stops before logging in, so nothing is written: that is the
# mode used right after a production deploy.
set -euo pipefail

BASE_URL="${1:-${BASE_URL:-http://localhost:8000}}"
BASE_URL="${BASE_URL%/}"
USERNAME="${SMOKE_USERNAME:-admin}"
PASSWORD="${SMOKE_PASSWORD:-senha123}"
TIMEOUT="${SMOKE_TIMEOUT:-60}"

pass() { printf '\033[32m✔\033[0m %s\n' "$*"; }
fail() { printf '\033[31m✘ %s\033[0m\n' "$*" >&2; exit 1; }

json_field() { # json_field <field>  (reads JSON from stdin)
  grep -o "\"$1\":\"[^\"]*\"" | head -n1 | cut -d'"' -f4
}

http_code() { curl -s -o /dev/null -w '%{http_code}' "$@"; }

echo "Smoke testing $BASE_URL"

deadline=$((SECONDS + TIMEOUT))
until [ "$(http_code "$BASE_URL/readyz")" = "200" ]; do
  [ $SECONDS -lt $deadline ] || fail "/readyz did not answer 200 within ${TIMEOUT}s"
  sleep 2
done
pass "ready"

[ "$(http_code "$BASE_URL/healthz")" = "200" ] || fail "/healthz is not 200"
pass "alive"

page=$(curl -sf "$BASE_URL/chamados")
grep -q '<div id="app">' <<<"$page" || fail "the SPA is not served on deep links"
pass "frontend served on deep links"

headers=$(curl -s -D - -o /dev/null "$BASE_URL/")
grep -qi '^content-security-policy:' <<<"$headers" || fail "missing Content-Security-Policy"
grep -qi '^x-frame-options: DENY' <<<"$headers" || fail "missing X-Frame-Options"
pass "security headers"

[ "$(http_code "$BASE_URL/api/tickets/tickets")" = "401" ] || fail "the API accepted an anonymous call"
pass "API refuses anonymous calls"

if [ "${SMOKE_READONLY:-0}" = "1" ]; then
  echo "Read-only smoke tests passed."
  exit 0
fi

cookies=$(mktemp)
trap 'rm -f "$cookies"' EXIT
login=$(curl -sf -c "$cookies" -X POST "$BASE_URL/api/employees/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}") || fail "login as $USERNAME failed"
token=$(json_field access_token <<<"$login")
[ -n "$token" ] || fail "login returned no access token"
grep -q '/api/employees/auth' "$cookies" || fail "refresh cookie is not scoped to /api/employees/auth"
pass "login (refresh cookie scoped to the auth routes)"

auth=(-H "Authorization: Bearer $token")

refreshed=$(curl -sf -b "$cookies" -c "$cookies" -X POST "$BASE_URL/api/employees/auth/refresh") || fail "refresh failed"
[ -n "$(json_field access_token <<<"$refreshed")" ] || fail "refresh returned no access token"
pass "session refresh"

deadline=$((SECONDS + TIMEOUT))
until [ "$(curl -sf "${auth[@]}" "$BASE_URL/api/tickets/responsibles" | grep -o '"id"' | wc -l)" -ge 1 ]; do
  [ $SECONDS -lt $deadline ] || fail "no support agent reached the tickets module"
  sleep 2
done
pass "support agents synced between modules"

ticket=$(curl -sf "${auth[@]}" -X POST "$BASE_URL/api/tickets/tickets" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Smoke test","description":"Criado pelo smoke test","priority":"Low"}') || fail "opening a ticket failed"
id=$(json_field ticket_id <<<"$ticket")
[ -n "$id" ] || fail "no ticket id in $ticket"
pass "ticket opened ($id)"

curl -sf "${auth[@]}" -X POST "$BASE_URL/api/tickets/tickets/$id/assign/auto" | grep -q '"assignee_id"' \
  || fail "automatic assignment failed"
pass "ticket assigned automatically"

curl -sf "${auth[@]}" "$BASE_URL/api/tickets/tickets/$id" | grep -q "\"ticket_id\":\"$id\"" || fail "ticket detail failed"
curl -sf "${auth[@]}" "$BASE_URL/api/tickets/tickets" | grep -q "$id" || fail "ticket missing from the list"
pass "ticket listed and detailed"

curl -sf -b "$cookies" -X POST "$BASE_URL/api/employees/auth/logout" -o /dev/null || fail "logout failed"
[ "$(http_code -b "$cookies" -X POST "$BASE_URL/api/employees/auth/refresh")" = "401" ] || fail "refresh still works after logout"
pass "logout revokes the session"

echo "All smoke tests passed."
