export function replaceVideoExtensionWithMp3(name: string): string {
    return name.replace(/\.(mp4|m4v|avi|mov|qt|wmv|mkv)$/i, '.mp3');
}
