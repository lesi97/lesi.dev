import crypto from 'crypto';

const ENCRYPTION_KEY = process.env.ENCRYPTION_KEY || '';
const ENCRYPTION_IV = process.env.ENCRYPTION_IV || '';

export function encrypt(text: string): string {
    const key = Buffer.from(ENCRYPTION_KEY, 'hex');
    const iv = Buffer.from(ENCRYPTION_IV, 'hex');
    const cipher = crypto.createCipheriv('aes-256-cbc', key, iv);
    let encrypted = cipher.update(text);
    encrypted = Buffer.concat([encrypted, cipher.final()]);
    return encrypted.toString('hex');
}

export function decrypt(encryptedText: string): string {
    const key = Buffer.from(ENCRYPTION_KEY, 'hex');
    const iv = Buffer.from(ENCRYPTION_IV, 'hex');
    const decipher = crypto.createDecipheriv('aes-256-cbc', key, iv);
    let decrypted = decipher.update(Buffer.from(encryptedText, 'hex'));
    decrypted = Buffer.concat([decrypted, decipher.final()]);
    return decrypted.toString();
}
