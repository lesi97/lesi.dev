import { Http_Status_Code } from '@/lib/server/enums/enums';
import { getApiDetails } from '../../_lib/query-api-details';
import { NextApiRequest, NextApiResponse } from 'next';

export default async function GET(req: NextApiRequest, res: NextApiResponse) {
    const { twitch_user } = req.query;
    if (typeof twitch_user !== 'string') {
        // Always return 200, this is used as a Twitch Nightbot command and if the status code is an error
        // The status code is returned instead of a custom error message, e.g a 400 returns SERVER RESPONDED WITH 400 or something similar
        res.setHeader('Content-Type', 'text/plain; charset=utf-8');
        return res.status(Http_Status_Code.OK).send('Invalid Twitch ID');
    }
    const envDetails = await getApiDetails('Twitch');

    //TODO

    // Always return 200, this is used as a Twitch Nightbot command and if the status code is an error
    // The status code is returned instead of a custom error message, e.g a 400 returns SERVER RESPONDED WITH 400 or something similar
    res.setHeader('Content-Type', 'text/plain; charset=utf-8');
    return res.status(Http_Status_Code.OK).send('Invalid Twitch ID');
}
