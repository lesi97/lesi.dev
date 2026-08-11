export const UUID_VERSIONS = ['1', '3', '4', '5', '6', '7', 'nil'] as const;

export type UuidVersion = (typeof UUID_VERSIONS)[number];

export const UUID_NAMESPACE_PRESETS = {
    dns: '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
    url: '6ba7b811-9dad-11d1-80b4-00c04fd430c8',
    oid: '6ba7b812-9dad-11d1-80b4-00c04fd430c8',
    x500: '6ba7b814-9dad-11d1-80b4-00c04fd430c8',
} as const;

export type UuidNamespacePreset = keyof typeof UUID_NAMESPACE_PRESETS;

const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
const NIL_UUID = '00000000-0000-0000-0000-000000000000';
const GREGORIAN_OFFSET_100NS = 122192928000000000n;
const MAX_BULK_UUIDS = 500;

let lastTimestampMilliseconds = 0;
let lastTimestampSubMilliseconds = 0;
let clockSequence: number | null = null;
let nodeId: Uint8Array | null = null;

export function isUuidVersion(value: string | undefined): value is UuidVersion {
    return UUID_VERSIONS.includes(value as UuidVersion);
}

export function isUuid(value: string) {
    return UUID_REGEX.test(value);
}

export async function generateUuids(
    version: UuidVersion,
    count: number,
    options: { namespace?: string; name?: string } = {}
) {
    const safeCount = Math.max(1, Math.min(MAX_BULK_UUIDS, Math.floor(count)));
    const uuids: string[] = [];

    for (let i = 0; i < safeCount; i++) {
        uuids.push(await generateUuid(version, options));
    }

    return uuids;
}

export async function generateUuid(
    version: UuidVersion,
    options: { namespace?: string; name?: string } = {}
): Promise<string> {
    switch (version) {
        case '1':
            return generateV1Uuid();
        case '3':
            return generateNameBasedUuid(3, options);
        case '4':
            return generateV4Uuid();
        case '5':
            return generateNameBasedUuid(5, options);
        case '6':
            return generateV6Uuid();
        case '7':
            return generateV7Uuid();
        case 'nil':
            return NIL_UUID;
        default:
            throw new Error(`Unsupported UUID version: ${version}`);
    }
}

function generateV1Uuid() {
    const timestamp = nextGregorianTimestamp();
    const bytes = new Uint8Array(16);
    const sequence = getClockSequence();

    writeBigEndian(bytes, 0, timestamp & 0xffffffffn, 4);
    writeBigEndian(bytes, 4, (timestamp >> 32n) & 0xffffn, 2);
    writeBigEndian(bytes, 6, ((timestamp >> 48n) & 0x0fffn) | 0x1000n, 2);
    bytes[8] = ((sequence >> 8) & 0x3f) | 0x80;
    bytes[9] = sequence & 0xff;
    bytes.set(getNodeId(), 10);

    return formatUuid(bytes);
}

function generateV4Uuid() {
    const bytes = randomBytes(16);
    setVersionAndVariant(bytes, 4);
    return formatUuid(bytes);
}

function generateV6Uuid() {
    const timestamp = nextGregorianTimestamp();
    const bytes = new Uint8Array(16);
    const sequence = getClockSequence();

    writeBigEndian(bytes, 0, timestamp >> 12n, 6);
    writeBigEndian(bytes, 6, (timestamp & 0x0fffn) | 0x6000n, 2);
    bytes[8] = ((sequence >> 8) & 0x3f) | 0x80;
    bytes[9] = sequence & 0xff;
    bytes.set(getNodeId(), 10);

    return formatUuid(bytes);
}

function generateV7Uuid() {
    const bytes = randomBytes(16);
    const timestamp = BigInt(Date.now());

    writeBigEndian(bytes, 0, timestamp, 6);
    setVersionAndVariant(bytes, 7);

    return formatUuid(bytes);
}

async function generateNameBasedUuid(version: 3 | 5, options: { namespace?: string; name?: string }) {
    const namespace = options.namespace?.trim();
    const name = options.name ?? '';

    if (!namespace || !isUuid(namespace)) {
        throw new Error('Enter a valid namespace UUID.');
    }

    if (!name) {
        throw new Error('Enter a name to generate this UUID.');
    }

    const namespaceBytes = parseUuid(namespace);
    const nameBytes = new TextEncoder().encode(name);
    const payload = new Uint8Array(namespaceBytes.length + nameBytes.length);
    payload.set(namespaceBytes);
    payload.set(nameBytes, namespaceBytes.length);

    const hash = version === 3 ? md5(payload) : await sha1(payload);
    const bytes = hash.slice(0, 16);
    setVersionAndVariant(bytes, version);

    return formatUuid(bytes);
}

function nextGregorianTimestamp() {
    const now = Date.now();
    let milliseconds = now;

    if (now < lastTimestampMilliseconds) {
        clockSequence = (getClockSequence() + 1) & 0x3fff;
        milliseconds = lastTimestampMilliseconds;
        lastTimestampSubMilliseconds += 1;
    } else if (now === lastTimestampMilliseconds) {
        lastTimestampSubMilliseconds += 1;
    } else {
        lastTimestampSubMilliseconds = 0;
    }

    if (lastTimestampSubMilliseconds >= 10000) {
        milliseconds = lastTimestampMilliseconds + 1;
        lastTimestampSubMilliseconds = 0;
    }

    lastTimestampMilliseconds = milliseconds;

    return BigInt(milliseconds) * 10000n + BigInt(lastTimestampSubMilliseconds) + GREGORIAN_OFFSET_100NS;
}

function getClockSequence() {
    if (clockSequence === null) {
        const bytes = randomBytes(2);
        clockSequence = (((bytes[0] << 8) | bytes[1]) & 0x3fff) >>> 0;
    }

    return clockSequence;
}

function getNodeId() {
    if (!nodeId) {
        nodeId = randomBytes(6);
        nodeId[0] |= 0x01;
    }

    return nodeId;
}

function getCrypto() {
    if (!globalThis.crypto) {
        throw new Error('This browser does not support the Crypto API.');
    }

    return globalThis.crypto;
}

function randomBytes(length: number) {
    const bytes = new Uint8Array(length);
    getCrypto().getRandomValues(bytes);
    return bytes;
}

async function sha1(bytes: Uint8Array) {
    const crypto = getCrypto();

    if (!crypto.subtle) {
        throw new Error('This browser does not support SHA-1 hashing.');
    }

    const buffer = new ArrayBuffer(bytes.byteLength);
    new Uint8Array(buffer).set(bytes);

    return new Uint8Array(await crypto.subtle.digest('SHA-1', buffer));
}

function setVersionAndVariant(bytes: Uint8Array, version: number) {
    bytes[6] = (bytes[6] & 0x0f) | (version << 4);
    bytes[8] = (bytes[8] & 0x3f) | 0x80;
}

function writeBigEndian(bytes: Uint8Array, offset: number, value: bigint, length: number) {
    let nextValue = value;

    for (let i = length - 1; i >= 0; i--) {
        bytes[offset + i] = Number(nextValue & 0xffn);
        nextValue >>= 8n;
    }
}

function parseUuid(uuid: string) {
    const hex = uuid.replace(/-/g, '');
    const bytes = new Uint8Array(16);

    for (let i = 0; i < bytes.length; i++) {
        bytes[i] = parseInt(hex.slice(i * 2, i * 2 + 2), 16);
    }

    return bytes;
}

function formatUuid(bytes: Uint8Array) {
    const hex = Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('');
    return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}

function md5(input: Uint8Array) {
    const shiftAmounts = [
        7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9,
        14, 20, 5, 9, 14, 20, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 6, 10, 15,
        21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
    ];
    const table = Array.from({ length: 64 }, (_, index) => Math.floor(Math.abs(Math.sin(index + 1)) * 2 ** 32));
    const paddedLength = getMd5PaddedLength(input.length);
    const bytes = new Uint8Array(paddedLength);
    const words = new Uint32Array(16);
    const bitLength = BigInt(input.length) * 8n;
    let a0 = 0x67452301;
    let b0 = 0xefcdab89;
    let c0 = 0x98badcfe;
    let d0 = 0x10325476;

    bytes.set(input);
    bytes[input.length] = 0x80;

    for (let i = 0; i < 8; i++) {
        bytes[paddedLength - 8 + i] = Number((bitLength >> BigInt(8 * i)) & 0xffn);
    }

    for (let offset = 0; offset < bytes.length; offset += 64) {
        for (let i = 0; i < words.length; i++) {
            words[i] =
                bytes[offset + i * 4] |
                (bytes[offset + i * 4 + 1] << 8) |
                (bytes[offset + i * 4 + 2] << 16) |
                (bytes[offset + i * 4 + 3] << 24);
        }

        let a = a0;
        let b = b0;
        let c = c0;
        let d = d0;

        for (let i = 0; i < 64; i++) {
            let f = 0;
            let g = 0;

            if (i < 16) {
                f = (b & c) | (~b & d);
                g = i;
            } else if (i < 32) {
                f = (d & b) | (~d & c);
                g = (5 * i + 1) % 16;
            } else if (i < 48) {
                f = b ^ c ^ d;
                g = (3 * i + 5) % 16;
            } else {
                f = c ^ (b | ~d);
                g = (7 * i) % 16;
            }

            const previousD = d;
            d = c;
            c = b;
            b = (b + leftRotate((a + f + table[i] + words[g]) >>> 0, shiftAmounts[i])) >>> 0;
            a = previousD;
        }

        a0 = (a0 + a) >>> 0;
        b0 = (b0 + b) >>> 0;
        c0 = (c0 + c) >>> 0;
        d0 = (d0 + d) >>> 0;
    }

    const hash = new Uint8Array(16);
    writeLittleEndian(hash, 0, a0);
    writeLittleEndian(hash, 4, b0);
    writeLittleEndian(hash, 8, c0);
    writeLittleEndian(hash, 12, d0);

    return hash;
}

function getMd5PaddedLength(inputLength: number) {
    let paddedLength = inputLength + 1;

    while (paddedLength % 64 !== 56) {
        paddedLength += 1;
    }

    return paddedLength + 8;
}

function leftRotate(value: number, amount: number) {
    return ((value << amount) | (value >>> (32 - amount))) >>> 0;
}

function writeLittleEndian(bytes: Uint8Array, offset: number, value: number) {
    bytes[offset] = value & 0xff;
    bytes[offset + 1] = (value >>> 8) & 0xff;
    bytes[offset + 2] = (value >>> 16) & 0xff;
    bytes[offset + 3] = (value >>> 24) & 0xff;
}
