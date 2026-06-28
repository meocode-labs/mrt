#!/usr/bin/env node
/**
 * Lazy fallback runner — only used if the `mrt` shim is invoked before
 * postinstall has had a chance to drop the real binary. It re-runs the
 * installer to fetch the binary, then exec's it.
 */
const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const ROOT = path.resolve(__dirname, '..');
const realBinaryGuess = path.join(ROOT, 'bin', process.platform === 'win32' ? 'mrt.exe' : 'mrt');

function tryExec() {
  if (fs.existsSync(realBinaryGuess)) {
    const p = spawn(realBinaryGuess, process.argv.slice(2), { stdio: 'inherit' });
    p.on('close', (code) => process.exit(code ?? 0));
    return true;
  }
  return false;
}

if (!tryExec()) {
  console.error('mrt: binary not installed yet. Run `npm install -g mrt` again,');
  console.error('or: curl -fsSL https://raw.githubusercontent.com/meocode-labs/mrt/main/install.sh | bash');
  process.exit(1);
}