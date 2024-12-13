export function isValidBungieId(bungie_id: string): boolean {
    return /^[a-zA-Z0-9]+#[0-9]{4}$/.test(bungie_id);
}
