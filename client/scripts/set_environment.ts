const { writeFile, existsSync, mkdirSync } = require('fs');

require('dotenv').config({ path: '../.env', override: true });

const environmentFileContent = `
export const environment = {
    API_HOST: '${process.env.API_HOST}',
    API_PORT: '${process.env.API_PORT}'
}
`;

const envDirectory = './src/environments/';
const envFile = 'environment.ts';

if (!existsSync(envDirectory)) {
    mkdirSync(envDirectory);
}

writeFile(`${envDirectory}${envFile}`, environmentFileContent, () => {});
