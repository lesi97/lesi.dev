import { useEffect, useMemo, useRef, useState } from 'react';
import { Button, Description, Icons, Input } from '@/components/ui';
import { usePageMeta } from '@/hooks';
import { cn } from '@/utils';

type TimestampStyle = {
    format: 'F' | 'f' | 'D' | 'd' | 't' | 'T' | 'R' | 's' | 'S';
    label: string;
    formatPreview: (date: Date, now: Date) => string;
};

const dateTimeFormat = (options: Intl.DateTimeFormatOptions) => {
    return new Intl.DateTimeFormat(undefined, options);
};

const relativeTimeFormat = new Intl.RelativeTimeFormat(undefined, {
    numeric: 'always',
});

const timestampStyles: TimestampStyle[] = [
    {
        format: 'F',
        label: 'Full date and time',
        formatPreview: (date) => dateTimeFormat({ dateStyle: 'full', timeStyle: 'short' }).format(date),
    },
    {
        format: 'f',
        label: 'Long date and time',
        formatPreview: (date) => dateTimeFormat({ dateStyle: 'long', timeStyle: 'short' }).format(date),
    },
    {
        format: 'D',
        label: 'Long date',
        formatPreview: (date) => dateTimeFormat({ dateStyle: 'long' }).format(date),
    },
    {
        format: 'd',
        label: 'Short date',
        formatPreview: (date) => dateTimeFormat({ dateStyle: 'short' }).format(date),
    },
    {
        format: 't',
        label: 'Short time',
        formatPreview: (date) => dateTimeFormat({ timeStyle: 'short' }).format(date),
    },
    {
        format: 'T',
        label: 'Medium time',
        formatPreview: (date) => dateTimeFormat({ timeStyle: 'medium' }).format(date),
    },
    {
        format: 'R',
        label: 'Relative time',
        formatPreview: (date, now) => formatRelativeTime(date, now),
    },
    {
        format: 's',
        label: 'Short date and time',
        formatPreview: (date) => dateTimeFormat({ dateStyle: 'short', timeStyle: 'short' }).format(date),
    },
    {
        format: 'S',
        label: 'Short date and medium time',
        formatPreview: (date) => dateTimeFormat({ dateStyle: 'short', timeStyle: 'medium' }).format(date),
    },
];

function formatDateTimeLocalValue(date: Date) {
    const localDate = new Date(date.getTime() - date.getTimezoneOffset() * 60_000);

    return localDate.toISOString().slice(0, 19);
}

function formatRelativeTime(date: Date, now: Date) {
    const diffSeconds = Math.round((date.getTime() - now.getTime()) / 1000);
    const absoluteDiffSeconds = Math.abs(diffSeconds);

    if (absoluteDiffSeconds < 1) {
        return '0 seconds ago';
    }

    const divisions: Array<{ amount: number; unit: Intl.RelativeTimeFormatUnit }> = [
        { amount: 60, unit: 'second' },
        { amount: 60, unit: 'minute' },
        { amount: 24, unit: 'hour' },
        { amount: 30, unit: 'day' },
        { amount: 12, unit: 'month' },
        { amount: Number.POSITIVE_INFINITY, unit: 'year' },
    ];

    let duration = diffSeconds;

    for (const division of divisions) {
        if (Math.abs(duration) < division.amount) {
            return relativeTimeFormat.format(Math.round(duration), division.unit);
        }

        duration /= division.amount;
    }

    return relativeTimeFormat.format(Math.round(duration), 'year');
}

async function copyText(text: string) {
    if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(text);
        return;
    }

    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.setAttribute('readonly', 'true');
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
}

export function DiscordTimestampConverter() {
    const [dateTimeValue, setDateTimeValue] = useState(() => formatDateTimeLocalValue(new Date()));
    const [now, setNow] = useState(() => new Date());
    const [copiedFormat, setCopiedFormat] = useState<TimestampStyle['format'] | null>(null);
    const codeRefs = useRef<Partial<Record<TimestampStyle['format'], HTMLElement | null>>>({});

    usePageMeta({
        title: 'Discord Timestamp Converter | Lesi',
        description: 'Create Discord timestamp markdown for each viewer local timezone',
    });

    useEffect(() => {
        const timerId = window.setInterval(() => setNow(new Date()), 1000);

        return () => window.clearInterval(timerId);
    }, []);

    const selectedDate = useMemo(() => {
        if (!dateTimeValue) {
            return null;
        }

        const date = new Date(dateTimeValue);

        return Number.isNaN(date.getTime()) ? null : date;
    }, [dateTimeValue]);

    const unixTime = selectedDate ? Math.floor(selectedDate.getTime() / 1000) : null;

    function highlightTimestampCode(format: TimestampStyle['format']) {
        const codeElement = codeRefs.current[format];

        if (!codeElement) {
            return;
        }

        const selection = window.getSelection();
        const range = document.createRange();

        range.selectNodeContents(codeElement);
        selection?.removeAllRanges();
        selection?.addRange(range);
    }

    async function copyTimestamp(format: TimestampStyle['format']) {
        if (unixTime === null) {
            return;
        }

        await copyText(`<t:${unixTime}:${format}>`);
        highlightTimestampCode(format);
        setCopiedFormat(format);
        window.setTimeout(() => setCopiedFormat(null), 1400);
    }

    function useCurrentTime() {
        const currentDate = new Date();

        setNow(currentDate);
        setDateTimeValue(formatDateTimeLocalValue(currentDate));
    }

    return (
        <>
            <Description
                title='Discord Timestamp Converter'
                subtitle='Create Discord markdown that displays times in each users local timezone'
            />

            <div className='mx-auto flex w-full flex-col gap-6'>
                <div className='flex flex-col sm:flex-row gap-4 justify-center sm:items-end'>
                    <label htmlFor='discord-timestamp-date' className='flex flex-col sm:flex-row gap-4 sm:items-center'>
                        Date and time:
                        <Input
                            id='discord-timestamp-date'
                            type='datetime-local'
                            step='1'
                            value={dateTimeValue}
                            variant='outline'
                            className='w-full sm:max-w-fit justify-start text-left text-base text-base-content '
                            onChange={(event) => setDateTimeValue(event.target.value)}
                        />
                    </label>

                    <Button type='button' variant='secondary' onClick={useCurrentTime}>
                        <Icons.Refresh className='mr-2 h-4 w-4' />
                        Now
                    </Button>
                </div>

                <section className='grid grid-cols-1 gap-2 md:grid-cols-2'>
                    {timestampStyles.map((style) => {
                        const code = unixTime === null ? '' : `<t:${unixTime}:${style.format}>`;
                        const isCopied = copiedFormat === style.format;
                        const CopyIcon = isCopied ? Icons.Check : Icons.Copy;

                        return (
                            <div
                                key={style.format}
                                className='grid min-h-[92px] cursor-pointer grid-cols-[minmax(0,1fr)_auto] grid-rows-[auto_1fr] gap-x-3 gap-y-2 rounded-lg border border-base-content/10 bg-base-200/70 p-3 transition-colors hover:bg-base-200 md:min-h-16 md:grid-cols-[minmax(0,145px)_minmax(0,1fr)_auto] md:grid-rows-1 md:items-center'
                                onClick={() => copyTimestamp(style.format)}>
                                <code
                                    ref={(element) => {
                                        codeRefs.current[style.format] = element;
                                    }}
                                    className={cn(
                                        'col-start-1 row-start-1 w-fit max-w-full overflow-x-auto rounded bg-base-300 px-2 py-1 font-mono text-sm  text-base-content md:col-auto md:row-auto',
                                        unixTime === null && 'text-base-content/40'
                                    )}>
                                    {code || '<t:UNIX:FORMAT>'}
                                </code>

                                <div className='col-start-1 row-start-2 min-w-0 self-start md:col-auto md:row-auto md:self-auto'>
                                    <p className='break-words text-base font-semibold leading-snug text-base-content md:text-sm'>
                                        {selectedDate ? style.formatPreview(selectedDate, now) : style.label}
                                    </p>
                                </div>

                                <button
                                    type='button'
                                    className='col-start-2 row-span-2 row-start-1 flex h-10 w-10 items-center justify-center self-center rounded-lg text-accent hover:bg-accent/10 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-accent disabled:cursor-not-allowed disabled:opacity-40 md:col-auto md:row-span-1'
                                    title={isCopied ? 'Copied' : `Copy ${style.label.toLowerCase()}`}
                                    aria-label={isCopied ? 'Copied timestamp' : `Copy ${style.label.toLowerCase()}`}
                                    disabled={unixTime === null}
                                    onClick={(event) => {
                                        event.stopPropagation();
                                        copyTimestamp(style.format);
                                    }}>
                                    <CopyIcon className='h-5 w-5' />
                                </button>
                            </div>
                        );
                    })}
                </section>
            </div>
        </>
    );
}
