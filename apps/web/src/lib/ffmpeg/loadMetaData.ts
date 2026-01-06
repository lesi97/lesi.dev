export function loadMetadata(video: HTMLVideoElement) {
    return new Promise((resolve, reject) => {
        video.onloadedmetadata = () => {
            console.log('Metadata loaded!');
            resolve(video.duration);
        };
        video.onerror = (e) => {
            console.error('Failed to load metadata:', e);
            reject(e);
        };
    });
}
