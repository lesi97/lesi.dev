import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';
import { mergeClassNames } from '@/utils';
import { Icons } from './Icons';

const inputVariants = cva(
    'inline-flex items-center justify-center whitespace-nowrap rounded-lg text-sm font-medium ring-offset-white transition-colors focus-visible:!outline-none focus-visible:ring-1 focus-visible:ring-neutral-950 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed select-none placeholder:italic',
    {
        variants: {
            variant: {
                default: 'bg-primary hover:bg-primary/80 text-primary-content placeholder:text-primary-content:50',
                secondary:
                    'bg-secondary hover:bg-secondary/80 text-secondary-content  placeholder:text-secondary-content:50',
                accent: 'bg-accent hover:bg-accent/80 text-accent-content placeholder:text-accent-content:50',
                outline:
                    'bg-transparent hover:bg-accent/10 border-2 border-base-content/80 placeholder:text-base-content/50 text-base-content',
                underline:
                    'bg-transparent hover:bg-accent/10 border border-b-accent border-t-0 border-l-0 border-r-0 rounded-none rounded-t-lg focus-within:!ring-0 focus-within:!outline-0 focus focus-visible:!ring-0 focus-visible:!outline-0 focus:!ring-0 focus:!outline-0 focus-visible:ring-offset-0 focus-visible:bg-accent/10',
            },
            size: {
                default: 'h-10 px-4 py-2',
                sm: 'h-9 px-4',
                xl: 'h-14  px-4',
                icon: 'h-10 w-10',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    }
);

export interface InputProps
    extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>,
        VariantProps<typeof inputVariants> {
    asChild?: boolean;
    error?: string | undefined | null;
    id: string;
    tags?: boolean;
    onTagsChange?: (tags: string[]) => void;
    passwordVisible?: boolean;
    passwordVisibilityHandler?: React.MouseEventHandler<HTMLButtonElement>;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
    (
        {
            className,
            variant,
            size,
            asChild = false,
            error,
            id,
            tags = false,
            onTagsChange,
            passwordVisible,
            passwordVisibilityHandler,
            ...props
        },
        ref
    ) => {
        const Component = asChild ? Slot : 'input';
        const [tagInput, setTagInput] = React.useState('');
        const [tagsList, setTagsList] = React.useState<string[]>([]);

        const handleTagKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
            if (e.key === 'Enter' && tagInput.trim()) {
                e.preventDefault();
                if (!tagsList.includes(tagInput.trim())) {
                    const updated = [...tagsList, tagInput.trim()];
                    setTagsList(updated);
                    onTagsChange?.(updated);
                }
                setTagInput('');
            } else if (e.key === 'Backspace' && tagInput === '') {
                e.preventDefault();
                const updated = tagsList.slice(0, -1);
                setTagsList(updated);
                onTagsChange?.(updated);
            }
        };

        const removeTag = (index: number) => {
            const updated = tagsList.filter((_, i) => i !== index);
            setTagsList(updated);
            onTagsChange?.(updated);
        };

        if (tags) {
            return (
                <div
                    className={mergeClassNames(
                        inputVariants({ variant, size, className }),
                        'flex flex-wrap items-center gap-2 px-4 py-1',
                        error && '!border-error/60'
                    )}>
                    {tagsList.map((tag, index) => (
                        <span
                            key={index}
                            className='flex items-center gap-1 rounded bg-base-300 px-2 py-1 text-xs text-base-content'>
                            {tag}
                            <button
                                type='button'
                                onClick={() => removeTag(index)}
                                className='text-base-content/60 hover:text-base-content'>
                                ×
                            </button>
                        </span>
                    ))}
                    <input
                        type='text'
                        id={id}
                        ref={ref}
                        {...props}
                        value={tagInput}
                        onChange={(e) => setTagInput(e.target.value)}
                        onKeyDown={handleTagKeyDown}
                        className='flex-grow border-none bg-transparent outline-none placeholder:text-base-content/50'
                    />
                    {error && <span className='w-full pl-1 pt-1 text-xs text-error'>{error}</span>}
                </div>
            );
        }

        function handleVisibilityToggle(e: React.MouseEvent<HTMLButtonElement>) {
            e.preventDefault();
            passwordVisibilityHandler?.(e);
            const inputElement = document.getElementById(id) as HTMLInputElement;
            const inputLength = inputElement?.value.length ?? 0;
            inputElement?.focus();
            setTimeout(() => {
                inputElement?.setSelectionRange(inputLength, inputLength);
            }, 0);
        }

        return (
            <>
                <Component
                    className={mergeClassNames(
                        inputVariants({ variant, size, className }),
                        error && '!border-error/60'
                    )}
                    id={id}
                    ref={ref}
                    {...props}
                />
                {passwordVisible === true ? (
                    <button
                        type='button'
                        className='absolute right-0 top-1/2 -translate-x-[8px] -translate-y-[5px] rounded bg-inherit p-2 hover:bg-base-200'
                        onClick={handleVisibilityToggle}>
                        <Icons.EyeClosed width={18} height={18} />
                    </button>
                ) : passwordVisible === false ? (
                    <button
                        type='button'
                        className='absolute right-0 top-1/2 -translate-x-[8px] -translate-y-[5px] rounded bg-inherit p-2 duration-500 ease-in-out hover:bg-base-300'
                        onClick={handleVisibilityToggle}>
                        <Icons.EyeOpen width={18} height={18} />
                    </button>
                ) : null}
                {error && <span className='pl-1 pt-1 text-xs text-error'>{error}</span>}
            </>
        );
    }
);

Input.displayName = 'Input';
export default Input;
