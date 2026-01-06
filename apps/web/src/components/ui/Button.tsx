import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/utils';

const buttonVariants = cva(
    'inline-flex items-center justify-center whitespace-nowrap rounded-lg text-sm font-medium ring-offset-white transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-neutral-950 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed peer-checked:before:text-white break-words overflow-hidden active:scale-90 transition-transform ',
    {
        variants: {
            variant: {
                default: 'bg-primary hover:bg-primary/80 text-primary-content select-none',
                secondary: 'bg-secondary hover:bg-secondary/80 text-secondary-content select-none',
                accent: 'bg-accent hover:bg-accent/80 text-accent-content select-none',
                destructive: 'bg-error hover:bg-error/70 text-error-content select-none',
                outline:
                    'bg-transparent hover:bg-base-300/80 border-2 border-base-content/80 hover:border-base-content/50 hover:text-base-content select-none',
                ghost: 'hover:bg-base-content/80 hover:text-base-100 select-none text-base-content',
                link: 'text-base-content underline-offset-2 underline hover:text-base-content/80 hover:underline w-fit !p-0',
                plainLink: 'text-base-content hover:text-base-content/80 hover:underline w-fit !p-0',
                gradient: 'bg-gradient-to-r hover:opacity-90 from-primary to-secondary text-base-content  select-none',
            },
            size: {
                default: 'h-10 px-4 py-2',
                sm: 'h-9 px-3',
                xl: 'h-14  px-8',
                icon: 'h-10 w-10',
                link: 'h-fit w-fit p-0 m-0',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    }
);

export interface ButtonProps
    extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {
    asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
    ({ className, variant, size, asChild = false, ...props }, ref) => {
        const Component = asChild ? Slot : 'button';
        return <Component className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />;
    }
);
Button.displayName = 'Button';

export default Button;
