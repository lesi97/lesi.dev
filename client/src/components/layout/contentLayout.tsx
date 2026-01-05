import { Outlet } from 'react-router-dom';
import { HexBackground } from './hexBackground';
import { DateTime } from './dateTime';
import { Nav } from './nav';

type ColourStop = { offset: number; colour: string };

type Props = {
    gradient: string;
    stops?: ColourStop[];
    className?: string;
};

export function ContentLayout({ gradient, stops, className }: Props) {
    return (
        <HexBackground gradient={gradient} stops={stops} maxDurationMs={10000}>
            <main
                className={`relative flex h-full w-full flex-col justify-between gap-8 rounded-lg bg-repeat px-8 py-8 ${className ?? ''}`}>
                <DateTime containerClassNames='place-self-end' />
                <Outlet />
                <Nav />
            </main>
        </HexBackground>
    );
}
