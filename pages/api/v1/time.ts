import { Http_Status_Code } from '@/lib/server/enums/enums';
import { format, toZonedTime } from 'date-fns-tz';

import { NextApiResponse, NextApiRequest } from 'next';

export default async function GET(req: NextApiRequest, res: NextApiResponse) {
    const { zone } = req.query;
    const activeZone = zone || 'UTC';

    const now = new Date();
    const zonedTime = toZonedTime(now, activeZone);
    const formattedTime = format(zonedTime, 'HH:mm:ss', { timeZone: activeZone });
    const formattedDate = format(zonedTime, 'dd-MM-yyyy', { timeZone: activeZone });

    const response = { message: { date: formattedDate, time: formattedTime } };
    res.setHeader('Content-Type', 'application/json; charset=utf-8');
    return res.status(Http_Status_Code.OK).send(response);
}
