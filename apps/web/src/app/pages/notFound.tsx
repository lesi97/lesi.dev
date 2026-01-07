import { Button, Description } from '@/components/ui';
import { Link } from 'react-router-dom';

export function NotFound() {
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
