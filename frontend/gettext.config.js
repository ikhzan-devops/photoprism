const glob = require("glob");
const fs = require("fs");
const path = require("path");

// Find all .po files in the frontend/src/locales
const poFiles = glob.sync("src/locales/*.po");
if (poFiles.length === 0) {
  console.error("No .po files found in src/locales");
  process.exit(1);
}

// Find output folder or create the new one in assets/static/locales
const outputDir = path.resolve(__dirname, "../assets/static/locales");
if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}

// Copy .po files from frontend/src/locales to assets/static/locales
poFiles.forEach((filePath) => {
  const fileName = path.basename(filePath);
  const destinationPath = path.join(outputDir, fileName);
  fs.copyFileSync(filePath, destinationPath);
});

// Find all languages codes from .po files (cut file names without .po)
const languageCodes = poFiles.map((filePath) => {
  const fileName = path.basename(filePath);
  return fileName.replace(".po", "");
});

// Transform files from .po to .json in the assets/static/locales
module.exports = {
  input: {
    path: path.resolve(__dirname, "../assets/static/locales"),
    include: ["**/*.po"],
  },
  output: {
    path: path.resolve(__dirname, "../assets/static/locales"),
    jsonPath: "",
    locales: languageCodes,
    splitJson: true,
  },
};
