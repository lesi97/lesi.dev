'use client';
import React from 'react';
import { motion } from 'framer-motion';
import { mergeClassNames } from '@/utils';
import { useSidebar } from './SidebarContext';

const DesktopSidebar = ({ className, children, ...props }: React.ComponentProps<typeof motion.div>) => {
    const { open, setOpen, animate } = useSidebar();
    return (
        <motion.div
            className={mergeClassNames(
                'hidden h-full w-[220px] flex-shrink-0 bg-neutral-100 px-4 py-4 md:flex md:flex-col',
                className
            )}
            animate={{
                width: animate ? (open ? '220px' : '60px') : '220px',
            }}
            onMouseEnter={() => setOpen(true)}
            onMouseLeave={() => setOpen(false)}
            {...props}>
            {children}
        </motion.div>
    );
};

export default DesktopSidebar;
