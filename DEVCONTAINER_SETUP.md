# Dev Container setup: go-face-recognition in Cursor

This guide gets every student the same Linux toolchain inside a container, whether
they use macOS (Apple Silicon or Intel), Windows or Linux. Cursor runs on your
computer. Go, cgo, dlib and all the class tools run inside the container.

What is inside the container (Debian 12 "bookworm", amd64 and arm64):

| Tool | Version | Command |
|---|---|---|
| Go (cgo on) | 1.26.10 | `go` |
| gcc/g++, pkg-config | 12.2 | used by cgo |
| dlib, BLAS/CBLAS (ATLAS), LAPACK, libjpeg | Debian packages | used by `github.com/Kagami/go-face` |
| gomutants | v0.5.1 (same as the class material) | `gomutants` |
| gobco | v1.3.4 | `gobco` |
| golangci-lint v2 | 2.15.0 (uses the repo's `.golangci.yml`) | `golangci-lint` |
| Node.js | 22.23.4 | `node`, `npx` |
| OpenSpec CLI | 1.3.1 | `openspec` |
| markdownlint-cli2 | 0.24.0 | `markdownlint-cli2` or `npx markdownlint-cli2` |
| Python 3 + PyYAML | 3.11 / 6.0 | `python3` |
| git, curl, wget, sudo | | |

The user inside the container is `vscode` (not root). `~/go/bin` is on `PATH`.
The image already contains a compiled copy of go-face (dlib), so your first
`go build` takes seconds instead of a minute or more.

---

## 0. Before class (do this at home; the download is about 1.7 GB)

You need about **15 GB free disk**, and the container should get **at least 4 GB of RAM**
(8 GB is better).

### macOS (Apple Silicon M1–M4 or Intel)

1. Download **Docker Desktop** from <https://www.docker.com/products/docker-desktop/>.
   Pick **"Mac with Apple chip"** or **"Mac with Intel chip"**
   ( menu → About This Mac shows which one you have).
2. Open the `.dmg`, drag Docker to Applications, start it, and accept the terms.
   You do not need to sign in.
3. Open Docker Desktop → **Settings → Resources**. Set **CPUs ≥ 4** and **Memory ≥ 8 GB**
   (4 GB is the minimum), then click **Apply & restart**.
4. Check it works. Open Terminal and run:
   ```bash
   docker run --rm hello-world
   ```
   You should see `Hello from Docker!`.

### Windows 10/11

1. Open **PowerShell as Administrator** and run `wsl --install`. Reboot if asked.
   (If WSL is already installed, run `wsl --update`.)
2. Install **Docker Desktop for Windows** from
   <https://www.docker.com/products/docker-desktop/>. During setup, keep
   **"Use WSL 2 instead of Hyper-V"** checked.
3. Start Docker Desktop. Open **Settings → General** and check that
   **"Use the WSL 2 based engine"** is on.
4. Memory: with WSL 2, Docker uses the WSL limits. By default that is half of your RAM,
   which is fine on a 16 GB laptop. On an 8 GB laptop, create `%UserProfile%\.wslconfig`
   with:
   ```ini
   [wsl2]
   memory=6GB
   processors=4
   ```
   Then run `wsl --shutdown` and start Docker Desktop again.
5. Check it works. In PowerShell, run `docker run --rm hello-world`.
6. Tip: the build is much faster if the repo is cloned **inside WSL**
   (`\\wsl$\Ubuntu\home\<you>\...`) rather than under `C:\`. Cloning under `C:\` also works.

*Alternative without Docker Desktop:* install Docker Engine inside a WSL 2 Ubuntu
distro (see the Linux steps). Then open the folder through Cursor's WSL connection
and run "Reopen in Container" from there. Use this only if Docker Desktop is not
allowed on your machine.

### Linux (Ubuntu/Debian/Fedora…)

1. Install Docker Engine by following <https://docs.docker.com/engine/install/>.
   Quick path for Ubuntu/Debian:
   ```bash
   curl -fsSL https://get.docker.com | sudo sh
   sudo usermod -aG docker $USER      # then log out and log in again
   ```
2. Check it works: `docker run --rm hello-world` (no `sudo`).
3. On native Linux there are no resource settings. Just keep at least 4 GB of RAM free.

### All systems: Cursor + git

1. Install **Cursor** from <https://cursor.com/download> and update it to the latest
   version (Help → Check for Updates).
2. Install the **Dev Containers** extension: open the Extensions view
   (`Ctrl+Shift+X` / `Cmd+Shift+X`), search for **Dev Containers**, and install the one
   **published by Anysphere** (`anysphere.remote-containers`). It is Cursor's own
   version, rebuilt in-house. **Uninstall Microsoft's `ms-vscode-remote.remote-containers`**
   if you have it, because it causes the known "cannot install Cursor server" errors.
   Sources: <https://forum.cursor.com/t/new-in-house-extensions-c-c-ssh-devcontainers-wsl-python/94531>
   and <https://forum.cursor.com/t/launching-with-devcontainer-failed-to-install-cursor-server/151404>.
3. Install git (macOS: `xcode-select --install`; Windows: <https://git-scm.com/download/win>
   or use git inside WSL).

---

## 1. Get the code

1. On GitHub, open <https://github.com/aurimrv/go-face-recognition> and click **Fork**.
   Untick "Copy the main branch only" if you want all the class branches.
2. Clone your fork:
   ```bash
   git clone https://github.com/<your-user>/go-face-recognition.git
   cd go-face-recognition
   ```
3. Make sure the folder contains `.devcontainer/devcontainer.json`.
   If it does not (for example, on a branch without it), unzip `devcontainer-kit.zip`
   from the course page into the repo root, so that you have `go-face-recognition/.devcontainer/...`.
   An untracked `.devcontainer/` folder stays in place when you `git switch` between branches.
   To keep it out of your commits, run:
   ```bash
   echo ".devcontainer/" >> .git/info/exclude
   ```

## 2. Open it in the container

1. Start **Docker Desktop**, and wait until it says *Engine running*.
2. In Cursor, go to **File → Open Folder…** and choose the `go-face-recognition` folder.
3. Cursor shows the popup *"Folder contains a Dev Container configuration file… Reopen in Container"*.
   Click it. If there is no popup, press `F1` (or `Ctrl/Cmd+Shift+P`) and run
   **Dev Containers: Reopen in Container**.
4. Wait. The first build downloads the base image and tools and compiles go-face:
   - about **1.5–3 minutes** on a fast connection and a recent laptop (amd64 measured from scratch: 92 s; the go-face/dlib compile alone takes about 57 s on 8 cores);
   - **5–10 minutes** on slow Wi-Fi or older laptops.
   To watch progress, click **"show log"** in the bottom-right notification.
   Later opens take only a few seconds.
5. When the green/blue box in the bottom-left corner says **Dev Container: go-face-recognition…**,
   you are inside. The `postCreateCommand` prints the tool versions and runs `go build ./...`
   once.
6. Open a terminal: **Terminal → New Terminal**. The prompt looks like
   `vscode@<id>:/workspaces/go-face-recognition$`.

## 3. Check that everything works (copy-paste)

```bash
go version
go vet ./...
go test ./... -cover
gomutants ./internal/entity
golangci-lint --version
npx markdownlint-cli2 --version
openspec --version
python3 -c 'import yaml; print("PyYAML OK", yaml.__version__)'
go run ./cmd
```

Expected output (measured on branch `11_final_mutation_score`; the numbers on other branches differ):

```text
go version go1.26.10 linux/amd64            # linux/arm64 on Apple Silicon
                                            # go vet: no output = OK
ok  github.com/andre-ols/go-face-recognition/cmd                5.4s  coverage: 100.0% of statements   (on a laptop: 5–20 s)
ok  github.com/andre-ols/go-face-recognition/internal/entity    0.0s  coverage: 100.0% of statements
ok  github.com/andre-ols/go-face-recognition/internal/usecases  3.9s  coverage: 100.0% of statements
gomutants v0.5.1 ... Killed: 45  Lived: 2  Timed out: 38  Efficacy: 95.74%   (about 13 s)
golangci-lint has version 2.15.0 built with go1.27.0 ...
markdownlint-cli2 v0.24.0 (markdownlint v0.42.0)   (followed by "Linting: 0 files", which is normal)
1.3.1
PyYAML OK 6.0
Found 8 faces ... Person: Trump ... Wrote images/result.jpg
```

On the **original commit** (`ad32735`, the starting point of the from-scratch exercise), first run:

```bash
go mod tidy          # adds the missing go-face requirement to go.mod
go vet ./...         # OK
go test ./... -cover # "coverage: 0.0% of statements" for every package (no tests yet)
golangci-lint run ./... # 10 issues (6 errcheck, 1 ineffassign, 3 staticcheck)
```

Notes:

- `gomutants` writes `mutation-report.json` and `.gomutants-cache.json`, and `go run ./cmd`
  rewrites `images/result.jpg`. These are tracked files: commit them or run
  `git restore mutation-report.json .gomutants-cache.json images/result.jpg`.
- `Timed out: 38` in gomutants is expected. These are the same numbers as in the committed report on branch 11.
- **Lint:** the repo ships a v2 `.golangci.yml` (errcheck, govet, ineffassign, staticcheck, unused),
  so everyone runs `golangci-lint run ./...` with the same rules. Don't create or replace it.
- The container sees your repo folder live. Edits in Cursor show up in the container immediately,
  and `git` works in both places.

---

## 4. Troubleshooting

| Symptom | Fix |
|---|---|
| `Cannot connect to the Docker daemon` / "Docker is not running" | Start Docker Desktop and wait for *Engine running*. On Linux: `sudo systemctl start docker`, and check that you are in the `docker` group (`groups`). |
| Build stops with `Killed`, `signal: killed`, or `c++: fatal error: Killed signal terminated program cc1plus` | Out of memory while compiling dlib/go-face. Give Docker **≥ 6–8 GB** (Docker Desktop → Resources, or `.wslconfig` on Windows), close the browser and other heavy apps, then run **Dev Containers: Rebuild Container**. Inside the container you can also limit parallel builds: `go build -p 1 ./...`. |
| First build is very slow | Normal on slow Wi-Fi (about 1.7 GB of downloads). Build at home before class. Don't cancel it: the next attempt reuses the finished layers. |
| Apple Silicon: is it emulated? | No. The image is built natively as **linux/arm64** (`go version` prints `linux/arm64`). Don't force `--platform linux/amd64`, because emulated amd64 is 10–20× slower. |
| `-march=native` (go-face) | go-face hard-codes `#cgo CXXFLAGS: ... -march=native`. **Apple clang on macOS** rejects it (Kagami/go-face#86), which is why we don't build natively on macOS. On **Linux arm64**, GCC's native detection can enable SHA3 instructions that the assembler rejects (`Error: selected processor does not support 'eor3 ...'`). We reproduced this while building for arm64. The image therefore sets `CXX=/usr/local/bin/cgo-g++`, a small wrapper that removes `-march=native` on arm64 only. Don't unset `CXX`. If you see the `eor3` error, you are not using this image: run **Rebuild Container**. |
| "Failed to install Cursor server" / endless "Opening Remote…" | Update Cursor, and make sure the extension is **Anysphere** Dev Containers (not Microsoft's). The image already has `wget`, which Cursor needs. Turn off VPN and try again. Run **Dev Containers: Rebuild Container**. Check the logs in **View → Output → Dev Containers**. |
| Corporate proxy / university network blocks downloads | Set the proxy in Docker Desktop (Settings → Resources → Proxies). For the build, add `"build": { "args": { "HTTP_PROXY": "http://proxy:port", "HTTPS_PROXY": "http://proxy:port" } }` to `devcontainer.json`, and add `"remoteEnv": { "HTTP_PROXY": "...", "HTTPS_PROXY": "..." }` for the Go/npm downloads inside. In Cursor, enable **Settings → "Disable HTTP/2"** if the connection hangs. Phone hotspot is a quick test. |
| `go: ... toolchain not available` | The container sets `GOTOOLCHAIN=local` on purpose. Go 1.26.10 builds this module (`go 1.22.1`). Don't change `go.mod`'s `go` line. |
| `no required module provides package github.com/Kagami/go-face` | You are on the original commit. Run `go mod tidy`. |
| Windows: very slow file access / tests | Clone the repo inside WSL (`~/go-face-recognition` in Ubuntu), open it from Cursor via WSL, then Reopen in Container. |
| `dubious ownership` from git | Already handled by the post-create script. Otherwise run `git config --global --add safe.directory '*'`. |
| Disk full | `docker system prune` (removes old containers and images). |

---

## 5. Fallback: GitHub Codespaces + Cursor Remote-SSH

Use this only if Docker can't run on your laptop. The same `.devcontainer` folder
works in Codespaces, and the free personal quota is enough for the class.

1. `gh codespace ssh` needs an SSH server in the container. In your fork, add this
   top-level entry to `.devcontainer/devcontainer.json`, then commit and push:
   ```json
   "features": { "ghcr.io/devcontainers/features/sshd:1": {} },
   ```
   (Verified: the image builds with this feature and `sshd` is present. It is left out of
   the default config to keep local builds minimal.)
2. Install the GitHub CLI (<https://cli.github.com>), then run:
   ```bash
   gh auth login
   gh auth refresh -h github.com -s codespace
   gh codespace create -R <your-user>/go-face-recognition -b <branch> -m standardLinux32gb
   gh codespace ssh --config > ~/.ssh/codespaces_config
   ```
   (`-m` picks a 4-core / 16 GB machine.)
3. Add `Include ~/.ssh/codespaces_config` at the top of `~/.ssh/config`.
4. In Cursor, install **Remote - SSH (Anysphere)**, run **Remote-SSH: Connect to Host…**,
   and pick the `cs.<name>...` host. Then open `/workspaces/go-face-recognition`.
5. Run the same checks as in section 3.
6. Stop the codespace afterwards: `gh codespace stop`. It stops on its own after 30 minutes idle.

The image builds with the sshd feature (verified with the devcontainer CLI), but the Codespaces creation and Remote-SSH steps were not run.
