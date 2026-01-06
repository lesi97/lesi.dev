import { cn } from '@/utils';

interface DescriptionProps {
    title: string;
    subtitle: string | React.ReactNode;
    className?: string;
}

export function Description({ title, subtitle, className }: DescriptionProps) {
    return (
        <div
            className={cn(
                'border-b-accent mb-5 border border-l-0 border-r-0 border-t-0 pb-4 text-center sm:border-0',
                className
            )}>
            <h1 className='mb-4 text-pretty text-2xl'>{title}</h1>
            <h2 className='text-pretty text-lg'>{subtitle}</h2>
        </div>
    );
}
