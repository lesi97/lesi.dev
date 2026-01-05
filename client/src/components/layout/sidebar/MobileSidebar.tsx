'use client';
import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Icons } from '@/components/ui';
import { mergeClassNames } from '@/utils';
import { useSidebar } from './SidebarContext';

const MobileSidebar = ({ className, children, ...props }: React.ComponentProps<'div'>) => {
    const { open, setOpen } = useSidebar();
    return (
        <>
            <div
                className={mergeClassNames(
                    'fixed top-0 z-10 flex h-10 w-full flex-row items-center justify-between bg-neutral-100 px-4 py-4 md:hidden'
                )}
                {...props}>
                <div className='z-10 flex w-full justify-end'>
                    <Icons.Menu className='cursor-pointer text-neutral-800' onClick={() => setOpen(!open)} />
                </div>
                <AnimatePresence>
                    {open && (
                        <motion.div
                            initial={{ x: '-100%', opacity: 0 }}
                            animate={{ x: 0, opacity: 1 }}
                            exit={{ x: '-100%', opacity: 0 }}
                            transition={{
                                duration: 0.3,
                                ease: 'easeInOut',
                            }}
                            className={mergeClassNames(
                                'fixed inset-0 z-[100] flex h-full w-full flex-col justify-between bg-white p-10',
                                className
                            )}>
                            <div
                                className='absolute right-10 top-10 z-50 cursor-pointer text-neutral-800'
                                onClick={() => setOpen(!open)}>
                                <Icons.X />
                            </div>
                            {children}
                        </motion.div>
                    )}
                </AnimatePresence>
            </div>
        </>
    );
};

export default MobileSidebar;
