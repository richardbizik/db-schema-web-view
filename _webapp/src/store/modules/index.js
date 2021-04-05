const requireModule = require.context('.', true, /\.js$/);
const index = {};

requireModule.keys().forEach(fileName => {
    if (fileName === './index.js') return;

    // Replace ./ and .js
    const path = fileName.replace(/(\.\/|\.js)/g, '');
    const [moduleName, imported] = path.split('/');

    if (!index[moduleName]) {
        index[moduleName] = {
            namespaced: true
        }
    }

    index[moduleName][imported] = requireModule(fileName).default
});

export default index