import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { supabase } from '@/lib/supabase/createClient';

export function AuthCallback() {
    const navigate = useNavigate();

    useEffect(() => {
        async function run() {
            const url = new URL(window.location.href);
            const code = url.searchParams.get('code');
            const next = url.searchParams.get('next') ?? '/';

            if (!code) {
                navigate('/error', { replace: true });
                return;
            }

            const { error } = await supabase.auth.exchangeCodeForSession(code);

            if (error) {
                throw error;
                console.error(error);
            }

            navigate(next, { replace: true });
        }

        run();
    }, [navigate]);

    return <div>Signing you in...</div>;
}
