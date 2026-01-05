import { mergeClassNames } from '@/utils';
import { Link } from 'react-router-dom';
import { useState } from 'react';
import { Triangle } from './triangle';

export function HeroShape() {
    const [hover, setHover] = useState<string | null>(null);
    const links = [
        { label: 'TOOLS', href: '/tools', className: 'mr-[10px]' },
        { label: 'AIM TRAINER', href: '/aim-trainer', className: 'mr-[32px]' },
        { label: 'IMAGE CONVERTERS', href: '/image', className: 'mr-[57px]' },
        { label: 'VIDEO CONVERTERS', href: '/video', className: 'mr-[80px]' },
        { label: 'TWITCH BOT', href: '/twitch', className: 'mr-[102px]' },
    ];

    const rhombus = [
        { label: 'MELCHIOR-1', className: '-top-[128px] left-[372px] rotate-[240deg]', blockFocus: true },
        { label: 'BALTHASAR-2', className: 'top-[103px]', blockFocus: false },
        { label: 'CASPER-3', className: 'rotate-[120deg] -left-[870px] top-[232px]', blockFocus: true },
    ];
    return (
        <div className='flex flex-row gap-32 h-fit w-fit items-center justify-center pl-[285px] pb-[140px] z-20'>
            {rhombus.map((rh, idx) => (
                <>
                    <RhombusContainer
                        key={rh.label}
                        label={rh.label}
                        className={rh.className}
                        links={links}
                        hover={hover}
                        setHover={(value) => setHover(value)}
                        blockFocus={rh.blockFocus}
                    />
                    {idx === 0 && <Triangle />}
                </>
            ))}
        </div>
    );
}

function RhombusContainer({
    label,
    height = '203px',
    width = '400px',
    className = '',
    links,
    hover,
    setHover,
    blockFocus,
}: {
    label: string;
    links: { label: string; href: string; className: string }[];
    height?: string;
    width?: string;
    className?: string;
    hover: string | null;
    setHover: (value: string | null) => void;
    blockFocus: boolean;
}) {
    return (
        <div
            className={mergeClassNames(
                'bg-primary aspect-[5/3] [clip-path:polygon(29%_0,100%_0,70%_100%,0_100%)] relative',
                className
            )}
            style={{ width, height }}>
            <RhombusLabel label={label} />
            {links.map((link) => (
                <RhombusLink
                    label={link.label}
                    href={link.href}
                    className={link.className}
                    isHovered={hover === link.href}
                    onMouseEnter={() => setHover(link.href)}
                    onMouseLeave={() => setHover(null)}
                    blockFocus={blockFocus}
                />
            ))}
        </div>
    );
}

function RhombusLink({
    label,
    href,
    height = '40px',
    width = '240px',
    className = '',
    isHovered,
    onMouseEnter,
    onMouseLeave,
    blockFocus,
}: {
    label: string;
    href: string;
    height?: string;
    width?: string;
    className?: string;
    isHovered: boolean;
    onMouseEnter: () => void;
    onMouseLeave: () => void;
    blockFocus: boolean;
}) {
    return (
        <div
            className={mergeClassNames(
                'place-self-end mr-10 mt-[0px] bg-base-300 text-primary [clip-path:polygon(7.5%_10%,100%_10%,92%_90%,0%_90%)] relative flex items-center text-center',
                className
            )}
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
            onFocus={onMouseEnter}
            onBlur={onMouseLeave}
            style={{ width, height }}>
            <Link
                to={href}
                className={mergeClassNames(
                    'w-full flex items-center text-center place-content-center focus:outline-none',
                    isHovered && 'text-accent'
                )}
                tabIndex={blockFocus ? -1 : 0}>
                {label}
            </Link>
        </div>
    );
}

function RhombusLabel({ label }: { label: string }) {
    return (
        <div className='absolute -left-[32px] place-content-center flex items-center text-center top-[84px] min-w-fit w-[240px] [clip-path:polygon(0%_0,92.3%_0,100%_100%,7.8%_100%)] bg-base-300 py-1 text-primary rotate-[120deg]'>
            {label}
        </div>
    );
}
