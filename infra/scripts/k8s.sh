#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CLUSTER_NAME="${CLUSTER_NAME:-sys-called}"
NAMESPACE="${NAMESPACE:-sys-called}"
RELEASE="${RELEASE:-sys-called}"
TAG="${TAG:-dev-$(date +%Y%m%d%H%M%S)}"
# IMAGE=<repo:tag> deploys an image that is already built (by CI, for
# instance) instead of building one from the working tree.
IMAGE="${IMAGE:-}"
CHART_DIR="$ROOT_DIR/infra/helm/sys-called"
KIND_CONFIG="$ROOT_DIR/infra/kind/cluster.yaml"
APP_URL="http://localhost:30080"

log() {
  printf '\033[1;34m==>\033[0m %s\n' "$*"
}

fail() {
  printf '\033[1;31merro:\033[0m %s\n' "$*" >&2
  exit 1
}

require() {
  local missing=()
  for tool in "$@"; do
    command -v "$tool" >/dev/null 2>&1 || missing+=("$tool")
  done
  if ((${#missing[@]})); then
    fail "ferramentas não encontradas: ${missing[*]}"
  fi
}

cluster_exists() {
  kind get clusters 2>/dev/null | grep -qx "$CLUSTER_NAME"
}

ensure_cluster() {
  if cluster_exists; then
    log "cluster kind '$CLUSTER_NAME' já existe"
  else
    log "criando cluster kind '$CLUSTER_NAME'"
    kind create cluster --name "$CLUSTER_NAME" --config "$KIND_CONFIG" --wait 120s
  fi
  kubectl config use-context "kind-$CLUSTER_NAME" >/dev/null
}

build_and_load_image() {
  local output
  if [[ -z "$IMAGE" ]]; then
    IMAGE="sys-called:$TAG"
    log "build de $IMAGE"
    docker build -q -t "$IMAGE" --build-arg VERSION="$TAG" "$ROOT_DIR" >/dev/null
  fi
  log "carregando $IMAGE no cluster"
  if ! output=$(kind load docker-image "$IMAGE" --name "$CLUSTER_NAME" 2>&1); then
    printf '%s\n' "$output" >&2
    fail "não foi possível carregar $IMAGE no cluster"
  fi
}

deploy() {
  log "instalando o chart '$RELEASE' no namespace '$NAMESPACE'"
  helm upgrade --install "$RELEASE" "$CHART_DIR" \
    --namespace "$NAMESPACE" --create-namespace \
    --values "$CHART_DIR/values-dev.yaml" \
    --set image.repository="${IMAGE%:*}" \
    --set image.tag="${IMAGE##*:}" \
    "$@" \
    --wait --timeout 10m
}

wait_for_app() {
  log "esperando o app responder em $APP_URL"
  for _ in $(seq 1 30); do
    if curl -sf "$APP_URL/readyz" >/dev/null; then
      return 0
    fi
    sleep 2
  done
  fail "o app não respondeu em $APP_URL"
}

cmd_up() {
  require docker kind kubectl helm curl
  ensure_cluster
  build_and_load_image
  deploy "$@"
  wait_for_app
  log "pronto: $APP_URL (usuários de desenvolvimento com a senha senha123)"
}

cmd_deploy() {
  require docker kind kubectl helm
  cluster_exists || fail "o cluster '$CLUSTER_NAME' não existe; rode '$0 up' primeiro"
  kubectl config use-context "kind-$CLUSTER_NAME" >/dev/null
  build_and_load_image
  deploy "$@"
}

cmd_test() {
  require helm curl
  kubectl config use-context "kind-$CLUSTER_NAME" >/dev/null
  log "helm test"
  helm test "$RELEASE" --namespace "$NAMESPACE" --logs
  log "smoke test em $APP_URL"
  "$ROOT_DIR/scripts/smoke-test.sh" "$APP_URL"
}

cmd_status() {
  require kubectl
  kubectl --context "kind-$CLUSTER_NAME" -n "$NAMESPACE" get pods,svc,pvc
}

cmd_logs() {
  require kubectl
  kubectl --context "kind-$CLUSTER_NAME" -n "$NAMESPACE" logs -f "deploy/$RELEASE" --all-containers --prefix
}

cmd_uninstall() {
  require helm kubectl
  log "removendo o release '$RELEASE' (os volumes e o secret são mantidos)"
  helm uninstall "$RELEASE" --namespace "$NAMESPACE" --kube-context "kind-$CLUSTER_NAME"
}

cmd_down() {
  require kind
  if cluster_exists; then
    log "apagando o cluster kind '$CLUSTER_NAME'"
    kind delete cluster --name "$CLUSTER_NAME"
  else
    log "cluster '$CLUSTER_NAME' não existe, nada a fazer"
  fi
}

usage() {
  cat <<USAGE
Uso: $(basename "$0") <comando> [argumentos extras do helm]

Comandos:
  up         cria o cluster kind (se preciso), faz o build da imagem e instala o chart
  deploy     refaz o build da imagem e atualiza o release num cluster existente
  test       roda o helm test e o smoke test contra o cluster
  status     mostra pods, services e volumes
  logs       acompanha os logs do app
  uninstall  remove o release, mantendo volumes e secret
  down       apaga o cluster kind inteiro

Variáveis: CLUSTER_NAME=$CLUSTER_NAME NAMESPACE=$NAMESPACE RELEASE=$RELEASE TAG=<gerada>
           IMAGE=<repo:tag> usa uma imagem pronta em vez de fazer o build
USAGE
}

main() {
  local command="${1:-}"
  shift || true
  case "$command" in
    up) cmd_up "$@" ;;
    deploy) cmd_deploy "$@" ;;
    test) cmd_test ;;
    status) cmd_status ;;
    logs) cmd_logs "$@" ;;
    uninstall) cmd_uninstall ;;
    down) cmd_down ;;
    -h | --help | help | "") usage ;;
    *) usage; fail "comando desconhecido: $command" ;;
  esac
}

main "$@"
