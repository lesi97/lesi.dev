export function loadVideoCropper() {
    return import('./videoCropper').then((module) => {
        return { default: module.VideoCropper };
    });
}
