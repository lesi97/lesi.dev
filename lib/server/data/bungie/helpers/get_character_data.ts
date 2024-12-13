import { Bungie_Characters } from '../bungie.types';
import { getCharacterTypeEnum } from './get_character_type_enum';

export function getCharacterData(characters, membership_id: string) {
    let last_played_character: Record<string, string>;
    const characters_array: Bungie_Characters[] = [];

    const character_data = Object.values(characters).map((character: any) => {
        const character_type = getCharacterTypeEnum(character.classType);
        if (
            !last_played_character ||
            new Date(character.dateLastPlayed) > new Date(last_played_character.dateLastPlayed)
        ) {
            last_played_character = character;
        }
        characters_array.push({
            character_id: character.characterId,
            character_type,
            time_played: character.minutesPlayedTotal,
        });
        return {
            membership_id: membership_id,
            character_id: character.characterId,
            character_type,
            minutes_played: character.minutesPlayedTotal,
        };
    });

    const sorted_characters = characters_array.sort((a, b) => b.time_played - a.time_played);

    return [character_data, sorted_characters, last_played_character];
}
