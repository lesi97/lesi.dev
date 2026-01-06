/**
 * Converts a time string (HH:MM:SS) into total seconds.
 *
 * @param {string} time - The time string in the format "HH:MM:SS".
 * @returns {number} - The total number of seconds.
 * @throws {Error} - Throws an error if the input format is invalid.
 */

export function timeToSeconds(time: string): number {
    const parts = time.split(':');

    if (parts.length !== 3 || parts.some((part) => isNaN(Number(part)))) {
        throw new Error('Invalid time format. Expected "HH:MM:SS".');
    }

    return Number(parts[0]) * 3600 + Number(parts[1]) * 60 + Number(parts[2]);
}
