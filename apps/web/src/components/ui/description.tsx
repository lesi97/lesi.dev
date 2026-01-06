import { cn } from '@/utils';
import { ReactNode, type HTMLAttributes } from 'react';

type DescriptionProps = HTMLAttributes<HTMLDivElement> & { title: string; subtitle: string | ReactNode };

export function Description({ title, subtitle, className, ...props }: DescriptionProps) {
    return (
        <div
            className={cn(
                'mb-5 border border-l-0 border-r-0 border-t-0 border-b-accent pb-4 text-center sm:border-0',
                className
            )}
            {...props}>
            <h1 className='mb-4 text-pretty text-2xl'>{title}</h1>
            <h2 className='text-pretty text-lg'>{subtitle}</h2>
        </div>
    );
}
