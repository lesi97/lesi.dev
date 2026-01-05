import { Outlet } from 'react-router-dom';
import { ReactNode } from 'react';

type Props = { className?: string };

export function PlainLayout({ className }: Props) {
    return (
        <main
            className={`relative flex h-full w-full flex-col justify-between gap-8 rounded-lg px-8 py-8 ${className ?? ''}`}>
            <Outlet />
        </main>
    );
}
