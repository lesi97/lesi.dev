import { PostgrestError } from '@supabase/postgrest-js';

export type Bungie_User = {
    username: string;
    username_code: number;
    membership_id: string;
    preferred_platform: number;
    friendly_name: string;
    last_played_character_id: string;
    characters: Bungie_Characters[];
};

export type Bungie_Characters = {
    character_id: string;
    character_type: string;
    time_played: number;
};

export type Bungie_Weapon = {
    name: string;
    id: string;
    hash_id: string;
    is_exotic: boolean | null;
    perks: Bungie_Perks[];
    mod: Bungie_Perks[];
    shader: Bungie_Perks[];
    masterwork: Bungie_Perks[];
    ornament: Bungie_Perks[];
    pvp_kills: number | null;
    pve_kills: number | null;
    trials_kills: number | null;
    active_tracker: 'pve' | 'pvp' | 'trials' | null;
};

export type Bungie_Perks = {
    name: string;
    item_type: string;
};
