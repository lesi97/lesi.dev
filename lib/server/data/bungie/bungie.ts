import { StringLiteral } from "typescript";
import { HTTP_Response } from "../../types/types";
import { selectSql, insertSql, upsertSql } from "../database/helpers";
import { getApiDetails } from "../database/database";
import { FilterArray, Api_Details_Response, PlugArray } from "../database/database.types";
import { Bungie_User, Bungie_Characters, Bungie_Weapon, } from "./bungie.types";
import * as helpers from './helpers'
import { Http_Status_Code } from "../../enums/enums";
import { fetchPerks } from "./helpers/fetch_perks";

export class Bungie_Current_User_Data {
    _bungie_id: string;
    _weapon_slot: 0 | 1 | 2;
    _user_details: Bungie_User;
    _weapon: Bungie_Weapon;
    _error: HTTP_Response | null;
    #_api_key: string | null = null;
    #_base_url: string | null = null;
    
    constructor(bungie_id: string, preferred_platform: number, weapon_slot: 0 | 1 | 2) {
        this._bungie_id = bungie_id;
        this._weapon_slot = weapon_slot || 0;
        this._user_details = {
            username: '',
            username_code: 0,
            membership_id: '',
            preferred_platform: preferred_platform,
            friendly_name: '',
            last_played_character_id: '',
            characters: []
        };
        this._weapon = {
            name: '',
            id: '',
            hash_id: '',
            is_exotic: null,
            perks: [],
            mod: [],
            shader: [],
            masterwork: [],
            ornament: [],
            pve_kills: null,
            pvp_kills: null,
            trials_kills: null,
            active_tracker: null,
        }
        this._error = null;
    }

    public get current_weapon(): HTTP_Response { return this.generate_current_weapon_string() };
    public get time(): HTTP_Response { return this.generate_play_time_string() };

    public static async create(bungie_id: string, preferred_platform: number, weapon_slot: 0 | 1 | 2) {
        const instance = new Bungie_Current_User_Data(bungie_id, preferred_platform, weapon_slot);
        await instance.initialize();
        if (!instance._error) {
            await instance.fetchBungieProfile();
        }
        return instance;
    }

    private async initialize() {
        const [username, username_code, id_validation_error] = helpers.validateId(this._bungie_id);
        if (id_validation_error) return this._error = id_validation_error;
        this._user_details.username = username as string;
        this._user_details.username_code = username_code as number;        
        
        const [api_details, api_query_error] = await getApiDetails('Bungie');
        if (api_query_error) return this._error = api_query_error;
        this.#_api_key = api_details.client_id;
        this.#_base_url = api_details.base_url;

        // await fetchPerks();
        
        await this.fetchUserFromDatabase();        
    }

    private async fetchUserFromDatabase() {
        const table = 'destiny_users';
        const columns = 'membership_id, preferred_platform, friendly_name';
        const filter: FilterArray = { bungie_id: { value: this._bungie_id, filter_type: 'ilike' } };
        
        const [user, user_query_error] = await selectSql(table, columns, filter, null );
        
        if (user_query_error && user_query_error !== 'No Data') {
            return this._error = user_query_error;
        }
        if (!user || user.length <= 0) {
            return await this.findBungieUserDetails();
        }

        this._user_details.membership_id = user[0].membership_id;
        this._user_details.preferred_platform = parseInt(user[0].preferred_platform);
        this._user_details.friendly_name = user[0].friendly_name;
        
        await this.fetchBungieProfile();
    }

    private async findBungieUserDetails() {        
        const headersList = {
            Accept: '*/*',
            'x-api-key': `${this.#_api_key}`,
        };
        const id = encodeURIComponent(this._bungie_id);
        const url = new URL(`${this.#_base_url}/Destiny2/SearchDestinyPlayer/-1/${id}/`);
        
        const response = await fetch(url, {
            method: 'GET',
            headers: headersList,
        });
        if (!response.ok) return this._error = {message: 'bungie api is currently down, gift a sub instead and try again later', status_code: 500};
        const data = await response.json();
        const main_profile = data.Response[0];

        this._user_details.membership_id = main_profile.membershipId.toString();
        this._user_details.friendly_name = main_profile.displayName;
        this._user_details.preferred_platform = main_profile.membershipType;

        const user = {
            membership_id: this._user_details.membership_id, 
            bungie_id: this._bungie_id, 
            preferred_platform: this._user_details.preferred_platform, 
            friendly_name: this._user_details.friendly_name 
        };
        this.fetchBungieProfile();
        insertSql('destiny_users', user);
    }

    private async fetchBungieProfile() {
        const headersList = {
            Accept: '*/*',
            'x-api-key': `${this.#_api_key}`,
        };
        const url = new URL(`${this.#_base_url}/Destiny2/${this._user_details.preferred_platform}/Profile/${this._user_details.membership_id}/?components=200,205,302,305,309`);
        const response = await fetch(url, { 
            method: 'GET', 
            headers: headersList 
        });
        if (!response.ok) return this._error = {message: 'bungie api is currently down, gift a sub instead and try again later', status_code: 500};
        
        const data = await response.json();
        
        const characters = data.Response.characters.data;
        const [character_data, characters_array, last_played_character] = helpers.getCharacterData(characters, this._user_details.membership_id);

        Object.assign(this._user_details, {
            characters: characters_array,
            last_played_character_id: last_played_character.characterId,
        });        

        const character_equipment = data.Response.characterEquipment.data[`${this._user_details.last_played_character_id}`].items[this._weapon_slot];
        this._weapon.id = character_equipment.itemInstanceId.toString();
        const plugHashes = data.Response.itemComponents.sockets.data[`${this._weapon.id}`].sockets.map((perk: any) => perk.plugHash).filter(Boolean);
        const plug_objectives = data.Response.itemComponents.plugObjectives.data[`${this._weapon.id}`]?.objectivesPerPlug;

        // console.log(plugHashes)

        const pvpIdentifiers = [38912240, 2302094943, 231101171, 3244015567];
        const pveIdentifiers = [2240097604, 2285636663, 231101171, 3915764595, 3915764594, 3915764593, 1187045864, 1690059054, 16638393, 1124054883, 3624435060, 16638392, 2617715132, 2617715133, 2302094943, 905869860, 2240097604 ];
        const trialsIdentifiers = [3915764595];
        // const gambitIdentifiers = []; // Source https://www.bungie.net/common/destiny2_content/json/en/DestinyPlugSetDefinition-7637b37a-babb-4b0a-aa1d-d259812e43ff.json

        // console.log(plug_objectives)
        Object.assign(this._weapon, {
            // id: character_equipment.itemInstanceId.toString(),
            hash_id: character_equipment.itemHash.toString(),
            pvp_kills: helpers.getKillCounts(plug_objectives, pvpIdentifiers),
            pve_kills: helpers.getKillCounts(plug_objectives, pveIdentifiers),
            trials_kills: helpers.getKillCounts(plug_objectives, trialsIdentifiers),
        });

        // console.log(this._weapon.pve_kills)

        await this.queryWeaponData(plugHashes);        

        const userCharactersTable = 'destiny_user_characters';
        const userCharactersConflict = 'membership_id, character_id';
        upsertSql(userCharactersTable, character_data, userCharactersConflict);
    }

    private async queryWeaponData(plugHashes) {
        const table = 'destiny_weapons';
        const columns = 'display_name, tier_type_name';
        const filter: FilterArray = { id: { value: this._weapon.hash_id, filter_type: 'eq' } };
        const [weaponData, weapon_data_query_error] = await selectSql(table, columns, filter);

        if (weapon_data_query_error) return this._error = weapon_data_query_error;
        if (weaponData.length === 0) {
            helpers.fetchWeapons();
            return this._error = `Weapon not found, please try again shortly, I'm updating my records 🤓`;
        } else {
            this._weapon.name = weaponData[0].display_name;
            this._weapon.is_exotic = weaponData[0].tier_type_name === 'Exotic';
        }     
        
        const perkTable = 'destiny_weapon_perks';
        const perkColumns = 'name, item_type';
        const perkFilter: FilterArray = { hash_id: { value: plugHashes, filter_type: 'in' } };
        const [perkData, error] = await selectSql(perkTable, perkColumns, perkFilter);

        // console.log(perkFilter)
        
        const allPerks = perkData || [];
        const weaponPerks = allPerks.filter((perk) => 
            !perk.item_type.includes('Shader') &&
            !perk.item_type.includes('Mod') &&
            !perk.item_type.includes('Ornament') &&
            !perk.name.includes('Masterwork') &&
            !perk.name.includes('Catalyst') &&
            !perk.name.includes('Default Shader')
        );

        const mods = allPerks
            .filter((perk) => perk.item_type.includes('Mod'))
            .map((mod) => ({
                ...mod,
                name: mod.item_type.includes('Enhanced') ? `Enhanced ${mod.name}` : mod.name,
        }));
        
        this._weapon.perks = weaponPerks;
        this._weapon.mod = allPerks.filter((perk) => perk.item_type.includes('Mod'));
        this._weapon.masterwork = allPerks.filter((perk) => perk.name.includes('Masterwork'));
        this._weapon.shader = allPerks.filter((perk) => perk.item_type === 'Shader');
        this._weapon.ornament = allPerks.filter((perk) => perk.item_type.includes('Ornament'));

        const weaponTable = 'destiny_weapon_kill_counts';
        const weaponInsertData = {
            membership_id: this._user_details.membership_id,
            weapon_id: this._weapon.id,
            pvp_kills: this._weapon.pvp_kills || null,
            pve_kills: this._weapon.pve_kills || null,
            trials_kills: this._weapon.trials_kills || null,
            weapon_hash: this._weapon.hash_id
        };
        const conflict = 'membership_id, weapon_id';
        upsertSql(weaponTable, weaponInsertData, conflict);
    }
    

    

    private formatPlayTime(play_time: number): { hours: number; minutes: number } {
        const hours = Math.floor(play_time / 60);
        const minutes = play_time % 60;
        return { hours, minutes };
    }

    private generate_current_weapon_string(): HTTP_Response {
        if (this._error) return this._error;

        const weapon = this._weapon;
        // if (!weapon.name) return 'error';
 
        const perksRaw = weapon.perks.map((perk) => perk.name).join(', ');
        const perks = `| Perks: ${perksRaw}`;

        const mod = weapon.is_exotic ? '' 
            : weapon.mod.length > 0 
                ? `| Mod: ${weapon.mod[0].name}`
                : `| Mod: None`;
                
        const masterwork = weapon.is_exotic ? '' 
        : weapon.masterwork.length > 0
            ? `| ${weapon.masterwork[0].name}`
            : ``;

        const shader = weapon.is_exotic ? '' 
        : weapon.shader.length > 0 
            ? `| Shader: ${weapon.shader[0].name}`
            : `| Shader: None}`;

        const ornament = weapon.ornament.length > 0
        ? `| Ornament: ${weapon.ornament[0].name}`
        : ``;

        const tracker = weapon.pvp_kills
        ? `| PVP Kill Count: ${weapon.pvp_kills.toLocaleString()}`
        : weapon.pve_kills
          ? `| PVE Kill Count: ${weapon.pve_kills.toLocaleString()}`
          : weapon.trials_kills
            ? `| Trials Kill Count: ${weapon.trials_kills.toLocaleString()}`
            : '';

        const message = `${this._bungie_id}: ${weapon.name} ${perks} ${mod} ${masterwork} ${shader} ${ornament} ${tracker}`.replace(/\s+/g, ' ').trim();
        return {message, status_code: Http_Status_Code.OK};
        
    }

    private generate_play_time_string() {
    // console.log(this._error)
        if (this._error) return this._error;
        
        // console.log(this._user_details.characters)
        
        const user_characters = this._user_details.characters;
        // if (user_characters.length < 0) return 'error';

        const characters_played_time = user_characters.map((character) => {
                const { hours, minutes } = this.formatPlayTime(character.time_played);

                const hoursText = hours === 1 ? '1 hour' : `${hours} hours`;
                const minutesText = minutes === 1 ? '1 minute' : `${minutes} minutes`;

                if (hours === 0) {
                    return `${character.character_type}: ${minutesText}`;
                } else if (minutes === 0) {
                    return `${character.character_type}: ${hoursText}`;
                } else {
                    return `${character.character_type}: ${hoursText} & ${minutesText}`;
                }
            })
            .join(' | ');
        const message = `${this._bungie_id} - ${characters_played_time}`;
        return { message, status_code: Http_Status_Code.OK }; 
    }
}