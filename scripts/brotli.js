const { readFileSync, writeFileSync } = require('node:fs');
const { brotliDecompressSync } = require('node:zlib');

function decodeBr(inputPath) {
    const outPath = inputPath.endsWith('.br')
        ? inputPath.slice(0, -3)
        : `${inputPath}.out`;
    const data = readFileSync(inputPath);
    const decoded = brotliDecompressSync(data);
    writeFileSync(outPath, decoded);
    console.log(`Wrote ${outPath}`);
}

decodeBr('./apps/web/public/aim-trainer/Build/aim-trainer.data.br');
decodeBr('./apps/web/public/aim-trainer/Build/aim-trainer.framework.js.br');
decodeBr('./apps/web/public/aim-trainer/Build/aim-trainer.wasm.br');
