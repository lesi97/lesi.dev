import { KeyboardEvent } from 'react';

/**
 * Prevent character input
 *
 * @export
 * @param {KeyboardEvent<HTMLInputElement>} e
 * @param {(string[] | string)} charsList
 * @example
 * ```js
 * const chars = 'e+-';
 * const charsArr = ['e', '+', '-'];
 * <input onChange={(e) => preventCharsInInput(e.target.value, chars)} />
 * <input onChange={(e) => preventCharsInInput(e.target.value, charsArr)} />
 */
export function preventCharsInInput(e: KeyboardEvent<HTMLInputElement>, charsList: string[] | string) {
    let chars: string[] = [];
    if (typeof charsList === 'string') {
        charsList.split('').forEach((char) => {
            chars.push(char);
        });
    }
    if (typeof charsList === 'object') {
        chars = charsList;
    }

    if (chars.includes(e.key.toLocaleLowerCase())) {
        e.preventDefault();
    }
}
