export function loadVideoToMp3() {
    return import('./videoToMp3').then((module) => {
        return { default: module.VideoToMp3 };
    });
}
