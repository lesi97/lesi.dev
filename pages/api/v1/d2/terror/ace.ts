import { Http_Status_Code } from '@/lib/server/enums/enums';
import { getKills, updateKills } from '../_lib/supabase/kill-counts';
import { getApiDetails } from '@/lib/server/data/database/database';
import { ApiDetailsResponse } from '@/types/db-queries';
import { NextApiRequest, NextApiResponse } from 'next';

export default async function GET(req: NextApiRequest, res: NextApiResponse) {
    // console.log("REQ", request.method);
    const apiDetails: ApiDetailsResponse = await getApiDetails('Bungie');

    if (apiDetails.error) {
        return console.error('Error fetching API details:', apiDetails.error.message);
    }

    const apiKey = apiDetails[0].client_id;
    const baseUrl = apiDetails[0].base_url;

    const headersList = {
        Accept: '*/*',
        'x-api-key': `${apiKey}`,
    };

    const endpointType = 'Destiny2';
    const membershipType = '3';
    const destinyMembershipId = 4611686018467358417n;

    // const warlock = 2305843009301476854;
    // const hunter = 2305843009321995500;
    // const titan = 2305843009369808628;

    const characterEquipment = 205;
    const itemPlugObjectives = 309;
    const components = `?components=${characterEquipment},${itemPlugObjectives}`;

    const weapon = 6917529207684081719n;
    const crucibleTracker = 38912240;

    const url = new URL(`${baseUrl}/${endpointType}/${membershipType}/Profile/${destinyMembershipId}/${components}`);

    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: headersList,
        });

        if (response.ok) {
            const data = await response.json();

            const killCountInt = parseInt(
                data.Response.itemComponents.plugObjectives.data[`${weapon}`].objectivesPerPlug[`${crucibleTracker}`][
                    '0'
                ].progress
            );
            const killCountHR = killCountInt.toLocaleString();

            updateKills('4611686018467358417', '6917529207684081719', '347366834', killCountInt).catch((error) =>
                console.error('Database update failed:', error)
            );

            // Always return 200, this is used as a Twitch Nightbot command and if the status code is an error
            // The status code is returned instead of a custom error message, e.g a 400 returns SERVER RESPONDED WITH 400 or something similar
            res.setHeader('Content-Type', 'text/plain; charset=utf-8');
            return res.status(Http_Status_Code.OK).send(`${killCountHR}`);
        }
    } catch (error) {
        const data = await getKills(destinyMembershipId, weapon);
        const killCountInt = data.pvp_kills;
        const killCountHR = killCountInt.toLocaleString();
        // Always return 200, this is used as a Twitch Nightbot command and if the status code is an error
        // The status code is returned instead of a custom error message, e.g a 400 returns SERVER RESPONDED WITH 400 or something similar
        res.setHeader('Content-Type', 'text/plain; charset=utf-8');
        return res.status(Http_Status_Code.OK).send(`${killCountHR}`);
    }
}
