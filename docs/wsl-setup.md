# Silkwave dev environment — Windows + WSL2 setup

A runbook for setting up (or rebuilding) the Silkwave development environment on Windows using WSL2. Designed to take ~30 minutes from a clean Windows install.

## Outcome

After completing this doc you will have:

- **Ubuntu 24.04 LTS** running under WSL2, configured with systemd and a non-root default user.
- The Silkwave repo at **`~/projects/silkwave`** (Linux ext4 — fast file watchers, fast `pnpm install`, fast `go build`).
- **VS Code on Windows** editing the project via the Remote-WSL extension.
- **zsh + zinit + Starship + Windows Terminal** as the shell stack.
- **`mise`** managing Node 22 + Go 1.25; **corepack** managing pnpm.
- **Native Docker Engine inside WSL** (no Docker Desktop).
- **ArangoDB** via `docker compose`, started on demand.
- **Claude Code, Sandcastle, `gh`** all running inside WSL.
- A fresh **ed25519 SSH key** for GitHub, stored only inside WSL.

## Prerequisites

- Windows 11 (any edition) with WSL2 already enabled at the kernel level.
- A password manager (Bitwarden, 1Password, KeePass) holding `ARANGO_ROOT_PASSWORD`, `ANTHROPIC_API_KEY`, AWS keys, etc. Values come from there.
- Access to the GitHub account `brandon-v-seeters` (to add the new SSH key).

---

## Step 1 — Confirm pre-existing backups

The old Arch distro was backed up to `C:\Users\brand\backups\wsl-archlinux-2026-05-06\` (SSH keys, gitconfig, zshrc). Do not delete this folder yet.

**Calendar reminder for 2026-11-06**: revisit and delete if no service has needed the old key in 6 months.

## Step 2 — Tear down the old environment

From PowerShell (admin not required):

```powershell
wsl --shutdown
wsl --unregister archlinux
wsl --unregister podman-machine-default
```

Then uninstall Docker Desktop via **Settings → Apps & Features → Docker Desktop → Uninstall**. Reboot.

After reboot, `wsl -l -v` should be empty (or show only `docker-desktop` until Docker Desktop's leftovers are cleared on its next launch — won't happen since it's uninstalled).

## Step 3 — Configure WSL2 globally

Create `C:\Users\brand\.wslconfig`:

```ini
[wsl2]
memory=12GB
swap=4GB
localhostForwarding=true
```

## Step 4 — Install Ubuntu 24.04

```powershell
wsl --install -d Ubuntu-24.04
```

When the new shell appears, set username `brand` and a password.

## Step 5 — Configure `/etc/wsl.conf` inside Ubuntu

From inside the Ubuntu shell:

```bash
sudo tee /etc/wsl.conf > /dev/null <<'EOF'
[boot]
systemd=true

[user]
default=brand

[automount]
options="metadata,umask=22,fmask=11"
EOF
```

From PowerShell, recycle the distro so the changes take effect:

```powershell
wsl --shutdown
wsl
```

Verify you're not root: `whoami` → `brand`. Verify systemd: `systemctl --version` should print without error.

## Step 6 — System packages

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y build-essential git curl wget unzip jq zsh ca-certificates rsync
```

## Step 7 — Set zsh as default shell

```bash
chsh -s "$(which zsh)"
```

Exit and reopen the shell. You should now be in zsh (prompt may look different).

## Step 8 — Install zinit + plugins

```bash
bash -c "$(curl --fail --show-error --silent --location https://raw.githubusercontent.com/zdharma-continuum/zinit/HEAD/scripts/install.sh)"
```

The installer auto-adds the bootstrap to `~/.zshrc`. Append the two plugins:

```bash
cat >> ~/.zshrc <<'EOF'

# zinit plugins
zinit light zsh-users/zsh-autosuggestions
zinit light zdharma-continuum/fast-syntax-highlighting
EOF
```

## Step 9 — Install Starship

```bash
curl -sS https://starship.rs/install.sh | sh -s -- -y
echo 'eval "$(starship init zsh)"' >> ~/.zshrc
```

## Step 10 — Install `mise`

```bash
curl https://mise.run | sh
echo 'eval "$(~/.local/bin/mise activate zsh)"' >> ~/.zshrc
exec zsh
```

Verify: `mise --version`.

## Step 11 — Install Docker Engine + Compose plugin

Use Docker's official apt repo (newer than Ubuntu's `docker.io` package):

```bash
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo tee /etc/apt/keyrings/docker.asc > /dev/null
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

sudo usermod -aG docker "$USER"
sudo systemctl enable --now docker
```

Recycle the distro from PowerShell so the new group membership takes effect:

```powershell
wsl --shutdown
wsl
```

Smoke test: `docker run --rm hello-world`.

## Step 12 — Install `gh` CLI

```bash
sudo mkdir -p -m 755 /etc/apt/keyrings
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null
sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update && sudo apt install -y gh

gh auth login   # GitHub.com → HTTPS → authenticate via browser
```

## Step 13 — Install Claude Code

```bash
curl -fsSL https://claude.ai/install.sh | bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
exec zsh
claude   # follow prompts to authenticate
```

## Step 14 — Generate SSH key for GitHub

```bash
ssh-keygen -t ed25519 -C "brandon@aizy.nl" -f ~/.ssh/id_ed25519 -N ""
cat ~/.ssh/id_ed25519.pub
```

Paste the public key at <https://github.com/settings/keys> → **New SSH key** → name "WSL Ubuntu", type "Authentication Key".

Verify:

```bash
ssh -T git@github.com
# Hi brandon-v-seeters! You've successfully authenticated...
```

## Step 15 — Configure git

```bash
git config --global user.name "Brandon van Seeters"
git config --global user.email "brandonvanseeters@gmail.com"
git config --global init.defaultBranch main
git config --global pull.rebase false
git config --global core.autocrlf input
```

## Step 16 — Install VS Code + Remote-WSL extension

On the **Windows** side:

1. Install **VS Code** from <https://code.visualstudio.com/>.
2. Install the **WSL** extension (extension ID `ms-vscode-remote.remote-wsl`).

From a WSL shell, `code` should now exist (proxied to Windows VS Code, opens the current dir in Remote-WSL mode).

## Step 17 — Migrate the project

```bash
mkdir -p ~/projects
rsync -av --progress \
  --exclude=node_modules \
  --exclude=tmp \
  --exclude='.git/index.lock' \
  /mnt/c/Users/brand/Documents/projects/silkwave/ ~/projects/silkwave/

cd ~/projects/silkwave
git status   # should still show branch SW-ISSUE-2 with prior WIP
```

(`rsync` over `cp -r` because it preserves permissions, supports excludes, and shows progress for large trees.)

## Step 18 — Pin toolchain versions with `mise`

```bash
cd ~/projects/silkwave
cat > .tool-versions <<'EOF'
nodejs 22.12.0
golang 1.25.0
EOF

mise install
mise current   # confirms node 22 + go 1.25 active in this directory
```

## Step 19 — Enable corepack + pin pnpm in frontend

```bash
corepack enable

cd ~/projects/silkwave/frontend
# Add "packageManager" to package.json (replace 9.14.4 with the version from
# the existing pnpm-lock.yaml's lockfileVersion if needed — major-version match
# is what matters)
npm pkg set packageManager="pnpm@9.14.4"

pnpm install   # corepack auto-installs the pinned pnpm version
```

## Step 20 — Backend dev tools

```bash
cd ~/projects/silkwave/backend
go install github.com/air-verse/air@latest
go mod download

echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
exec zsh
```

## Step 21 — Restore secrets

For each `.env` file, pull values from your password manager. Generate a fresh `ARANGO_ROOT_PASSWORD` (it's only used by your local dev DB, no production blast).

- `~/projects/silkwave/backend/.env` — `ARANGO_ROOT_PASSWORD`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, etc.
- `~/projects/silkwave/.sandcastle/.env` — `ANTHROPIC_API_KEY` and any Sandcastle-specific config.

Verify both are gitignored (they should be already).

## Step 22 — Smoke test

In four terminal panes (VS Code Remote-WSL → split):

```bash
# Pane 1: ArangoDB
cd ~/projects/silkwave/backend
docker compose up -d
docker compose logs -f arangodb   # wait for "is ready for business"
# Browser on Windows → http://localhost:8529 → log in as root

# Pane 2: Backend
cd ~/projects/silkwave/backend
air

# Pane 3: Frontend
cd ~/projects/silkwave/frontend
pnpm dev
# Browser on Windows → http://localhost:5173

# Pane 4: Sandcastle (optional, only when running agents)
cd ~/projects/silkwave
npm install
npm run sandcastle
```

If all four work, setup is complete.

---

## Daily workflow

- **Open**: launch Windows Terminal → Ubuntu profile, or `wsl` from any prompt. From there: `code ~/projects/silkwave`. VS Code opens in Remote-WSL mode. Split the integrated terminal into 3 panes for ArangoDB / `air` / `pnpm dev`.
- **End of day**: `Ctrl-C` `air` and `pnpm dev`. Optionally `cd backend && docker compose down` if you want ArangoDB's RAM back. (`restart: unless-stopped` keeps it alive across daemon restarts otherwise.)
- **Sync**: `git push` over the SSH key. `gh pr create` for PRs.
- **Agents**: `npm run sandcastle` from project root.

## Troubleshooting

- **`docker: command not found`** — `sudo systemctl status docker`. If "inactive": `sudo systemctl start docker`. If systemd itself isn't running, re-check `/etc/wsl.conf` and `wsl --shutdown` from PowerShell.
- **`Permission denied (publickey)` pushing to GitHub** — confirm `cat ~/.ssh/id_ed25519.pub` matches what's at github.com/settings/keys. If you regenerated the key, the old one needs to be removed from GitHub.
- **Vite HMR misses changes** — confirm you're editing files at `~/projects/silkwave/...`, not `/mnt/c/...`. `pwd` and `git rev-parse --show-toplevel` both. Code on `/mnt/c` requires polling watchers, which are unreliable.
- **WSL eating all RAM** — confirm `~/.wslconfig` is at `C:\Users\brand\.wslconfig` (not inside WSL) and that `wsl --shutdown` was run after creating it. Verify with `wsl -d Ubuntu-24.04 -- free -h`.
- **`docker run hello-world` "permission denied"** — your user isn't in the `docker` group yet. Confirm `groups` lists `docker`. If not, `wsl --shutdown` from PowerShell and reopen.
- **Claude Code or Sandcastle "can't reach docker daemon"** — same root cause as above; user is outside `docker` group, or `dockerd` isn't started.

## Rebuilding from scratch

If the Ubuntu distro corrupts or you want to reset:

```powershell
wsl --shutdown
wsl --unregister Ubuntu-24.04
```

Then re-run this doc from Step 4. Total time: ~30 min, assuming password-manager values are accessible.

## Reminders

- **2026-11-06** — decide whether to delete `C:\Users\brand\backups\wsl-archlinux-2026-05-06\`. If the old key has not been needed by any service in 6 months, delete it.
