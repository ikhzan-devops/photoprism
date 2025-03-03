const glob = require('glob');
const poFiles = glob.sync('src/locales/*.po');

const languageCodes = poFiles.map(filePath => {
  const fileName = filePath.split('/').pop();
  return fileName.replace('.po', '');
});

module.exports = {
  output: {
    path: "./src/locales",
    jsonPath: "",
    locales: languageCodes,
    splitJson: true,
  },
};
