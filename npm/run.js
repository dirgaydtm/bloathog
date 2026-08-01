#!/usr/bin/env node
const { spawnSync } = require("node:child_process");
const path = require("node:path");
const os = require("node:os");
const fs = require("node:fs");

// Determine the correct binary name based on the OS
const binaryName = os.platform() === "win32" ? "bloathog.exe" : "bloathog";
const binaryPath = path.join(__dirname, "..", "bin", binaryName);

// Ensure the binary was successfully downloaded during postinstall
if (!fs.existsSync(binaryPath)) {
	console.error("bloathog binary not found.");
	console.error("Please reinstall the package: npm install -g bloathog");
	process.exit(1);
}

// Forward all command-line arguments to the Go binary
const result = spawnSync(binaryPath, process.argv.slice(2), {
	stdio: "inherit",
});

process.exit(result.status || 0);
