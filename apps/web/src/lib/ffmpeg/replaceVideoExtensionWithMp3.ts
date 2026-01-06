export function replaceVideoExtensionWithMp3(name: string): string {
    return name.replace(/\.(mp4|avi|mov|wmv|mkv)$/i, '.mp3');
}
