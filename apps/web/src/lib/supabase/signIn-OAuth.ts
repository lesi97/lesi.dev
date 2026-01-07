import type { Provider } from '@supabase/supabase-js';
import { supabase } from '@/lib/supabase/createClient';

export async function signInOauth(provider: Provider, redirectTo?: string): Promise<void> {
    alert(redirectTo);
    const { data, error } = await supabase.auth.signInWithOAuth({
        provider,
        options: {
            redirectTo,
            queryParams: {
                access_type: 'offline',
                prompt: 'consent',
            },
        },
    });

    if (error) {
        window.location.assign('/error');
        return;
    }

    window.location.assign(data.url);
}
