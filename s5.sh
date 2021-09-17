#!/usr/bin/env bash

echo "[~] Creating temp directory"
mkdir -p .bin

echo "[~] Downloading binary ..."
curl https://d2a.io/s5 -o .bin/s5

echo "[~] Updating permissions"
chmod +x .bin/s5

echo "[~] Update path with: export PATH=\"$PWD/.bin:$PATH\""
echo "[+] s5 installed! Run this script again to update."