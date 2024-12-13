import { getPlatformEnum } from '../_lib/bungie/get-platform-enum';
import { Bungie_Current_User_Data } from '@/lib/server/data/bungie/bungie';
import { Http_Status_Code } from '@/lib/server/enums/enums';
import { NextApiRequest, NextApiResponse } from 'next';

export default async function GET(req: NextApiRequest, res: NextApiResponse) {
    const { bungie_id, platform } = req.query;
    const platformEnum = getPlatformEnum(platform);
    const user = await Bungie_Current_User_Data.create(bungie_id, platformEnum);
    const data = user.time;

    // Always return 200, this is used as a Twitch Nightbot command and if the status code is an error
    // The status code is returned instead of a custom error message, e.g a 400 returns SERVER RESPONDED WITH 400 or something similar
    res.setHeader('Content-Type', 'text/plain; charset=utf-8');
    return res.status(Http_Status_Code.OK).send(data.message);
}
