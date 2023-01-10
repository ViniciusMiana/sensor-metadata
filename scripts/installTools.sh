#!/bin/bash

set -e

while getopts "e:" opt
   do
     case $opt in
       e) multi+=("$OPTARG");;
     esac
done

if [ ${#multi[@]} -eq 0 ]; then
  echo "Usage: $0 -e swagger -e linter -e mockery -e swagger-cli -e trivy"
  echo "\t-e: Name of the exec you want to run. Available: swagger, linter, mockery, swagger-cli, trivy"
  exit 1
fi

if [ "$GO_EXEC_PATH" == "" ]; then GO_EXEC_PATH=".tmp_exec_path"; fi
if [ "$LINT_VERSION" == "" ]; then LINT_VERSION="1.50.1"; fi
if [ "$SWAGGER_VERSION" == "" ]; then SWAGGER_VERSION="0.29.0"; fi
if [ "$SWAGGER_CLI_VERSION" == "" ]; then SWAGGER_CLI_VERSION="4.0.4"; fi
if [ "$MOCKERY_VERSION" == "" ]; then MOCKERY_VERSION="2.14.0"; fi
if [ "$TRIVY_VERSION" == "" ]; then TRIVY_VERSION="0.32.1"; fi

tmpdir=`mktemp -d`
mkdir -p $GO_EXEC_PATH/bin

uname_os() {
  os=`uname -s | tr '[:upper:]' '[:lower:]'`
  case "$os" in
    cygwin_nt*) os="windows";;
    mingw*) os="windows";;
    msys_nt*) os="windows";;
  esac
  echo "$os"
}

uname_arch() {
  arch=`uname -m`
  case $arch in
    x86_64) arch="amd64";;
    x86) arch="386";;
    i686) arch="386";;
    i386) arch="386";;
    aarch64) arch="arm64";;
    armv5*) arch="armv5";;
    armv6*) arch="armv6";;
    armv7*) arch="armv7";;
    arm64) arch="arm64";;
  esac
  echo ${arch}
}

os=`uname_os`
arch=`uname_arch`

download_swagger() {
  echo "[STATUS] Downloading swagger $SWAGGER_VERSION $os $arch"
  url="https://github.com/go-swagger/go-swagger/releases/download/v$SWAGGER_VERSION/swagger_${os}_${arch}"
  rm -rf $GO_EXEC_PATH/bin/swagger
  curl -s -o $GO_EXEC_PATH/bin/swagger -L'#' "$url" && chmod +x $GO_EXEC_PATH/bin/swagger
}

download_golangci_lint() {
  echo "[STATUS] Downloading golangci-lint $LINT_VERSION $os $arch"
  url="https://github.com/golangci/golangci-lint/releases/download/v$LINT_VERSION/golangci-lint-${LINT_VERSION}-${os}-${arch}.tar.gz"
  curl -s -o $tmpdir/linter.tar.gz -L'#' "$url"

  (cd "${tmpdir}" && tar --no-same-owner -xzf "linter.tar.gz")
  rm -rf $GO_EXEC_PATH/bin/golangci-lint
  cp $tmpdir/golangci-lint-$LINT_VERSION-$os-${arch}/golangci-lint $GO_EXEC_PATH/bin/golangci-lint
}

download_swagger_cli() {
  echo "[STATUS] Downloading swagger-cli $SWAGGER_CLI_VERSION $os any"
  if [ "`command -v npm`" == "" ]; then
    echo "[WARN] ignoring swagger-cli download npm must be installed be download the swagger-cli"
    exit 0
  fi
  mkdir -p $GO_EXEC_PATH/bin/node
  cd $GO_EXEC_PATH/bin/node
  npm install --quiet swagger-cli@$SWAGGER_CLI_VERSION --prefix . &>/dev/null
  cd ..
  ln -sf node/node_modules/swagger-cli/swagger-cli.js ./swagger-cli

}

download_mockery() {
  m_arch=$arch
  if [ "$arch" == "amd64" ]; then
    m_arch="x86_64"
  fi
  url="https://github.com/vektra/mockery/releases/download/v${MOCKERY_VERSION}/mockery_${MOCKERY_VERSION}_${os}_${m_arch}.tar.gz"
  echo "[STATUS] Downloading mockery $MOCKERY_VERSION $os $arch"
  curl -s -o $tmpdir/mockery.tar.gz -L'#' "$url"
  (cd "${tmpdir}" && tar --no-same-owner -xzf "mockery.tar.gz")
  rm -rf $GO_EXEC_PATH/bin/mockery
  cp $tmpdir/mockery $GO_EXEC_PATH/bin/mockery
}

download_trivy() {
  m_os=$os
  if [ "$os" == "darwin" ]; then
    m_os="macos"
  fi
  m_arch=$arch
  if [ "$arch" == "amd64" ]; then
    m_arch="64bit"
  fi
  url="https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_${m_os}-${m_arch}.tar.gz"
  echo "[STATUS] Downloading trivy $TRIVY_VERSION $m_os $m_arch"
  curl -s -o $tmpdir/trivy.tar.gz -L'#' "$url"
  (cd "${tmpdir}" && tar --no-same-owner -xzf "trivy.tar.gz")
  rm -rf $GO_EXEC_PATH/bin/trivy
  cp $tmpdir/trivy $GO_EXEC_PATH/bin/trivy
  url="https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/html.tpl"
  curl -s -o $GO_EXEC_PATH/bin/trivy-html.tpl -L'#' "$url"
}

if [[ " ${multi[@]} " =~ " swagger " ]]; then
  tool="$GO_EXEC_PATH/bin/swagger"
  if [[ ! -f ${tool} ]] || [[ "$(${tool} version | awk 'NR==1{ print $2 }')" != "v${SWAGGER_VERSION}" ]]; then
    download_swagger
  fi
fi

if [[ " ${multi[@]} " =~ " linter " ]]; then
  tool="$GO_EXEC_PATH/bin/golangci-lint"
  if [[ ! -f ${tool} ]] || [[ "$(${tool} version | awk 'match($0,/[0-9.]{1,8}/){print substr($0,RSTART,RLENGTH)}')" != "${LINT_VERSION}" ]]; then
    download_golangci_lint
  fi
fi

if [[ " ${multi[@]} " =~ " mockery " ]]; then
  tool="$GO_EXEC_PATH/bin/mockery"
  if [[ ! -f ${tool} ]] || [[ "$(${tool} --quiet --version)" != "v${MOCKERY_VERSION}" ]]; then
    download_mockery
  fi
fi

if [[ " ${multi[@]} " =~ " swagger-cli " ]]; then
  tool="$GO_EXEC_PATH/bin/swagger-cli"
  if [[ ! -f ${tool} ]] || [[ "$(${tool} --version)" != "${SWAGGER_CLI_VERSION}" ]]; then
    download_swagger_cli
  fi
fi

if [[ " ${multi[@]} " =~ " trivy " ]]; then
  tool="$GO_EXEC_PATH/bin/trivy"
  if [[ ! -f ${tool} ]] || [[ "$(${tool} --version | awk 'NR==1{ print $2 }')" != "${TRIVY_VERSION}" ]]; then
    download_trivy
  fi
fi
