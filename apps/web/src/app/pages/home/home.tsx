import { HeroShape, ArcedText, SinOverlay, Rings } from '.';

export function Home() {
    return (
        <>
            <div className='h-full w-full items-center justify-center overflow-hidden flex relative'>
                <HeroShape />
                <ArcedText />
                <Rings />
                <SinOverlay />
            </div>
        </>
    );
}
