export function getPlatformEnum(platform: string | undefined): number {
    switch (platform) {
        case 'xb':
            return 1;
        case 'ps':
            return 2;
        case 'pc':
            return 3;
        case 'bnet':
            return 4;
        case 'st':
            return 5;
        case 'demon':
            return 10;
        default:
            return -1;
    }
}
