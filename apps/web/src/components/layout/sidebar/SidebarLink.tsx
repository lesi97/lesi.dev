'use client';
import React from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useSidebar } from './SidebarContext';
import { mergeClassNames } from '@/utils';
import { LinkProps } from 'react-router-dom';

export interface Links {
    label: string;
    href: string;
    icon: React.JSX.Element | React.ReactNode;
    onClick?: () => void;
}

const SidebarLink = ({ link, className, ...props }: { link: Links; className?: string; props?: LinkProps }) => {
    const { open, animate } = useSidebar();
    const handleClick = async (event: React.MouseEvent) => {
        if (link.onClick) {
            event.preventDefault();
            await link.onClick();
        }
    };
    return (
        <Link
            to={link.href}
            className={mergeClassNames('group/sidebar relative flex items-center justify-start gap-2 py-2', className)}
            onClick={handleClick}
            {...props}>
            {link.icon}
            <motion.span
                animate={{
                    display: animate ? (open ? 'inline-block' : 'none') : 'inline-block',
                    opacity: animate ? (open ? 1 : 0) : 1,
                }}
                className='!m-0 inline-block overflow-hidden whitespace-pre !p-0 text-sm text-neutral-700 transition duration-150 group-hover/sidebar:translate-x-1'>
                {link.label}
            </motion.span>
        </Link>
    );
};

export default SidebarLink;
