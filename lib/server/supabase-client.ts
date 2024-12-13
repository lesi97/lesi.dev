import { createClient } from '@supabase/supabase-js';
import path from 'path';

const isLocal = process.env.NODE_ENV !== 'production';

const currentDir = __dirname;

if (isLocal) {
    require('dotenv').config();
} else {
    require('dotenv').config({ path: '/srv/lesi.dev/private/.env' });
}

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL || '';
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY || '';
const supabaseServiceRoleKey = process.env.SUPABASE_SERVICE_ROLE_KEY || '';

export const supabase = createClient(supabaseUrl, supabaseServiceRoleKey);
