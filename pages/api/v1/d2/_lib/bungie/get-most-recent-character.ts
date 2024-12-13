export function getMostRecentCharacter(characters: any) {
    let lastPlayedChar = '';

    Object.values(characters).forEach((character: any) => {
        if (!lastPlayedChar || new Date(character.dateLastPlayed) > new Date(lastPlayedChar.dateLastPlayed)) {
            lastPlayedChar = character;
        }
    });

    return lastPlayedChar;
}
