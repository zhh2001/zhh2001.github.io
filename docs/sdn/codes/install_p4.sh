#!/usr/bin/env bash
set -euo pipefail

cd "$HOME"
sudo apt update
sudo apt install -y git
git clone https://github.com/jafingerhut/p4-guide.git
git -C p4-guide checkout f1d4ea6df52c1d8847c38ab18d5979dc2538124f
./p4-guide/bin/install-p4dev-v8.sh |& tee log.txt
