import { minify as minifyJS } from 'terser';
import xmlFormatter from 'xml-formatter';
import prettier from 'prettier/standalone';
import cssbeautify from 'cssbeautify';
import { css as beautifyCss } from 'js-beautify';
import type { Plugin } from 'prettier';
import * as parserBabel from 'prettier/parser-babel';

export function useMinifier() {
    function minify(code: string) {
        const type = determineCodeType(code);
        switch (type) {
            case 'json':
                return minifyJson(code);
            case 'xml':
                return minifyXml(code);
            case 'js':
                return minifyJs(code);
            case 'css':
                return minifyCss(code);
            case 'unknown':
                return;
        }
    }

    function unminify(code: string) {
        const type = determineCodeType(code);
        switch (type) {
            case 'json':
                return unminifyJson(code);
            case 'xml':
                return unminifyXml(code);
            case 'js':
                return unminifyJs(code);
            case 'css':
                return unminifyCss(code);
            case 'unknown':
                return;
        }
    }

    function determineCodeType(code: string) {
        const trimmedCode = code?.trim();

        if (trimmedCode.startsWith('{') && trimmedCode.endsWith('}')) {
            try {
                JSON.parse(trimmedCode);
                return 'json';
            } catch (e) {
                // Not valid JSON
            }
        }

        if (trimmedCode.startsWith('<') && trimmedCode.endsWith('>')) {
            return 'xml';
        }

        if (/^[\s\S]*[a-zA-Z0-9\s.#:_-]+\s*{[\s\S]*}$/m.test(trimmedCode)) {
            return 'css';
        }

        if (trimmedCode.includes('{') && trimmedCode.includes('}')) {
            return 'js';
        }

        return 'unknown';
    }

    function minifyJson(code: string) {
        return JSON.stringify(JSON.parse(code));
    }

    function unminifyJson(code: string) {
        return JSON.stringify(JSON.parse(code), null, 2);
    }

    function minifyXml(code: string) {
        return code.replace(/>\s+</g, '><').replace(/>\s+/g, '>').replace(/\s+</g, '<').trim();
    }

    function unminifyXml(code: string) {
        return xmlFormatter(code, { collapseContent: true, indentation: '\t', lineSeparator: '\n' });
    }

    async function minifyJs(code: string) {
        const jsResult = await minifyJS(code);
        return jsResult.code;
    }

    function unminifyJs(code: string) {
        return prettier.format(code, { parser: 'babel', plugins: [parserBabel as unknown as Plugin<any>] });
    }

    function minifyCss(code: string) {
        return (cssbeautify as unknown as (code: string, opts?: any) => string)(code, {
            indent: '',
            autosemicolon: true,
        }).replace(/\n\s*/g, '');
    }

    function unminifyCss(code: string) {
        return beautifyCss(code, { indent_size: 4 });
    }

    return { minify, unminify };
}
