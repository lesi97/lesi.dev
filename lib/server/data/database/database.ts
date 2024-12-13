import { supabase } from '@/lib/server/supabase-client';
import { encrypt, decrypt } from '@/lib/server/encryption';
import * as Database from './database.types';
import { selectSql, upsertSql } from './helpers';
import { PostgrestError } from '@supabase/supabase-js';

export async function getApiDetails(
    application: string
): Promise<[Record<string, string>[], null] | [null, PostgrestError] | [null, string]> {
    const table = 'api_keys';
    const columns =
        'name, client_id, client_secret, access_token, refresh_token, refresh_token_expiry, base_url, redirect_url';
    const filter = { name: { value: application, filter_type: 'eq' } };
    const [data, error] = await selectSql(table, columns, filter);

    if (error) return error;

    const unencryptedData = {
        name: data[0].name,
        client_id: data[0].client_id ? decrypt(data[0].client_id) : null,
        client_secret: data[0].client_secret ? decrypt(data[0].client_secret) : null,
        access_token: data[0].access_token ? decrypt(data[0].access_token) : null,
        refresh_token: data[0].refresh_token ? decrypt(data[0].refresh_token) : null,
        refresh_token_expiry: data[0].refresh_token_expiry ? data[0].refresh_token_expiry : null,
        base_url: data[0].base_url ? data[0].base_url : null,
        redirect_url: data[0].redirect_url ? data[0].redirect_url : null,
    };

    return [unencryptedData, error];
}

export async function updateApiDetails(data: Database.Api_Details) {
    const table = 'api_keys';
    const conflict = 'name';
    const dataToUpdate: Database.Api_Details = {
        name: data.name ? data.name : null,
        client_id: data.client_id ? encrypt(data.client_id) : null,
        client_secret: data.client_secret ? encrypt(data.client_secret) : null,
        access_token: data.access_token ? encrypt(data.access_token) : null,
        refresh_token: data.refresh_token ? encrypt(data.refresh_token) : null,
        refresh_token_expiry: data.refresh_token_expiry ? data.refresh_token_expiry : null,
        base_url: data.base_url ? data.base_url : null,
        redirect_url: data.redirect_url ? data.redirect_url : null,
    };

    const { error } = await upsertSql(table, dataToUpdate, conflict);

    if (error) return [error];
    return [null];
}
