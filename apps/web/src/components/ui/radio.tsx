import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/utils';

const radioVariants = cva(
    'inline-flex items-center justify-center whitespace-nowrap rounded-lg text-sm font-medium ring-offset-white transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-neutral-950 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer relative peer-checked:*:flex border-2 select-none group-focus:ring-2 group-focus:ring-neutral/50 group-focus:ring-offset-2 ',
    {
        variants: {
            variant: {
                default:
                    'bg-primary/40 hover:bg-primary/80 peer-checked:bg-primary peer-checked:hover:bg-primary/80 border-primary text-primary-content/60 peer-checked:text-primary-content hover:text-primary-content peer-checked:border-primary peer-checked:hover:border-primary/80 active:scale-95 transition-transform ',
                secondary:
                    'bg-secondary/40 hover:bg-secondary/80 peer-checked:bg-secondary peer-checked:hover:bg-secondary/80 border-secondary peer-checked:border-secondary peer-checked:hover:border-secondary/80 text-secondary-content/60 peer-checked:text-secondary-content hover:text-secondary-content active:scale-95 transition-transform ',
                accent: 'bg-accent/40 hover:bg-accent/80 peer-checked:bg-accent peer-checked:hover:bg-accent/80 border-accent peer-checked:border-accent peer-checked:hover:border-accent/80 text-accent-content/60 peer-checked:text-accent-content hover:text-accent-content active:scale-95 transition-transform',
                outline:
                    'bg-transparent hover:bg-neutral/50 border-neutral/80 hover:border-neutral/50 hover:text-base-content active:scale-95 transition-transform',
                'circle-primary':
                    'peer-checked:*:hidden *:hidden bg-primary/40 hover:bg-primary peer-checked:bg-primary border border-transparent peer-checked:border-primary-content',
                'circle-secondary':
                    'peer-checked:*:hidden *:hidden bg-secondary/40 hover:bg-secondary border border-transparent peer-checked:bg-secondary peer-checked:border-secondary-content',
                'circle-accent':
                    'peer-checked:*:hidden *:hidden bg-accent/40 hover:bg-accent border border-transparent peer-checked:bg-accent peer-checked:border-neutral peer-checked:border-accent-content',
                'circle-outline':
                    'peer-checked:*:hidden *:hidden bg-transparent hover:bg-neutral/50 border-neutral/80 hover:border-neutral/50 hover:text-base-content',
            },
            size: {
                default: 'h-10 px-4 py-2',
                sm: 'w-4 h-4 rounded',
                circle: 'w-6 h-6 rounded-full peer-checked:before:w-full peer-checked:before:content-["."] peer-checked:before:h-full peer-checked:before:flex peer-checked:before:flex-row peer-checked:before:items-center peer-checked:before:justify-center peer-checked:before:absolute peer-checked:before:-top-[29px] peer-checked:before:text-8xl',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    }
);

export interface RadioProps
    extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>, VariantProps<typeof radioVariants> {
    asChild?: boolean;
    label?: React.ReactNode;
}

export const Radio = React.forwardRef<HTMLInputElement, RadioProps>(
    ({ className, variant, size, asChild = false, id, label, ...props }, ref) => {
        const Component = asChild ? Slot : 'input';
        return (
            <label
                htmlFor={id}
                className='group focus-within:outline-0 focus-visible:outline-0 focus-visible:ring-0'
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        (e.target as HTMLLabelElement).click();
                    }
                }}
                tabIndex={0}>
                <Component id={id} type='radio' className='peer hidden' ref={ref} {...props} />
                <div className={cn(radioVariants({ variant, size, className }))}>
                    <span>{label}</span>
                </div>
            </label>
        );
    }
);
Radio.displayName = 'Radio';
