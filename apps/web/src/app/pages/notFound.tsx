import { Button, Description } from '@/components/ui';
import { usePageMeta } from '@/hooks';
import { Link } from 'react-router-dom';

export function NotFound() {
    usePageMeta({
        title: 'Page Not Found | Lesi',
        description: 'Page Not Found',
    });
    return (
        <>
            <Description
                title='Page Not Found'
                subtitle={
                    <>
                        <Link to='/'>
                            <Button variant='secondary'>Return to home</Button>
                        </Link>
                    </>
                }
            />
        </>
    );
}
