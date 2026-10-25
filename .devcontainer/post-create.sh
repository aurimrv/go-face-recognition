#!/usr/bin/env bash
# Runs once after the container is created. Never fails the container start.
set -u
git config --global --add safe.directory '*' || true
echo "== Tool versions"
go version
gcc --version | head -1
node --version
golangci-lint --version
command -v gomutants && echo "gomutants OK"
npx --no-install markdownlint-cli2 --version 2>/dev/null | head -1 || true
openspec --version
python3 -c 'import yaml; print("python3", __import__("sys").version.split()[0], "PyYAML", yaml.__version__)'
echo "== go mod download + go build ./..."
if [ ! -f go.mod ]; then
  echo "post-create: no go.mod in $(pwd), skipping build check."
elif go mod download && go build ./...; then
  echo "post-create: build OK. Next: go test ./... -cover"
else
  echo "post-create: WARNING build failed (on the original commit run 'go mod tidy' first)."
fi
exit 0
