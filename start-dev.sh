#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="${ROOT_DIR}/frontend"

BACKEND_PID=""
FRONTEND_PID=""

cleanup() {
  if [[ -n "${BACKEND_PID}" ]] && kill -0 "${BACKEND_PID}" 2>/dev/null; then
    kill "${BACKEND_PID}" 2>/dev/null || true
  fi
  if [[ -n "${FRONTEND_PID}" ]] && kill -0 "${FRONTEND_PID}" 2>/dev/null; then
    kill "${FRONTEND_PID}" 2>/dev/null || true
  fi
}

trap cleanup EXIT INT TERM

printf "\n➡️  Starting Go backend on http://localhost:8080 ...\n"
(
  cd "${ROOT_DIR}"
  go run .
) &
BACKEND_PID=$!

printf "\n➡️  Starting frontend dev server on http://localhost:5173 ...\n"
if [[ ! -d "${FRONTEND_DIR}/node_modules" ]]; then
  (
    cd "${FRONTEND_DIR}"
    npm install
  )
fi
(
  cd "${FRONTEND_DIR}"
  npm run dev
) &
FRONTEND_PID=$!

EXIT_CODE=0
while true; do
  if ! kill -0 "${BACKEND_PID}" 2>/dev/null; then
    wait "${BACKEND_PID}" || EXIT_CODE=$?
    break
  fi

  if ! kill -0 "${FRONTEND_PID}" 2>/dev/null; then
    wait "${FRONTEND_PID}" || EXIT_CODE=$?
    break
  fi

  sleep 1
done

exit "${EXIT_CODE}"
