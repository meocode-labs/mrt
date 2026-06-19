#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const http = require('http');
const os = require('os');

const VERSION = process.env.npm_package_version || '1.0.0';
const REPO = 'meocodelabs/mrt';
const API_URL = `https://api.github.com/repos/${REPO}/releases/latest`;

const PLATFORMS = {
  darwin: { arch: { x64: 'apple-darwin-x64', arm64: 'apple-darwin-arm64' } },
  linux: { arch: { x64: 'unknown-linux-musl-x64', arm64: 'unknown-linux-musl-arm64' } },
  win32: { arch: { x64: 'pc-windows-x64', arm64: 'pc-windows-arm64' } }
};

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();
  
  const platformKey = PLATFORMS[platform];
  if (!platformKey) {
    console.error(`Unsupported platform: ${platform}`);
    process.exit(1);
  }
  
  let archKey = arch === 'arm64' ? 'arm64' : 'x64';
  if (!platformKey.arch[archKey]) {
    archKey = 'x64';
  }
  
  return {
    platform,
    arch: archKey,
    binaryName: platform === 'win32' ? 'meo.exe' : 'meo',
    compressedArch: platformKey.arch[archKey]
  };
}

function getBinaryUrl(platform) {
  return `https://github.com/${REPO}/releases/download/v${VERSION}/mrt_v${VERSION}_${platform.compressedArch}.tar.gz`;
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    const protocol = url.startsWith('https') ? https : http;
    
    const request = protocol.get(url, { headers: { 'User-Agent': 'meo-npm-installer' } }, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        downloadFile(response.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      
      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode}`));
        return;
      }
      
      response.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve();
      });
    });
    
    request.on('error', reject);
  });
}

function extractTarGz(tarPath, destDir) {
  const { spawn } = require('child_process');
  
  return new Promise((resolve, reject) => {
    const tar = spawn('tar', ['-xzf', tarPath, '-C', destDir]);
    tar.on('close', (code) => {
      if (code === 0) resolve();
      else reject(new Error(`tar exited with code ${code}`));
    });
    tar.on('error', reject);
  });
}

async function install() {
  console.log('╔══════════════════════════════════════════════════╗');
  console.log('║  ✦ Meo Reduce Token (MRT) Installer ✦          ║');
  console.log('╚══════════════════════════════════════════════════╝\n');
  
  console.log(`Platform: ${os.platform()} ${os.arch()}`);
  console.log(`Version: ${VERSION}\n`);
  
  const platform = getPlatform();
  const url = getBinaryUrl(platform);
  
  console.log(`Downloading: ${url}\n`);
  
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'mrt-'));
  const tarPath = path.join(tempDir, 'mrt.tar.gz');
  const binDir = path.join(tempDir, 'bin');
  
  try {
    await downloadFile(url, tarPath);
    fs.mkdirSync(binDir, { recursive: true });
    
    console.log('Extracting...');
    await extractTarGz(tarPath, binDir);
    
    const extractedFile = path.join(binDir, platform.binaryName);
    const targetPath = path.join(process.env.PREFIX || '/usr/local', 'bin', platform.binaryName);
    const targetDir = path.dirname(targetPath);
    
    if (!fs.existsSync(targetDir)) {
      fs.mkdirSync(targetDir, { recursive: true });
    }
    
    fs.copyFileSync(extractedFile, targetPath);
    fs.chmodSync(targetPath, 0o755);
    
    console.log(`\n✓ Installed to: ${targetPath}`);
    console.log('\nRun "meo --help" to get started!\n');
    
  } catch (error) {
    console.error(`\n✗ Installation failed: ${error.message}`);
    console.error('\nYou can try installing manually:');
    console.error(`  curl -fsSL ${url} | tar -xz`);
    console.error(`  mv meo /usr/local/bin/\n`);
    process.exit(1);
  } finally {
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
}

install();
