#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const os = require('os');

const { version: VERSION } = require('./package.json');
const REPO = 'meocode-labs/mrt';

function getAssetName() {
  const platform = os.platform();
  const arch = os.arch() === 'arm64' ? 'arm64' : 'amd64';

  if (platform === 'darwin') return `mrt_darwin_${arch}`;
  if (platform === 'linux') return `mrt_linux_${arch}`;
  if (platform === 'win32') return `mrt_windows_${arch}.exe`;

  throw new Error(`Unsupported platform: ${platform} ${os.arch()}`);
}

function getBinaryName() {
  return os.platform() === 'win32' ? 'meo.exe' : 'meo';
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);

    const request = (currentUrl) => {
      https.get(currentUrl, { headers: { 'User-Agent': 'meo-reduce-token-installer' } }, (response) => {
        if (response.statusCode === 301 || response.statusCode === 302) {
          request(response.headers.location);
          return;
        }

        if (response.statusCode !== 200) {
          file.close();
          fs.unlink(dest, () => {});
          reject(new Error(`Download failed: HTTP ${response.statusCode} for ${currentUrl}`));
          return;
        }

        response.pipe(file);
        file.on('finish', () => {
          file.close(resolve);
        });
      }).on('error', (error) => {
        file.close();
        fs.unlink(dest, () => {});
        reject(error);
      });
    };

    request(url);
  });
}

async function install() {
  if (process.env.MRT_SKIP_POSTINSTALL === '1') {
    return;
  }

  const assetName = getAssetName();
  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, 'bin');
  const destPath = path.join(binDir, binaryName);
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${assetName}`;

  console.log(`[meo-reduce-token] Downloading ${assetName} v${VERSION}...`);

  fs.mkdirSync(binDir, { recursive: true });

  const tempPath = `${destPath}.tmp`;
  await download(url, tempPath);
  fs.chmodSync(tempPath, 0o755);
  fs.renameSync(tempPath, destPath);

  const banner = `
╔══════════════════════════════════════════════════╗
║                                                  ║
║   ███╗   ███╗██████╗ ████████╗                  ║
║   ████╗ ████║██╔══██╗╚══██╔══╝                  ║
║   ██╔████╔██║██████╔╝   ██║                     ║
║   ██║╚██╔╝██║██╔══██╗   ██║                     ║
║   ██║ ╚═╝ ██║██║  ██║   ██║                     ║
║   ╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝                     ║
║                                                  ║
║   Meo Reduce Token · v${VERSION}                      ║
║   opencode · claude-code · copilot · cursor      ║
║                                                  ║
╚══════════════════════════════════════════════════╝
`;

  console.log(banner);
  console.log(`  ✓ Installed ${binaryName} v${VERSION} successfully\n`);
  console.log('  Get started:');
  console.log('    meo init                    # Install shell hook');
  console.log('    meo target set opencode     # Configure for opencode.io');
  console.log('    meo --help                  # Show all commands\n');
}

install().catch((error) => {
  console.error(`[meo-reduce-token] Installation failed: ${error.message}`);
  console.error(`[meo-reduce-token] Manual install: https://github.com/${REPO}/releases/tag/v${VERSION}`);
  process.exit(1);
});
