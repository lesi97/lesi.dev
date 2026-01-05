import { Icons } from '@/components/ui';
import { ReactNode } from 'react';

export function Tools() {
    return (
        <div className='flex flex-row flex-1 justify-between gap-8'>
            <Card
                title='ASPECT RATIO CALCULATOR'
                icon={<Icons.GoldenRatio strokeWidth={4} className='text-accent h-full' />}>
                <>
                    <input className='w-full' />
                </>
            </Card>

            <Card title='PASSWORD GENERATOR' icon={<Icons.AllSeeingEye className='text-accent h-full' />}>
                <></>
            </Card>

            <Card title={'WEIGHT\nCONVERTER'} icon={<Icons.Scales className='h-full' />}>
                <></>
            </Card>
        </div>
    );
}

function Card({ title, icon, children }: { title: string; icon: ReactNode; children: ReactNode }) {
    return (
        <div data-id-name='Card' className='border-2 border-primary rounded-lg w-full bg-base-300 z-20'>
            <div
                data-id-name='Card Title'
                className='w-full border-b-2 overflow-hidden border-primary px-8 py-2 flex gap-4 h-1/5 items-center'>
                <span className='text-4xl w-3/4 text-primary whitespace-pre-line leading-none tracking-widest'>
                    {title}
                </span>
                <div className='w-1/4 h-full py-4 flex flex-row items-center justify-end'>{icon}</div>
            </div>

            <div className='w-full p-4'>{children}</div>
        </div>
    );
}
