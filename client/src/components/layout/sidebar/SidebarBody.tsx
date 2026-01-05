'use client';
import React from 'react';
import { motion } from 'framer-motion';
import { DesktopSidebar, MobileSidebar } from '.';

const SidebarBody = (props: React.ComponentProps<typeof motion.div>) => {
    return (
        <>
            <DesktopSidebar {...props} />
            <MobileSidebar {...(props as React.ComponentProps<'div'>)} />
        </>
    );
};

export default SidebarBody;
