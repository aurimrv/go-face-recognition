# Dev Container for go-face-recognition (Go + cgo + dlib)
# Multi-arch: linux/amd64 (Intel/AMD PCs, Intel Macs) and linux/arm64 (Apple Silicon).
# Base: Debian 12 "bookworm" (same distro as the repo's golang:1.22.1 image);
# bookworm still ships libatlas-base-dev, which provides libcblas for go-face.
FROM debian:bookworm

ARG GO_VERSION=1.26.10
ARG NODE_VERSION=22.23.4
ARG GOLANGCI_LINT_VERSION=2.15.0
# gomutants: pinned. v0.5.1 is the version used in the class material.
# NOTE: from v0.7.0 the module moved to github.com/gomutants/gomutants, so
# "go install github.com/szhekpisov/gomutants@latest" no longer works.
ARG GOMUTANTS_VERSION=v0.5.1
ARG GOBCO_VERSION=v1.3.4
ARG OPENSPEC_VERSION=1.3.1
ARG MARKDOWNLINT_CLI2_VERSION=0.24.0
ARG GO_FACE_VERSION=v0.0.0-20210630145111-0c14797b4d0e
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8

# 1) OS packages: toolchain, cgo deps of go-face (dlib, BLAS/CBLAS, LAPACK, JPEG),
#    git, curl/wget (Cursor downloads its server with wget), python3 + PyYAML.
RUN apt-get -o Acquire::Retries=5 update \
    && apt-get -o Acquire::Retries=5 install -y --no-install-recommends \
        build-essential pkg-config \
        libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg-dev \
        git curl wget ca-certificates xz-utils unzip sudo less procps \
        python3 python3-yaml \
    && rm -rf /var/lib/apt/lists/*

# 2) Go from the official go.dev tarball (arch picked automatically).
RUN set -eux; arch="$(dpkg --print-architecture)"; \
    case "$arch" in amd64|arm64) ;; *) echo "unsupported arch $arch"; exit 1;; esac; \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz" -o /tmp/go.tgz; \
    tar -C /usr/local -xzf /tmp/go.tgz; rm /tmp/go.tgz; /usr/local/go/bin/go version

# 3) Node.js (needed by OpenSpec CLI >= 20.19 and markdownlint-cli2).
RUN set -eux; arch="$(dpkg --print-architecture)"; \
    case "$arch" in amd64) n=x64;; arm64) n=arm64;; esac; \
    curl -fsSL "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-${n}.tar.xz" -o /tmp/node.txz; \
    tar -C /usr/local --strip-components=1 -xJf /tmp/node.txz; rm /tmp/node.txz; \
    node --version; npm --version; \
    npm install -g "@fission-ai/openspec@${OPENSPEC_VERSION}" "markdownlint-cli2@${MARKDOWNLINT_CLI2_VERSION}"; \
    npm cache clean --force

# 4) golangci-lint v2 (prebuilt binary); the project standard (config "version: 2").
RUN set -eux; arch="$(dpkg --print-architecture)"; \
    curl -fsSL "https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-${arch}.tar.gz" \
      | tar -xz -C /tmp; \
    install -m 0755 "/tmp/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-${arch}/golangci-lint" /usr/local/bin/; \
    rm -rf /tmp/golangci-lint-*; golangci-lint --version

# 5) Non-root user.
RUN groupadd --gid "${USER_GID}" "${USERNAME}" \
    && useradd --uid "${USER_UID}" --gid "${USER_GID}" -m -s /bin/bash "${USERNAME}" \
    && echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > "/etc/sudoers.d/${USERNAME}" \
    && chmod 0440 "/etc/sudoers.d/${USERNAME}"

# 5b) C++ compiler wrapper for cgo. go-face hard-codes "#cgo CXXFLAGS: -march=native".
#     On arm64 (Apple Silicon VMs, QEMU) GCC's native detection can enable SHA3
#     instructions (eor3) that the assembler then rejects:
#       "Error: selected processor does not support `eor3 ...'"
#     The wrapper drops -march=native on aarch64 only (GCC default armv8-a).
#     On amd64, -march=native is kept (verified to work).
RUN printf '%s\n' '#!/bin/bash' \
      'if [ "$(uname -m)" = "aarch64" ]; then' \
      '  args=(); for a in "$@"; do [ "$a" = "-march=native" ] || args+=("$a"); done' \
      '  exec g++ "${args[@]}"' \
      'fi' \
      'exec g++ "$@"' > /usr/local/bin/cgo-g++ \
    && chmod 0755 /usr/local/bin/cgo-g++
ENV CXX=/usr/local/bin/cgo-g++

ENV GOPATH=/home/${USERNAME}/go \
    GOTOOLCHAIN=local \
    CGO_ENABLED=1
ENV PATH=/usr/local/go/bin:${GOPATH}/bin:${PATH}
# Login shells (/etc/profile resets PATH) must also see Go and the Go tools.
RUN printf '%s\n' 'export GOPATH="$HOME/go"' 'export GOTOOLCHAIN=local' \
      'case ":$PATH:" in *":/usr/local/go/bin:"*) ;; *) export PATH="/usr/local/go/bin:$HOME/go/bin:$PATH";; esac' \
      > /etc/profile.d/10-go.sh && chmod 0644 /etc/profile.d/10-go.sh

USER ${USERNAME}

# 6) Go tools + pre-warm the build cache with the cgo package go-face (dlib),
#    so the students' first "go build" does not recompile it (slow, RAM hungry).
RUN set -eux; \
    go install "github.com/szhekpisov/gomutants@${GOMUTANTS_VERSION}"; \
    go install "github.com/rillig/gobco@${GOBCO_VERSION}"; \
    mkdir -p /tmp/warm && cd /tmp/warm; \
    go mod init warm; \
    go get "github.com/Kagami/go-face@${GO_FACE_VERSION}"; \
    go build github.com/Kagami/go-face; \
    cd / && rm -rf /tmp/warm; \
    gomutants --help >/dev/null 2>&1 || true; \
    ls -l "${GOPATH}/bin"

ENV NPM_CONFIG_UPDATE_NOTIFIER=false
WORKDIR /workspaces
CMD ["sleep", "infinity"]
