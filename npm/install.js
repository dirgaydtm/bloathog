const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const https = require("node:https");
const { execSync } = require("node:child_process");

const version = require("../package.json").version;

// Map Node.js platform/arch to GoReleaser format
const osMap = { Linux: "Linux", Darwin: "Darwin", Windows_NT: "Windows" };
const archMap = { x64: "x86_64", arm64: "arm64", ia32: "i386" };

const osName = osMap[os.type()];
const archName = archMap[os.arch()];

if (!osName || !archName) {
	console.error(`Unsupported platform: ${os.type()} ${os.arch()}`);
	process.exit(1);
}

// Construct the download URL based on the OS and architecture
const ext = osName === "Windows" ? "zip" : "tar.gz";
const filename = `bloathog_${osName}_${archName}.${ext}`;
const url = `https://github.com/dirgaydtm/bloathog/releases/download/v${version}/${filename}`;

const binDir = path.join(__dirname, "..", "bin");
const archivePath = path.join(__dirname, filename);

// Ensure the bin directory exists
if (!fs.existsSync(binDir)) fs.mkdirSync(binDir, { recursive: true });

// Downloads a file from a URL, following redirects automatically.
function download(url, dest) {
	return new Promise((resolve, reject) => {
		https
			.get(url, (res) => {
				// Handle redirects (e.g., GitHub releases redirecting to S3)
				if (
					res.statusCode >= 300 &&
					res.statusCode < 400 &&
					res.headers.location
				) {
					return download(res.headers.location, dest)
						.then(resolve)
						.catch(reject);
				}
				if (res.statusCode !== 200)
					return reject(new Error(`Status: ${res.statusCode}`));

				const file = fs.createWriteStream(dest);
				res.pipe(file);
				file.on("finish", () => {
					file.close();
					resolve();
				});
				file.on("error", (err) => fs.unlink(dest, () => reject(err)));
			})
			.on("error", reject);
	});
}

console.log(`Downloading bloathog v${version} for ${osName} ${archName}...`);

// Start the download and extraction process
download(url, archivePath)
	.then(() => {
		console.log("Extracting binary...");
		const tarFlag = ext === "zip" ? "-xf" : "-xzf";
		execSync(`tar ${tarFlag} "${archivePath}" -C "${binDir}"`, {
			stdio: "inherit",
		});

		// Clean up the downloaded archive
		fs.unlinkSync(archivePath);
		console.log("bloathog installed successfully.");
	})
	.catch((err) => {
		console.error("Installation failed:", err.message);
		process.exit(1);
	});
