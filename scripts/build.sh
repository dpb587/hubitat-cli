#!/bin/bash

set -euo pipefail

set -eu

repository=github.com/dpb587/hubitat-cli
cli=$( basename "${repository}" )
version="${1:-0.0.0}"

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/.."

commit=$( git rev-parse HEAD | cut -c-10 )

if [[ $( git clean -dnx | wc -l ) -gt 0 ]] ; then
  commit="${commit}+dirty"

  if [[ "${version}" != "0.0.0" ]]; then
    echo "ERROR: building an official version requires a clean repository"
    git clean -dnx

    exit 1
  fi
fi

built=$( date -u +%Y-%m-%dT%H:%M:%S+00:00 )

echo "build: cli=${cli}, version=${version}, commit=${commit}, built=${built}"

mkdir -p tmp/build
rm -fr tmp/build/*

export CGO_ENABLED=0

function build {
  os="$1"
  arch="$2"

  name=$cli-$version-$os-$arch

  if [ "$os" == "windows" ]; then
    name=$name.exe
  fi

  echo "build: ${name}"
  GOOS=$os GOARCH=$arch go build \
    -ldflags "
      -s -w
      -X ${repository}/cmd/cmdflags.VersionName=${version}
      -X ${repository}/cmd/cmdflags.VersionCommit=${commit}
      -X ${repository}/cmd/cmdflags.VersionBuilt=${built}
    " \
    -o tmp/build/${name} \
    .
}

build darwin amd64
build darwin arm64
build linux amd64
build linux arm64
build windows amd64

cd tmp/build

echo 'build: checksums (sha256)'
sha256sum * | tee ../build-checksums.txt

cd ..
(
  if [ -e "../docs/releases/v${version}.md" ]; then
    sed '1{/^---$/!q;};1,/^---$/d' "../docs/releases/v${version}.md" | sed -e '2,$b' -e '/^$/d'
    echo
  fi

  echo "**Assets (sha256)**"
  echo ""
  sed 's/^/    /' build-checksums.txt
) > build-notes.md
