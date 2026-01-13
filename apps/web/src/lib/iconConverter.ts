interface IconInput {
    png: Blob;
    ignoreSize: boolean;
    bpp?: number;
}

const config = {
    MaxSize: 256,
    FileHeaderSize: 6,
    ImageHeaderSize: 16,
    IcoMime: 'image/x-icon',
};

class PngIcoConverter {
    async convertToBlobAsync(inputs: IconInput[], mime = config.IcoMime): Promise<Blob> {
        const arr = await this.convertAsync(inputs);
        return new Blob([arr as BlobPart], {
            type: mime,
        });
    }

    async convertAsync(inputs: IconInput[]): Promise<Uint8Array> {
        const inLen = inputs.length;

        const headersLen = config.FileHeaderSize + config.ImageHeaderSize * inLen;
        const totalLen = headersLen + (await this.sumInputLen(inputs));
        const arr = new Uint8Array(totalLen);

        // File Header
        arr.set([0, 0, 1, 0, ...this.to2Bytes(inLen)], 0);

        // Image Headers & Data
        let imgPos = headersLen;
        for (let i = 0; i < inputs.length; i++) {
            const currPos = config.FileHeaderSize + config.ImageHeaderSize * i;
            const input = inputs[i];
            const pngBlob = await this.convertToPngBlob(input.png);
            const img = await this.loadImageAsync(pngBlob);
            const w = img.naturalWidth;
            const h = img.naturalHeight;

            if (!input.ignoreSize && (w > config.MaxSize || h > config.MaxSize)) {
                throw new Error('INVALID_SIZE');
            }
            // Header
            arr.set(
                [
                    w > config.MaxSize ? 0 : w,
                    h > config.MaxSize ? 0 : h,
                    0,
                    0,
                    0,
                    0,
                    ...(input.bpp ? this.to2Bytes(input.bpp) : [0, 0]),
                    ...this.to4Bytes(pngBlob.size),
                    ...this.to4Bytes(imgPos),
                ],
                currPos
            );
            // Image data
            const buffer = await pngBlob.arrayBuffer();
            arr.set(new Uint8Array(buffer), imgPos);
            imgPos += pngBlob.size;
        }
        return arr;
    }

    loadImageAsync(blob: Blob): Promise<HTMLImageElement> {
        return new Promise((resolve, reject) => {
            const img = new Image();
            img.onload = () => resolve(img);
            img.onerror = () => reject(new Error('INVALID_IMAGE'));
            img.src = URL.createObjectURL(blob);
        });
    }

    async convertToPngBlob(input: Blob): Promise<Blob> {
        // Ensure we have a proper Blob to load the image with.
        const blobToUse = this.toBlob(input, input.type);
        const img = await this.loadImageAsync(blobToUse);

        // Create a canvas to draw the image on
        const canvas = document.createElement('canvas');
        canvas.width = img.naturalWidth;
        canvas.height = img.naturalHeight;

        const ctx = canvas.getContext('2d');
        if (!ctx) {
            throw new Error('Could not get canvas context');
        }
        ctx.drawImage(img, 0, 0);

        // Convert the canvas to a PNG Blob
        return new Promise<Blob>((resolve, reject) => {
            canvas.toBlob((blob) => {
                if (blob) {
                    resolve(blob);
                } else {
                    reject(new Error('Canvas is empty or toBlob conversion failed'));
                }
            }, 'image/png');
        });
    }

    toBlob(input: BlobPart, type = 'image/png'): Blob {
        return input instanceof Blob
            ? input
            : new Blob([input], {
                  type: type,
              });
    }

    to2Bytes(n: number): number[] {
        return [n & 255, (n >> 8) & 255];
    }

    to4Bytes(n: number): number[] {
        return [n & 255, (n >> 8) & 255, (n >> 16) & 255, (n >> 24) & 255];
    }

    async sumInputLen(inputs: IconInput[]): Promise<number> {
        let total = 0;
        for (const input of inputs) {
            const pngBlob = await this.convertToPngBlob(input.png);
            total += pngBlob.size;
        }
        return total;
    }
}

export class ConvertApp {
    static currBlob: Blob | null = null;

    static onDownload(img: File): void {
        if (!this.currBlob) return;

        const url = URL.createObjectURL(this.currBlob);
        const a = document.createElement('a');
        a.href = url;

        const originalName = img.name || 'favicon';
        const nameWithoutExtension = originalName.replace(/\.[^/.]+$/, '');
        const name = `${nameWithoutExtension}.ico`;

        a.download = name;
        a.click();

        URL.revokeObjectURL(url);
        a.remove();
    }

    static async convert(img: File): Promise<void> {
        const converter = new PngIcoConverter();
        const ignoreSize = true;
        const inputs: IconInput[] = [{ png: img, ignoreSize }];

        try {
            this.currBlob = await converter.convertToBlobAsync(inputs);
            ConvertApp.onDownload(img);
        } catch (e) {
            console.error(e);
        }
    }
}
