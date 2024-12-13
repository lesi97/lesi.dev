import { Http_Status_Code } from '@/lib/server/enums/enums';
import { getApiDetails, updateApiDetails } from './_lib/query-api-details';
import { NextApiResponse, NextApiRequest } from 'next';

export default async function POST(req: NextApiRequest, res: NextApiResponse) {
    if (req.method !== 'POST') {
        const response = { error: 'Only POST requests allowed' };
        res.setHeader('Content-Type', 'application/json; charset=utf-8');
        return res.status(Http_Status_Code.MethodNotAllowed).send(response);
    }
    // const data = await req.json();
    // const envDetails = await getApiDetails('Anilist');
    // console.log(envDetails);
    //TODO

    const response = { message: 'Success' };
    res.setHeader('Content-Type', 'application/json; charset=utf-8');
    return res.status(Http_Status_Code.MethodNotAllowed).send(response);
}
