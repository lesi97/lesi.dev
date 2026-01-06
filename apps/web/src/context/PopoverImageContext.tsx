import React, { createContext, useContext, useState, ReactNode } from 'react';
import { Icons, Button } from '../components/ui';

type PopoverImageType = {
    setPopOverSrc: React.Dispatch<React.SetStateAction<string>>;
    setPopOverAlt: React.Dispatch<React.SetStateAction<string>>;
};

export const PopoverImageContext = createContext<PopoverImageType | undefined>(undefined);

export function PopoverImageProvider({ children }: { children: ReactNode }) {
    const [popOverSrc, setPopOverSrc] = useState<string>('/_static/images/lesi.webp');
    const [popOverAlt, setPopOverAlt] = useState<string>('');
    return (
        <PopoverImageContext.Provider value={{ setPopOverSrc, setPopOverAlt }}>
            <div
                popover='auto'
                id='large-image-preview-modal'
                className='peer fixed max-h-full max-w-full rounded bg-base-100 p-2 drop-shadow-2xl'>
                <div className='relative flex h-full w-full justify-end'>
                    <Button
                        variant='accent'
                        size='icon'
                        className='absolute right-4 top-4'
                        style={{ top: '16px' }}
                        onClick={() => {
                            document.getElementById('large-image-preview-modal')?.hidePopover();
                        }}>
                        <Icons.X />
                    </Button>

                    <img src={popOverSrc} width={1920} height={960} alt={popOverAlt} className='!max-h-full w-auto' />
                </div>
            </div>
            {children}
        </PopoverImageContext.Provider>
    );
}

export function usePopoverImage() {
    const context = useContext(PopoverImageContext);
    if (context === undefined) {
        throw new Error('usePopoverImage must be used within a PopoverImageProvider');
    }
    return context;
}
