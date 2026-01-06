import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/utils';
import { Icons } from './icons';
import { useState, useRef, useEffect, Fragment } from 'react';

const sliderContainerVariants = cva(
    'inline-flex rounded-full text-sm font-medium ring-offset-white transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-neutral-950 focus-visible:ring-offset-2 relative disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed select-none',
    {
        variants: {
            variant: {
                default:
                    'border-secondary bg-secondary text-neutral-50 hover:bg-secondary/80 *:bg-primary *:*:bg-secondary',
                secondary:
                    'border-primary  bg-secondary text-neutral-50 hover:bg-secondary/80 *:bg-secondary *:*:bg-primary',
            },
            size: {
                default: 'w-full h-2',
                sm: 'w-full h-1 *:*:-top-1.5',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    }
);

export interface SliderProps
    extends
        Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size' | 'onChange'>,
        VariantProps<typeof sliderContainerVariants> {
    asChild?: boolean;
    onChange?: (value: number) => void;
}

const Slider = ({
    className,
    variant,
    size,
    asChild = false,
    id,
    onChange,
    min = 8,
    max = 100,
    ...props
}: SliderProps) => {
    const [value, setValue] = useState(16);
    const sliderContainerRef = useRef(null);
    const [isDragging, setIsDragging] = useState(false);

    const handleDrag = (e: MouseEvent | Touch) => {
        if (!sliderContainerRef.current) {
            return;
        }
        if (!isDragging) {
            setIsDragging(true);
        }
        const rect = (sliderContainerRef.current as HTMLDivElement).getBoundingClientRect();
        const newValue = ((e.clientX - rect.left) / rect.width) * (max as number);
        if (newValue < (min as number) || newValue > (max as number)) {
            return;
        }
        const roundedValue = Math.round(Math.min(max as number, Math.max(0, newValue)));
        setValue(roundedValue);
        if (onChange) {
            onChange(roundedValue);
        }
    };

    const handleMouseUp = () => setIsDragging(false);

    function handleKeyDown(e: React.KeyboardEvent<HTMLDivElement>) {
        if (e.key === 'ArrowLeft' && parseInt(value.toString()) !== parseInt(min.toString())) {
            setValue(value - 1);
            if (onChange) {
                onChange(value - 1);
            }
        }
        if (e.key === 'ArrowRight' && parseInt(value.toString()) !== parseInt(max.toString())) {
            setValue(value + 1);
            if (onChange) {
                onChange(value + 1);
            }
        }
    }

    useEffect(() => {
        if (isDragging) {
            window.addEventListener('mousemove', handleDrag);
            window.addEventListener('mouseup', handleMouseUp);
        } else {
            window.removeEventListener('mousemove', handleDrag);
            window.removeEventListener('mouseup', handleMouseUp);
        }
        return () => {
            window.removeEventListener('mousemove', handleDrag);
            window.removeEventListener('mouseup', handleMouseUp);
        };
    }, [isDragging]);

    const Component = asChild ? Slot : 'input';
    return (
        <Fragment>
            <div className={cn(sliderContainerVariants({ variant, size, className }))}>
                <Component type='range' className='peer hidden' value={value} disabled {...props} />
                <div
                    className='relative w-full rounded-full'
                    onMouseDown={(e) => {
                        (e.target as HTMLDivElement).focus();
                        handleDrag(e.nativeEvent);
                    }}
                    onMouseMove={(e) => e.buttons === 1 && handleDrag(e.nativeEvent)}
                    onTouchStart={(e) => {
                        (e.target as HTMLDivElement).focus();
                        handleDrag(e.touches[0] as Touch);
                    }}
                    onTouchMove={(e) => e.touches[0].clientX && handleDrag(e.touches[0] as Touch)}
                    ref={sliderContainerRef}>
                    <div
                        className='absolute -top-1 h-4 w-4 cursor-pointer rounded-full bg-secondary'
                        style={{ left: `${(value / (max as number)) * 100}%`, transform: 'translateX(-50%)' }}
                        tabIndex={0}
                        onMouseDown={(e) => (e.target as HTMLDivElement).focus()}
                        onKeyDown={handleKeyDown}></div>
                </div>
            </div>
        </Fragment>
    );
};

Slider.displayName = 'Slider';

export default Slider;
