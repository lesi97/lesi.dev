'use client';
import React, { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Icons } from '@/components/ui';
import { SidebarBody, SidebarLink } from '.';
import { SidebarProvider } from './SidebarContext';
import { mergeClassNames } from '@/utils';

export default function SidebarNavigation() {
    const [open, setOpen] = useState(false);

    const profileRef = useRef<HTMLDivElement>(null);
    const avatarRef = useRef<HTMLDivElement>(null);

    const links = [
        {
            label: 'Home',
            href: `/editor`,
            icon: <Icons.Menu className='h-5 w-5 flex-shrink-0 text-neutral-700' />,
            requiresAdmin: false,
            viewOnMobile: true,
            viewOnDesktop: true,
        },
        {
            label: 'Sale',
            href: `/editor/sale`,
            icon: <Icons.Menu className='h-5 w-5 flex-shrink-0 text-neutral-700' />,
            requiresAdmin: false,
            viewOnMobile: true,
            viewOnDesktop: true,
        },
    ];

    return (
        <>
            <SidebarProvider open={open} setOpen={setOpen}>
                <SidebarBody className='bg-primary justify-between gap-10 text-white'>
                    <div className='flex flex-1 flex-col overflow-y-auto overflow-x-hidden'>
                        <Link
                            to={`/editor`}
                            className='relative z-20 flex items-center space-x-2 py-1 text-sm font-normal'>
                            {open ? <Logo /> : <LogoIcon />}
                        </Link>
                        <div className='mt-8 flex flex-col gap-2'>
                            {links.map((link, idx) => {
                                return (
                                    <SidebarLink
                                        key={idx}
                                        link={link}
                                        className={mergeClassNames(
                                            link.viewOnMobile ? 'flex' : 'hidden',
                                            link.viewOnDesktop ? 'lg:flex' : 'lg:hidden'
                                        )}
                                    />
                                );
                            })}
                        </div>
                    </div>
                </SidebarBody>
            </SidebarProvider>
        </>
    );
}

export const Logo = () => {
    return (
        <span className='relative z-20 flex items-center space-x-2 py-1 text-sm font-normal'>
            <Icons.Menu height={20} width={20} />
            <motion.span initial={{ opacity: 0 }} animate={{ opacity: 1 }} className='whitespace-pre font-medium'>
                <Icons.Menu height={20} width={80} />
            </motion.span>
        </span>
    );
};

export const LogoIcon = () => {
    return (
        <span className='relative z-20 flex items-center space-x-2 py-1 text-sm font-normal'>
            <Icons.Menu height={20} width={20} />
        </span>
    );
};
