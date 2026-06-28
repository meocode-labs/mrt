#!/usr/bin/env node
/**
 * MRT (Meo Reduce Token) NPM postinstall script.
 * Downloads the pre-compiled Go binary for the current platform from
 * the GitHub release and overwrites the `bin/mrt` shim with the real
 * binary. The npm link to `./bin/mrt` will then point at the real one.
 */
const fs = require('fs');
const path = require('path');
const https = require('https');
const http = require('http');
const os = require('os');

const { version: VERSION } = require('./package.json');
const REPO = process.env.MRT_REPO || 'meocode-labs/mrt';
const BASE_URL = (process.env.MRT_BASE_URL || 'https://github.com').replace(/\/$/, '');

function getAssetName() {
  const platform = os.platform();
  const arch = os.arch() === 'arm64' ? 'arm64' : 'amd64';

  if (platform === 'darwin') return `mrt_darwin_${arch}`;
  if (platform === 'linux')  return `mrt_linux_${arch}`;
  if (platform === 'win32')  return `mrt_windows_${arch}.exe`;

  throw new Error(`Unsupported platform: ${platform} ${os.arch()}`);
}

function download(url, dest, redirectsLeft = 5) {
  return new Promise((resolve, reject) => {
    const request = (currentUrl) => {
      const proto = currentUrl.startsWith('https:') ? https : http;
      proto.get(currentUrl, { headers: { 'User-Agent': 'mrt-npm-installer' } }, (response) => {
        if (response.statusCode === 301 || response.statusCode === 302) {
          if (redirectsLeft === 0) return reject(new Error('Too many redirects'));
          response.resume();
          return request(response.headers.location);
        }
        if (response.statusCode !== 200) {
          fs.unlink(dest, () => {});
          return reject(new Error(`HTTP ${response.statusCode} for ${currentUrl}`));
        }
        const file = fs.createWriteStream(dest);
        response.pipe(file);
        file.on('finish', () => file.close(resolve));
        file.on('error', (e) => { fs.unlink(dest, () => {}); reject(e); });
      }).on('error', (error) => {
        fs.unlink(dest, () => {});
        reject(error);
      });
    };
    request(url);
  });
}

async function install() {
  if (process.env.MRT_SKIP_POSTINSTALL === '1') return;

  const assetName = getAssetName();
  const binaryName = os.platform() === 'win32' ? 'mrt.exe' : 'mrt';
  const binDir = path.join(__dirname, 'bin');
  const destPath = path.join(binDir, binaryName);
  const oldBinaryPath = path.join(binDir, 'meo' + (os.platform() === 'win32' ? '.exe' : ''));
  const url = `${BASE_URL}/${REPO}/releases/download/v${VERSION}/${assetName}`;

  console.log(`[mrt] Downloading ${assetName} v${VERSION}...`);

  fs.mkdirSync(binDir, { recursive: true });
  const tempPath = `${destPath}.tmp`;
  try {
    await download(url, tempPath);
  } catch (err) {
    throw new Error(`Could not download from ${url}: ${err.message}`);
  }
  fs.chmodSync(tempPath, 0o755);
  fs.renameSync(tempPath, destPath);

  // Remove old shim / leftover from v1.2.0 (`bin/meo`)
  if (fs.existsSync(oldBinaryPath)) {
    try {
      fs.unlinkSync(oldBinaryPath);
      console.log(`[mrt] ✓ Removed leftover v1.2.0 \`bin/meo\` from this package`);
    } catch (e) {
      console.warn(`[mrt] ! Could not remove ${oldBinaryPath}: ${e.message}`);
    }
  }

  // Notice about globally installed old `meo` command (v1.2.0).
  // We can't easily find where the previous `meo-reduce-token` package
  // lived, so we just warn the user.
  console.log(`
╔══════════════════════════════════════════════════╗
║                                                  ║
║   ███╗   ███╗██████╗ ████████╗                  ║
║   ████╗ ████║██╔══██╗╚══██╔══╝                  ║
║   ██╔████╔██║██████╔╝   ██║                     ║
║   ██║╚██╔╝██║██╔══██╗   ██║                     ║
║   ██║ ╚═╝ ██║██║  ██║   ██║                     ║
║   ╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝                     ║
║                                                  ║
║   Meo Reduce Token · v${VERSION.padEnd(28)}║
║   opencode · claude-code · copilot · cursor      ║
║                                                  ║
╚══════════════════════════════════════════════════╝

  ✓ Installed ${binaryName} v${VERSION} successfully

  Get started:
    mrt init                    # Install shell hook
    mrt target set opencode     # Configure for opencode.io
    mrt --help                  # Show all commands
`);

  // Migration notice: v1.2.0 had the command as `meo`.
  console.log(`  ⚠ v1.3.0 renamed the command from \`meo\` to \`mrt\`.`);
  console.log(`    If you previously installed \`meo-reduce-token\`, remove the old binary:`);
  if (process.platform === 'win32') {
    console.log(`      npm uninstall -g meo-reduce-token`);
  } else {
    console.log(`      npm uninstall -g meo-reduce-token 2>/dev/null`);
    console.log(`      rm -f "$(npm prefix -g)/bin/meo"`);
  }
  console.log(`    Or set MRT_AUTO_MIGRATE=1 before install to do it automatically.\n`);
}

install().catch((error) => {
  console.error(`[mrt] Installation failed: ${error.message}`);
  console.error(`[mrt] Manual install: https://github.com/${REPO}/releases/tag/v${VERSION}`);
  process.exit(1);
});