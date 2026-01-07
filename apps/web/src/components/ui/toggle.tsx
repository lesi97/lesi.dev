import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/utils';
import { Icons } from './icons';

const toggleVariants = cva(
    'inline-flex items-center justify-center whitespace-nowrap rounded-lg text-sm font-medium ring-offset-white transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-neutral-950 disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer relative *:invisible peer-checked:*:visible peer-checked:*:inline-flex border-2 select-none group-focus:ring-2 group-focus:ring-offset-2 ',
    {
        variants: {
            variant: {
                default:
                    'bg-neutral hover:bg-primary/80 peer-checked:bg-primary peer-checked:hover:bg-primary/80 border-primary text-primary-content peer-checked:border-primary peer-checked:hover:border-primary/80 group-focus:ring-neutral/50 ',
                secondary:
                    'bg-neutral hover:bg-secondary/80 peer-checked:bg-secondary peer-checked:hover:bg-secondary/80 border-secondary peer-checked:border-secondary peer-checked:hover:border-secondary/80 text-secondary-content group-focus:ring-neutral/50 ',
                accent: 'bg-neutral hover:bg-accent/80 peer-checked:bg-accent peer-checked:hover:bg-accent/80 border-accent peer-checked:border-accent peer-checked:hover:border-accent/80 text-accent-content group-focus:ring-neutral/50 ',
                outline:
                    'bg-transparent hover:bg-neutral/50 border-neutral/80 hover:border-neutral/50 hover:text-base-content group-focus:ring-neutral/50 ',
            },
            size: {
                default: 'w-8 h-8 rounded',
                sm: 'w-4 h-4 rounded',
                md: 'w-6 h-6 rounded',
                lg: 'w-10 h-10 rounded',
                xl: 'w-12 h-12 rounded',
                '2xl': 'w-14 h-14 rounded',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    }
);

export interface ToggleProps
    extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>, VariantProps<typeof toggleVariants> {
    asChild?: boolean;
}

const Toggle = React.forwardRef<HTMLInputElement, ToggleProps>(
    ({ className, variant, size, asChild = false, id, ...props }, ref) => {
        const Component = asChild ? Slot : 'input';
        return (
            <label
                htmlFor={id}
                className='group h-fit w-fit focus-within:outline-0 focus-visible:outline-0 focus-visible:ring-0'
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        (e.target as HTMLLabelElement).click();
                    }
                }}
                tabIndex={0}>
                <Component id={id} type='checkbox' className='peer hidden' ref={ref} {...props} />
                <div className={cn(toggleVariants({ variant, size, className }), '')}>
                    <Icons.Checkmark className='h-full w-full' />
                </div>
            </label>
        );
    }
);
Toggle.displayName = 'Toggle';

export default Toggle;
