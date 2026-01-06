import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/utils';

const textareaVariants = cva(
    'inline-flex items-center justify-start whitespace-pre-wrap rounded-lg text-sm font-medium ring-offset-white transition-colors focus-visible:!outline-none focus-visible:ring-1 focus-visible:ring-neutral-950 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed select-text placeholder:italic',
    {
        variants: {
            variant: {
                default: 'bg-primary text-primary-content placeholder:text-primary-content/50',
                secondary: 'bg-secondary text-secondary-content placeholder:text-secondary-content/50',
                accent: 'bg-accent text-accent-content placeholder:text-accent-content/50',
                outline:
                    'bg-transparent border-2 border-base-content/80 text-base-content placeholder:text-base-content/50',
                underline:
                    'bg-transparent border border-b-accent border-t-0 border-l-0 border-r-0 rounded-none rounded-t-lg focus-visible:!ring-0',
            },
            size: {
                default: 'min-h-[100px] p-3',
                sm: 'min-h-[80px] p-2',
                xl: 'min-h-[140px] p-4',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    }
);

export interface TextareaProps
    extends Omit<React.TextareaHTMLAttributes<HTMLTextAreaElement>, 'size'>, VariantProps<typeof textareaVariants> {
    id: string;
    error?: string;
    tags?: boolean;
    defaultTags?: string[];
    onTagsChange?: (tags: string[]) => void;
}

const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
    ({ id, className, variant, size, error, tags, defaultTags, onTagsChange, ...props }, ref) => {
        const [currentTag, setCurrentTag] = React.useState('');
        const [tagsList, setTagsList] = React.useState<string[]>(defaultTags || []);

        const commitTag = (tag: string) => {
            const trimmed = tag.trim();
            if (trimmed && !tagsList.includes(trimmed)) {
                const updated = [...tagsList, trimmed];
                setTagsList(updated);
                onTagsChange?.(updated);
            }
            setCurrentTag('');
        };

        const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                commitTag(currentTag);
            } else if (e.key === 'Backspace' && currentTag === '') {
                e.preventDefault();
                const updated = tagsList.slice(0, -1);
                setTagsList(updated);
                onTagsChange?.(updated);
            }
        };

        const handleBlur = () => {
            commitTag(currentTag);
        };

        const removeTag = (index: number) => {
            const updated = tagsList.filter((_, i) => i !== index);
            setTagsList(updated);
            onTagsChange?.(updated);
        };

        if (tags) {
            return (
                <div
                    className={cn(
                        textareaVariants({ variant, size, className }),
                        'flex flex-wrap items-start gap-2 overflow-y-auto px-4 py-2'
                    )}>
                    {tagsList.map((tag, i) => (
                        <span
                            key={i}
                            className='bg-base-300 text-base-content flex items-center gap-1 rounded px-2 py-1 text-xs'>
                            {tag}
                            <button
                                type='button'
                                onClick={() => removeTag(i)}
                                className='text-base-content/60 hover:text-base-content'>
                                ×
                            </button>
                        </span>
                    ))}
                    <textarea
                        id={id}
                        value={currentTag}
                        onChange={(e) => setCurrentTag(e.target.value)}
                        onKeyDown={handleKeyDown}
                        onBlur={handleBlur}
                        className='placeholder:text-base-content/50 min-w-[120px] flex-grow resize-none border-none bg-transparent outline-none'
                        rows={1}
                        {...props}
                    />
                    {error && <span className='text-error w-full pl-1 pt-1 text-xs'>{error}</span>}
                </div>
            );
        }

        return (
            <>
                <textarea
                    id={id}
                    className={cn(textareaVariants({ variant, size, className }), error && '!border-error/60')}
                    ref={ref}
                    {...props}
                />
                {error && <span className='text-error pl-1 pt-1 text-xs'>{error}</span>}
            </>
        );
    }
);

Textarea.displayName = 'Textarea';
export default Textarea;
