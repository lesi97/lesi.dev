import { Outlet } from 'react-router-dom';
import { Footer } from './footer';

export function WideLayout() {
    return (
        <>
            <main className='relative top-8 mb-8 flex h-fit w-full justify-center rounded-lg bg-base-100 px-8 py-8 shadow xl:w-fit xl:min-w-50%'>
                <div className='flex w-full flex-col gap-4'>
                    <div>
                        <Outlet />
                    </div>
                </div>
            </main>
            <Footer />
        </>
    );
}
