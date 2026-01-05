import { useRef, useEffect } from 'react';

export function Ltoe() {
    const memeRef = useRef<HTMLImageElement>(null);

    const randomPosition = () => {
        if (!memeRef.current) {
            return;
        }
        const windowHeight = window.innerHeight - memeRef.current.offsetHeight;
        const windowWidth = window.innerWidth - memeRef.current.offsetWidth;
        const randomTop = Math.floor(Math.random() * windowHeight);
        const randomLeft = Math.floor(Math.random() * windowWidth);

        return { top: randomTop, left: randomLeft };
    };

    useEffect(() => {
        if (!memeRef.current) {
            return;
        }
        const memeElement = memeRef.current;
        const position = randomPosition();
        if (!position) return;
        memeElement.style.top = `${position.top}px`;
        memeElement.style.left = `${position.left}px`;

        const handleMouseEnter = () => {
            memeElement.style.transition = 'all 0.2s';
            memeElement.style.opacity = '1';
            memeElement.style.width = '800px';
            memeElement.style.height = '800px';
            setTimeout(() => {
                memeElement.style.transition = 'all 1s';
                memeElement.style.width = '650px';
                memeElement.style.height = '650px';
            }, 200);
        };

        const handleMouseLeave = () => {
            memeElement.style.transition = 'all 0.5s';
            memeElement.style.width = '1px';
            memeElement.style.height = '1px';
        };

        memeElement.addEventListener('mouseenter', handleMouseEnter);
        memeElement.addEventListener('mouseleave', handleMouseLeave);

        return () => {
            memeElement.removeEventListener('mouseenter', handleMouseEnter);
            memeElement.removeEventListener('mouseleave', handleMouseLeave);
        };
    }, []);

    return (
        <a target='_blank' rel='noreferrer' href='https://twitch.tv/L2KxD' tabIndex={-1} className='z-40'>
            <img
                ref={memeRef}
                className='absolute left-0 top-0 z-40 h-[1px] w-[1px]'
                src='/_static/images/ltoe.webp'
                alt='The sacred image is missing, plz forgiv'
                height={800}
                width={800}
            />
        </a>
    );
}
