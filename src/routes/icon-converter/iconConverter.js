const MaxSize = 256; // 1 << 8;
const FileHeaderSize = 6;
const ImageHeaderSize = 16;
const IcoMime = "image/x-icon";

export class PngIcoConverter {
    async convertToBlobAsync(inputs, mime = IcoMime) {
        const arr = await this.convertAsync(inputs);
        return new Blob([arr], {
            type: mime,
        });
    }

    async convertAsync(inputs) {
        const inLen = inputs.length;

        const headersLen = FileHeaderSize + ImageHeaderSize * inLen;
        const totalLen = headersLen + await this.sumInputLen(inputs);
        const arr = new Uint8Array(totalLen);

        // File Header
        arr.set([0, 0, 1, 0, ...this.to2Bytes(inLen)], 0);

        // Image Headers & Data
        let imgPos = headersLen;
        for (let i = 0; i < inputs.length; i++) {
            const currPos = FileHeaderSize + ImageHeaderSize * i;
            const input = inputs[i];
            const pngBlob = await this.convertToPngBlob(input.png);
            const img = await this.loadImageAsync(pngBlob);
            const w = img.naturalWidth, h = img.naturalHeight;

            if (!input.ignoreSize && (w > MaxSize || h > MaxSize)) {
                throw new Error("INVALID_SIZE");
            }
            // Header
            arr.set([
                w > MaxSize ? 0 : w,
                h > MaxSize ? 0 : h,
                0,
                0,
                0, 0,
                ...(input.bpp ? this.to2Bytes(input.bpp) : [0, 0]),
                ...this.to4Bytes(pngBlob.size),
                ...this.to4Bytes(imgPos),
            ], currPos);
            // Image
            const buffer = await pngBlob.arrayBuffer();
            arr.set(new Uint8Array(buffer), imgPos);
            imgPos += pngBlob.size;
        }
        return arr;
    }

    loadImageAsync(blob) {
        return new Promise((resolve, reject) => {
            const img = new Image();
            img.onload = () => resolve(img);
            img.onerror = () => reject(new Error("INVALID_IMAGE"));
            img.src = URL.createObjectURL(blob);
        });
    }

    async convertToPngBlob(input) {
        const img = await this.loadImageAsync(this.toBlob(input, input.type));

        // Create a canvas to draw the image on
        const canvas = document.createElement("canvas");
        canvas.width = img.naturalWidth;
        canvas.height = img.naturalHeight;

        const ctx = canvas.getContext("2d");
        ctx.drawImage(img, 0, 0);

        // Convert the canvas to a PNG Blob
        return new Promise((resolve) => {
            canvas.toBlob((blob) => {
                resolve(blob);
            }, "image/png");
        });
    }

    toBlob(input, type = "image/png") {
        return input instanceof Blob ? input : new Blob([input], {
            type: type,
        });
    }

    to2Bytes(n) {
        return [n & 255, (n >> 8) & 255];
    }

    to4Bytes(n) {
        return [n & 255, (n >> 8) & 255, (n >> 16) & 255, (n >> 24) & 255];
    }

    async sumInputLen(inputs) {
        let total = 0;
        for (const input of inputs) {
            const pngBlob = await this.convertToPngBlob(input.png);
            total += pngBlob.size;
        }
        return total;
    }
}

export class ConvertApp {
    static currBlob = null;

    static onDownload(img) {
        if (!this.currBlob) return;

        const url = URL.createObjectURL(this.currBlob);
        const a = document.createElement("a");
        a.href = url;

        // Remove the original file extension and append .ico
        const originalName = img?.name || "favicon";
        const nameWithoutExtension = originalName.replace(/\.[^/.]+$/, "");
        const name = `${nameWithoutExtension}.ico`;

        a.download = name;
        a.click();

        URL.revokeObjectURL(url); // Clean up
    }

    static async convert(img) {
        const converter = new PngIcoConverter();
        const ignoreSize = true;
        const inputs = [{ png: img, ignoreSize }];

        try {
            this.currBlob = await converter.convertToBlobAsync(inputs);
            ConvertApp.onDownload(img);
        } catch (e) {
            console.error(e);
        }
    }
}
