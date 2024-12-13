import { PostgrestError } from '@supabase/supabase-js';

export type Filter = {
    value: any;
    filter_type?: 'eq' | 'ilike' | 'like' | 'in' | 'gte' | 'lte';
};

export type FilterArray = Record<string, Filter>;

export type SelectSqlResponse<T> = {
    data: T[] | null;
    error: boolean | null;
};

export type Api_Details = {
    name: string | null;
    client_id: string | null;
    client_secret: string | null;
    access_token: string | null;
    refresh_token: string | null;
    refresh_token_expiry: number | null;
    base_url: string | null;
    redirect_url: string | null;
};

export type Api_Details_Response = {
    data: Api_Details[] | null;
    error: PostgrestError | null;
};

export type Plug = {
    name: string;
    item_type: string;
};

export type PlugArray = Plug[];
