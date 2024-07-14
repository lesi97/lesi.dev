import { minify as minifyJS } from 'terser';
import xmlFormatter from 'xml-formatter';
import prettier from 'prettier/standalone';
import parserBabel from 'prettier/parser-babel';
import cssbeautify from 'cssbeautify';
import { css as beautifyCss } from 'js-beautify';

export function minify(code) {
    const type = determineCodeType(code);
    switch (type) {
        case "json":
            return minifyJson(code);
        case "xml":
            return minifyXml(code);
        case "js":
            return minifyJs(code);
        case "css":
            return minifyCss(code);
        case "unknown":
            return;
    }
}

export function unminify(code) {
    const type = determineCodeType(code);
    switch (type) {
        case "json":
            return unminifyJson(code);
        case "xml":
            return unminifyXml(code);
        case "js":
            return unminifyJs(code);
        case "css":
            return unminifyCss(code);
        case "unknown":
            return;
    }
}

function determineCodeType(code) {
    const trimmedCode = code?.trim();

    if (trimmedCode.startsWith("{") && trimmedCode.endsWith("}")) {
        try {
            JSON.parse(trimmedCode);
            return "json";
        } catch (e) {
            // Not valid JSON
        }
    }

    if (trimmedCode.startsWith("<") && trimmedCode.endsWith(">")) {
        return "xml";
    }

    if (/^[\s\S]*[a-zA-Z0-9\s.#:_-]+\s*{[\s\S]*}$/m.test(trimmedCode)) {
        return "css";
    }

    if (trimmedCode.includes("{") && trimmedCode.includes("}")) {
        return "js";
    }

    return "unknown";
}


function minifyJson(code) {
    return JSON.stringify(JSON.parse(code));
}

function unminifyJson(code) {
    return JSON.stringify(JSON.parse(code), null, 2);
}

function minifyXml(code) {
    return code.replace(/\>\s+\</g, '><').trim();
}

function unminifyXml(code) {
    return xmlFormatter(code);
}

async function minifyJs(code) {
    const jsResult = await minifyJS(code);
    return jsResult.code;
}

function unminifyJs(code) {
    return prettier.format(code, { parser: 'babel', plugins: [parserBabel] });
}

function minifyCss(code) {
    return cssbeautify(code, { indent: '', autosemicolon: true }).replace(/\n\s*/g, '');
}

function unminifyCss(code) {
    return beautifyCss(code, { indent_size: 4 });
}