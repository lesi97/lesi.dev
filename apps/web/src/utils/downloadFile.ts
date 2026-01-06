/**
 * Creates a temporary link to download a file
 *
 * @export
 * @param {string} name - Name of the file
 * @param {string} url - Url of the file in memory
 * @throws {Error} - Throws an error if the input format is invalid.
 */

export function downloadFile(name: string, url: string) {
    if (!name || !url) {
        throw new Error('File name or URL expected.');
    }
    const downloadLink = document.createElement('a');
    downloadLink.download = name;
    downloadLink.href = url;
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
}
