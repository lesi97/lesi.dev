import { Link } from 'react-router-dom';
import type { ReactNode } from 'react';
import { useMemo, useState } from 'react';
import { useSeason } from '@/context/SeasonContext';
import { cn } from '@/utils';
import { Icons } from '../ui';
import { useIsMobile } from '@/hooks';

type NavLink = {
    to: string;
    label: string;
    children?: Array<{ to: string; label: string }>;
};

const links: NavLink[] = [
    { to: '/aspect-ratio-calculator', label: 'Aspect Ratio Calculator' },
    {
        to: '/file-converters',
        label: 'File Converters',
        children: [
            { to: '/pdf-to-png', label: 'PDF To PNG' },
            { to: '/ico-converter', label: 'image to Icon' },
            { to: '/video-to-mp3', label: 'Video To MP3' },
            { to: '/video-cropper', label: 'Video Cropper' },
        ],
    },
    { to: '/password-generator', label: 'Password Generator' },
    {
        to: '/uuid-generator/v4',
        label: 'UUID Generator',
        children: [
            { to: '/uuid-generator/v1', label: 'UUID v1' },
            { to: '/uuid-generator/v4', label: 'UUID v4' },
            { to: '/uuid-generator/v6', label: 'UUID v6' },
            { to: '/uuid-generator/v7', label: 'UUID v7' },
            { to: '/uuid-generator/nil', label: 'Nil UUID' },
        ],
    },
    { to: '/weight-converter', label: 'Weight Converter' },
    { to: '/minifier', label: 'Minifier' },
    { to: '/countdown', label: 'Countdown' },
    { to: '/aim-trainer', label: 'Aim Trainer' },
    { to: '/settings', label: 'Settings' },
];

type SeasonImage = {
    src: string;
    alt: string;
};

const seasonImages: Record<string, { src: string; alt: string }> = {
    Valentines: {
        src: '/images/lesi-valentines.webp',
        alt: 'Lesi-Valentines',
    },
    Halloween: {
        src: '/images/lesi-halloween.webp',
        alt: 'Lesi-Halloween',
    },
    'Lesi-Birthday': {
        src: '/images/lesi-birthday.webp',
        alt: 'Lesi-Birthday',
    },
    Christmas: {
        src: '/images/lesi-christmas.webp',
        alt: 'Lesi-Christmas',
    },
    'New-Years': {
        src: '/images/lesi-newyear.webp',
        alt: 'Lesi-New-Years',
    },
};

type SeasonLogoProps = {
    season: string | null;
};

function SeasonLogo({ season }: SeasonLogoProps) {
    const size = 47;

    const image = season ? seasonImages[season] : undefined;

    const src = image?.src ?? '/images/lesi.webp';
    const alt = image?.alt ?? 'Lesi';

    return <img src={src} alt={alt} height={size} width={size} className='h-full w-auto' />;
}

type NavItemProps = {
    item: NavLink;
    onNavigate?: () => void;
};

function NavItem({ item, onNavigate }: NavItemProps) {
    const hasChildren = Boolean(item.children && item.children.length > 0);
    const isMobile = useIsMobile();
    const [open, setOpen] = useState(false);

    function handleParentClick() {
        if (isMobile && hasChildren) {
            setOpen((v) => !v);
        } else {
            onNavigate?.();
        }
    }

    return (
        <li
            className={cn(
                'group list-none text-primary-content hover:bg-secondary/30',
                isMobile ? 'flex flex-col' : 'flex items-center max-h-47px'
            )}>
            <div className='relative w-full'>
                {hasChildren ? (
                    <button
                        type='button'
                        className={cn(
                            'relative flex w-full flex-row items-center justify-between gap-2 px-4 py-4 text-left text-[17px] text-primary-content',
                            'xl:h-full xl:max-h-[47px] xl:py-0'
                        )}
                        onClick={handleParentClick}
                        aria-expanded={open}>
                        <span>{item.label}</span>
                        <Icons.Chevron
                            className={cn(
                                'h-4 w-4 mr-1 text-accent transition-transform duration-300 ease-in-out',
                                isMobile ? (open ? 'rotate-0' : 'rotate-180') : 'rotate-180 group-hover:rotate-0'
                            )}
                        />
                    </button>
                ) : (
                    <Link
                        to={item.to}
                        onClick={onNavigate}
                        className='relative flex w-full flex-row items-center gap-2 px-4 py-4 text-center text-[17px] text-primary-content no-underline xl:h-full xl:max-h-[47px] xl:py-0'>
                        <span>{item.label}</span>
                    </Link>
                )}

                {hasChildren ? (
                    <div
                        className={cn(
                            isMobile
                                ? cn(
                                      open ? 'block' : 'hidden',
                                      'w-full bg-secondary/20',
                                      'border-l-2 border-accent/40',
                                      'pl-4'
                                  )
                                : 'min-w-parent absolute left-0 top-[35px] hidden w-max min-w-full divide-y divide-solid divide-neutral shadow-xl xl:flex-col xl:group-hover:flex'
                        )}>
                        {item.children?.map((child) => {
                            return (
                                <Link
                                    key={child.to}
                                    to={child.to}
                                    onClick={onNavigate}
                                    className={cn(
                                        'flex h-[47px] w-full items-center px-4 py-4 text-left text-[17px] no-underline',
                                        isMobile
                                            ? 'text-primary-content hover:bg-secondary/20'
                                            : 'bg-secondary text-secondary-content last:rounded-b-md hover:bg-secondary/90 xl:py-0'
                                    )}>
                                    {child.label}
                                </Link>
                            );
                        })}
                    </div>
                ) : null}
            </div>
        </li>
    );
}

type NavListProps = {
    expanded: boolean;
    children?: ReactNode;
};

function NavList({ expanded, children }: NavListProps) {
    return (
        <ul
            className={cn(
                'absolute top-[47px] z-20 m-0 h-fit w-full bg-inherit p-0 xl:relative xl:top-0 xl:!flex xl:h-full xl:max-h-47px xl:flex-row',
                expanded ? 'flex flex-col' : 'hidden'
            )}>
            {children}
        </ul>
    );
}

export function Nav() {
    const { season } = useSeason();
    const [navExpanded, setNavExpanded] = useState(false);

    const closeNav = useMemo(() => {
        return () => setNavExpanded(false);
    }, []);

    return (
        <header className='w-full place-items-start place-self-start'>
            <nav className='relative flex h-47px w-full flex-row justify-between bg-primary xl:justify-start'>
                <Link className='relative h-full w-fit overflow-hidden' to='/' onClick={closeNav}>
                    <SeasonLogo season={season} />
                </Link>

                <button className='px-4 xl:hidden' onClick={() => setNavExpanded((v) => !v)}>
                    <Icons.Burger />
                </button>

                <NavList expanded={navExpanded}>
                    {links.map((item) => {
                        return <NavItem key={item.to} item={item} onNavigate={closeNav} />;
                    })}
                </NavList>
            </nav>
        </header>
    );
}
