import { Link } from 'react-router-dom';
import { cn } from '@/utils';

export function TwitchMessageContainer({
    children,
    user,
    className,
}: {
    children: React.ReactNode;
    user: string;
    className?: string;
}) {
    const badges = [];
    const isNightbot = user === 'Nightbot';
    const nightbot = { href: 'http://nightbot.tv/', colour: '[#7c7ce1]' };
    const me = { href: 'https://twitch.tv/c_lesi', colour: 'accent' };

    if (isNightbot) {
        badges.push(
            {
                alt: 'Twitch Bot',
                src: '/images/icons/Twitch-Bot.png',
                width: 16,
                height: 16,
                className: 'inline-flex h-[16px] w-[16px]',
            },
            {
                alt: 'Twitch Verified',
                src: '/images/icons/Twitch-Verified.png',
                width: 16,
                height: 16,
                className: 'inline-flex h-[16px] w-[16px]',
            }
        );
        <Link to='http://nightbot.tv/' target='_blank' className='inline-block'>
            <span className='font-bold text-[#7c7ce1]'>Nightbot:</span>
        </Link>;
    } else {
        badges.push(
            {
                alt: 'Twitch Mod',
                src: '/images/icons/Twitch-Mod.png',
                width: 16,
                height: 16,
                className: 'inline-flex h-[16px] w-[16px]',
            },
            {
                alt: 'Twitch Prime',
                src: '/images/icons/Twitch-Prime.png',
                width: 16,
                height: 16,
                className: 'inline-flex h-[16px] w-[16px]',
            }
        );
    }

    return (
        <>
            <div
                className={cn(
                    'flex flex-row place-content-start items-start justify-start gap-x-1 break-all',
                    className
                )}>
                <span className='flex h-full min-w-fit flex-row items-center gap-x-1 break-all'>
                    {badges.map((badge, idx) => {
                        return (
                            <img
                                key={idx}
                                alt={badge.alt}
                                src={badge.src}
                                width={badge.width}
                                height={badge.height}
                                className={badge.className}
                            />
                        );
                    })}
                    <Link to={isNightbot ? nightbot.href : me.href} target='_blank' className='inline-block'>
                        <span className={cn('font-bold', isNightbot ? `text-${nightbot.colour}` : `text-${me.colour}`)}>
                            {user}:
                        </span>
                    </Link>
                </span>
                <span className='inline-block w-fit'>{children}</span>
            </div>
        </>
    );
}
