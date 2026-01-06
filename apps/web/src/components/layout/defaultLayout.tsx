import { Outlet } from 'react-router-dom';

export function DefaultLayout() {
    return (
        <main className='relative top-8 mb-8 flex h-fit w-11/12 justify-center rounded-lg bg-base-100 px-8 py-8 shadow peer-[data-state="open"]:opacity-0 2xl:w-50% 2xl:min-w-50%'>
            <div className='flex w-11/12 flex-col gap-4'>
                <Outlet />
            </div>
        </main>
    );
}
