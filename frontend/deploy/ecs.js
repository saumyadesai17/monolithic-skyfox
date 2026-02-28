const fs = require('fs');

const templateName = process.argv[2];

const template = require('./'+templateName+'.js');

let data = JSON.stringify(template.config, null, 2);

let outputFileName = template.name + '.json'

fs.writeFileSync(outputFileName, data);