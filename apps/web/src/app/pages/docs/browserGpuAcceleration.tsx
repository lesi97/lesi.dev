'use client';
import { Icons, Button } from '@/components/ui';
import { useState } from 'react';
import { cn } from '@/utils';
import { usePopoverImage } from '@/context/PopoverImageContext';

type BrowserList = 'Chrome' | 'Opera' | 'Edge' | 'Firefox';

export function BrowserGpuAcceleration() {
    const [isExpanded, setIsExpanded] = useState<BrowserList | null>(null);
    const { setPopOverSrc, setPopOverAlt } = usePopoverImage();

    return (
        <div className='flex w-full flex-col gap-2'>
            <h1 className='mb-4 w-full text-center text-xl'>How to enable your browser&#39;s GPU Acceleration</h1>
            <details className='w-full transition-all ease-in-out' open={isExpanded === 'Chrome'}>
                <summary
                    className='mb-4 grid w-[200px] cursor-pointer grid-cols-[20px_135px_50px] flex-row items-center gap-4'
                    onClick={(e) => {
                        e.preventDefault();
                        setIsExpanded((prev) => (prev === 'Chrome' ? null : 'Chrome'));
                    }}>
                    <Icons.Browsers.GoogleChrome />
                    <Button variant='plainLink' size='link' className='text-lg'>
                        Google Chrome
                    </Button>
                    <Icons.Chevron
                        width={12}
                        height={12}
                        className={cn(
                            'transition-transform duration-500 ease-in-out',
                            isExpanded === 'Chrome' ? 'rotate-0' : 'rotate-180'
                        )}
                    />
                </summary>
                <span className='flex w-full flex-col gap-4 overflow-hidden transition-all duration-500 ease-in-out'>
                    <p>
                        Within Google Chrome, copy the below Google Chrome URL and paste it in your URL bar to go
                        directly to the required setting.
                    </p>
                    <p>
                        Note: You cannot click on the link as system links are prevented as clickable links<span></span>
                    </p>
                    <p className='flex flex-row items-center gap-4 font-bold text-secondary'>
                        chrome://settings/system
                        <Button
                            variant='accent'
                            size='icon'
                            onClick={() => {
                                navigator.clipboard.writeText('chrome://settings/system');
                            }}>
                            <Icons.Copy width={18} height={18} />
                        </Button>
                    </p>
                    <p>Look for the setting &#34;Use graphics acceleration when available&#34; and enable it</p>
                    <p>Close Google Chrome and re-open it</p>
                    <button
                        className='mb-8 w-full'
                        popoverTarget='large-image-preview-modal'
                        onClick={() => {
                            setPopOverSrc(
                                '/static/images/help/browser-gpu-acceleration/enable-chrome-gpu-acceleration.png'
                            );
                            setPopOverAlt(
                                'Visual representation of the Google Chrome settings page showing the "Use graphics acceleration when available" option'
                            );
                        }}>
                        <img
                            src='/static/images/help/browser-gpu-acceleration/enable-chrome-gpu-acceleration.png'
                            alt='Visual representation of the Google Chrome settings page showing the "Use graphics acceleration when available" option'
                            width={1920}
                            height={1080}
                            className='w-full'
                        />
                    </button>
                </span>
            </details>

            <details className='w-full transition-all ease-in-out' open={isExpanded === 'Firefox'}>
                <summary
                    className='mb-4 grid w-[200px] cursor-pointer grid-cols-[20px_135px_50px] items-center gap-4'
                    onClick={(e) => {
                        e.preventDefault();
                        setIsExpanded((prev) => (prev === 'Firefox' ? null : 'Firefox'));
                    }}>
                    <Icons.Browsers.Firefox />
                    <Button variant='plainLink' size='link' className='text-lg'>
                        Firefox
                    </Button>
                    <Icons.Chevron
                        width={12}
                        height={12}
                        className={cn(
                            'transition-transform duration-500 ease-in-out',
                            isExpanded === 'Firefox' ? 'rotate-0' : 'rotate-180'
                        )}
                    />
                </summary>
                <span className='flex w-full flex-col gap-4'>
                    <p>
                        Within Firefox, copy the below Firefox URL and paste it in your URL bar to go directly to the
                        required setting.
                    </p>
                    <p>
                        Note: You cannot click on the link as system links are prevented as clickable links<span></span>
                    </p>
                    <p className='flex flex-row items-center gap-4 font-bold text-secondary'>
                        about:preferences{' '}
                        <Button
                            variant='accent'
                            size='icon'
                            onClick={() => {
                                navigator.clipboard.writeText('about:preferences');
                            }}>
                            <Icons.Copy width={18} height={18} />
                        </Button>
                    </p>
                    <p>You should be in General</p>
                    <p>Scroll down to Performance and uncheck &#34;Use recommended perfomance settings&#34;</p>
                    <p>You should then see &#34;Use hardware acceleration when available&#34;, ensure this is ticked</p>
                    <p>Close Firefox and re-open it</p>
                    <button
                        className='mb-8 w-full'
                        popoverTarget='large-image-preview-modal'
                        onClick={() => {
                            setPopOverSrc(
                                '/static/images/help/browser-gpu-acceleration/enable-firefox-gpu-acceleration.png'
                            );
                            setPopOverAlt(
                                'Visual representation of the Firefox settings page showing the "Use hardware acceleration when available" option'
                            );
                        }}>
                        <img
                            src='/static/images/help/browser-gpu-acceleration/enable-firefox-gpu-acceleration.png'
                            alt='Visual representation of the Firefox settings page showing the "Use hardware acceleration when available" option'
                            width={1920}
                            height={1080}
                            className='w-full'
                        />
                    </button>
                </span>
            </details>

            <details className='w-full transition-all ease-in-out' open={isExpanded === 'Edge'}>
                <summary
                    className='mb-4 grid w-[200px] cursor-pointer grid-cols-[20px_135px_50px] items-center gap-4'
                    onClick={(e) => {
                        e.preventDefault();
                        setIsExpanded((prev) => (prev === 'Edge' ? null : 'Edge'));
                    }}>
                    <Icons.Browsers.MicrosoftEdge />
                    <Button variant='plainLink' size='link' className='text-lg'>
                        Microsoft Edge
                    </Button>
                    <Icons.Chevron
                        width={12}
                        height={12}
                        className={cn(
                            'transition-transform duration-500 ease-in-out',
                            isExpanded === 'Edge' ? 'rotate-0' : 'rotate-180'
                        )}
                    />
                </summary>
                <span className='flex w-full flex-col gap-4'>
                    <p>
                        Within Microsoft Edge, copy the below Microsoft Edge URL and paste it in your URL bar to go
                        directly to the required setting.
                    </p>
                    <p>
                        Note: You cannot click on the link as system links are prevented as clickable links<span></span>
                    </p>
                    <p className='flex flex-row items-center gap-4 font-bold text-secondary'>
                        edge://settings/system{' '}
                        <Button
                            variant='accent'
                            size='icon'
                            onClick={() => {
                                navigator.clipboard.writeText('edge://settings/system');
                            }}>
                            <Icons.Copy width={18} height={18} />
                        </Button>
                    </p>
                    <p>Look for the setting &#34;Use graphics acceleration when available&#34; and enable it</p>
                    <p>Close Microsoft Edge and re-open it</p>
                    <button
                        className='mb-8 w-full'
                        popoverTarget='large-image-preview-modal'
                        onClick={() => {
                            setPopOverSrc(
                                '/static/images/help/browser-gpu-acceleration/enable-edge-gpu-acceleration.png'
                            );
                            setPopOverAlt(
                                'Visual representation of the Microsoft Edge settings page showing the "Use hardware acceleration when available" option'
                            );
                        }}>
                        <img
                            src='/static/images/help/browser-gpu-acceleration/enable-edge-gpu-acceleration.png'
                            alt='Visual representation of the Microsoft Edge settings page showing the "Use hardware acceleration when available" option'
                            width={1920}
                            height={1080}
                            className='w-full'
                        />
                    </button>
                </span>
            </details>

            <details className='w-full transition-all ease-in-out' open={isExpanded === 'Opera'}>
                <summary
                    className='mb-4 grid w-[200px] cursor-pointer grid-cols-[20px_135px_50px] items-center gap-4'
                    onClick={(e) => {
                        e.preventDefault();
                        setIsExpanded((prev) => (prev === 'Opera' ? null : 'Opera'));
                    }}>
                    <Icons.Browsers.OperaGX />
                    <Button variant='plainLink' size='link' className='text-lg'>
                        Opera GX
                    </Button>
                    <Icons.Chevron
                        width={12}
                        height={12}
                        className={cn(
                            'transition-transform duration-500 ease-in-out',
                            isExpanded === 'Opera' ? 'rotate-0' : 'rotate-180'
                        )}
                    />
                </summary>
                <span className='flex w-full flex-col gap-4'>
                    <p>
                        Within Opera GX (Or regular Opera), copy the below Opera GX URL and paste it in your URL bar to
                        go directly to the required setting.
                    </p>
                    <p>
                        Note: You cannot click on the link as system links are prevented as clickable links<span></span>
                    </p>
                    <p className='flex flex-row items-center gap-4 text-wrap break-all font-bold text-secondary'>
                        opera://settings/startPage?search=Use+graphics+acceleration+when+available{' '}
                        <Button
                            variant='accent'
                            size='icon'
                            className='min-h-fit min-w-fit px-3 py-2'
                            onClick={() => {
                                navigator.clipboard.writeText(
                                    'opera://settings/startPage?search=Use+graphics+acceleration+when+available'
                                );
                            }}>
                            <Icons.Copy width={18} height={18} />
                        </Button>
                    </p>
                    <p>Look for the setting &#34;Use graphics acceleration when available&#34; and enable it</p>
                    <p>Close Opera GX and re-open it</p>
                    <button
                        className='mb-8 w-full'
                        popoverTarget='large-image-preview-modal'
                        onClick={() => {
                            setPopOverSrc(
                                '/static/images/help/browser-gpu-acceleration/enable-opera-gpu-acceleration.png'
                            );
                            setPopOverAlt(
                                'Visual representation of the Opera GX settings page showing the "Use hardware acceleration when available" option'
                            );
                        }}>
                        <img
                            src='/static/images/help/browser-gpu-acceleration/enable-opera-gpu-acceleration.png'
                            alt='Visual representation of the Opera GX settings page showing the "Use hardware acceleration when available" option'
                            width={1920}
                            height={1080}
                            className='w-full'
                        />
                    </button>
                </span>
            </details>
            <div className='w-full text-right text-sm'>
                Last updated: March 07<sup>th</sup> 2025
            </div>
        </div>
    );
}
