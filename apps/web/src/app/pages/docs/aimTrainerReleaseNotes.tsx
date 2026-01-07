'use client';
import { ReactNode, useState } from 'react';
import { Icons, Button } from '@/components/ui';
import { cn } from '@/utils';

export function AimTrainerReleaseNotes() {
    const [isExpanded, setIsExpanded] = useState<string | null>('1.0.1');

    function Summary({
        description,
        isExpanded,
        currentVer,
    }: {
        description: ReactNode;
        isExpanded: boolean;
        currentVer: string;
    }) {
        return (
            <summary
                className='mb-4 grid w-fit cursor-pointer grid-cols-[235px_50px] flex-row items-center gap-4'
                onClick={(e) => {
                    e.preventDefault();
                    setIsExpanded((prev) => (prev === currentVer ? null : currentVer));
                }}>
                <Button variant='plainLink' size='link' className='text-xl'>
                    {description}
                </Button>
                <Icons.Chevron
                    width={12}
                    height={12}
                    className={cn(
                        'transition-transform duration-500 ease-in-out',
                        isExpanded ? 'rotate-0' : 'rotate-180'
                    )}
                />
            </summary>
        );
    }

    return (
        <div data-id-version='1.0.1'>
            <details open={isExpanded === '1.0.1'}>
                <Summary
                    description={
                        <span>
                            March 08<sup>th</sup> 2025 - v1.0.1
                        </span>
                    }
                    currentVer='1.0.1'
                    isExpanded={isExpanded === '1.0.1'}></Summary>
                <div className='mb-6'>
                    <h2 className='mb-2 text-lg'>Bug Fixes</h2>
                    <ul className='ml-5 list-disc space-y-2'>
                        <li>Fixed issue with the Show Weapon toggle not saving the preference between sessions</li>
                        <li>The same issue with the Show FPS toggle has also been fixed</li>
                    </ul>
                </div>

                <h2 className='mb-2 text-lg'>New Features</h2>
                <ul className='ml-5 list-disc space-y-2'>
                    <li>The game timer will now start once the first box has been hit</li>
                    <li>Added R as an additional keybind to restart the game</li>
                </ul>
            </details>
        </div>
    );
}
