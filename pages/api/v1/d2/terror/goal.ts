export const maxDuration = 60;
import { NextApiResponse, NextApiRequest } from 'next';
import { NextResponse } from 'next/server';
import { NextRequest } from 'next/server';
import { Http_Status_Code } from '@/lib/server/enums/enums';
import { headers } from 'next/headers';

import { insertSql, selectSql, upsertSql } from '@/lib/server/data/database/helpers';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
    const method = req.method;

    if (method === 'GET') {
        return await handleGet(req, res);
    }

    if (method === 'POST') {
        return await handlePost(req, res);
    }

    const response = { message: null, error: 'Invalid Method' };
    res.setHeader('Content-Type', 'text/plain; charset=utf-8');
    return res.status(Http_Status_Code.MethodNotAllowed).send(response);
}

async function handleGet(req: NextApiRequest, res: NextApiResponse) {
    const { subs } = req.query;
    const headersList = req.headers;
    const referer = headersList.referer;
    const origin = headersList.origin;
    const allowedOrigin = process.env.NEXT_PUBLIC_ROOT_URL;

    const table = 'terror_goal_messages';
    const columns = '*';

    const [message, error] = await selectSql(table, columns, null);
    const first = message[0];
    console.log(referer);
    // if (referer.includes(process.env.NEXT_PUBLIC_ROOT_URL)) {
    if (referer?.includes('terror/goal')) {
        const response = JSON.stringify(message);
        //console.log('response', response);
        res.setHeader('Content-Type', 'application/json; charset=utf-8,');
        // res.setHeader('Access-Control-Allow-Origin', '*');
        return res.status(Http_Status_Code.OK).send(response);
    }

    let response = '';
    if (first.message.includes('[subs]') && subs) {
        response = first.message.replace('[subs]', subs);
    } else if (first.message.includes('[subs]') && !subs) {
        response = first.message.replace('[subs]', '50');
    } else {
        response = `${subs ? `${subs} ` : ''}${first.message}`;
    }

    // Always return 200, this is used as a Twitch Nightbot command and if the status code is an error
    // The status code is returned instead of a custom error message, e.g a 400 returns SERVER RESPONDED WITH 400 or something similar
    res.setHeader('Content-Type', 'text/plain; charset=utf-8');
    return res.status(Http_Status_Code.OK).send(response);
}

async function handlePost(req: NextApiRequest, res: NextApiResponse) {
    const { message } = await req.json();
    const table = 'terror_goal';
    if (!message) {
        const response = { message: null, error: 'Invalid request: message is required' };
        res.setHeader('Content-Type', 'application/json; charset=utf-8');
        return res.status(Http_Status_Code.BadRequest).send(response);
    }
    const data = { message };
    const [error] = await insertSql(table, data);
    if (error) {
        const response = { message: null, error: 'Error updating message' };
        res.setHeader('Content-Type', 'application/json; charset=utf-8');
        return res.status(Http_Status_Code.BadRequest).send(response);
    }

    const response = { message: 'Message updated successfully', error: null };
    res.setHeader('Content-Type', 'application/json; charset=utf-8');
    return res.status(Http_Status_Code.OK).send(response);
}
