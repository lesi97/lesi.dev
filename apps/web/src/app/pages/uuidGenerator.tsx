import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { Button, Description, Icons, Input, TextArea } from '@/components/ui';
import { cn, downloadFile } from '@/utils';
import { generateUuids, type UuidVersion } from '@/utils/uuid';
import { usePageMeta } from '@/hooks';

const UUID_ROUTE_VERSIONS = ['1', '4', '6', '7', 'nil'] as const satisfies readonly UuidVersion[];

type UuidRouteVersion = (typeof UUID_ROUTE_VERSIONS)[number];

const VERSION_DETAILS: Record<UuidRouteVersion, { title: string; subtitle: string }> = {
    '1': {
        title: 'UUID v1 Generator',
        subtitle: 'Generate timestamp based UUIDs',
    },
    '4': {
        title: 'UUID v4 Generator',
        subtitle: 'Generate random UUIDs',
    },
    '6': {
        title: 'UUID v6 Generator',
        subtitle: 'Generate reordered timestamp UUIDs',
    },
    '7': {
        title: 'UUID v7 Generator',
        subtitle: 'Generate unix time UUIDs',
    },
    nil: {
        title: 'Nil UUID',
        subtitle: 'Generate a nil UUID',
    },
};

export function UuidGenerator({ version }: { version: UuidRouteVersion }) {
    return <UuidGeneratorContent version={version} />;
}

function UuidGeneratorContent({ version }: { version: UuidRouteVersion }) {
    const [uuids, setUuids] = useState<string[]>([]);
    const [quantity, setQuantity] = useState(1);
    const [error, setError] = useState('');
    const [copiedTarget, setCopiedTarget] = useState<'single' | 'single-json' | 'bulk' | 'bulk-json' | null>(null);
    const supportsBulk = version !== 'nil';
    const canGenerateSingle = version !== 'nil';
    const details = VERSION_DETAILS[version];
    const currentUuid = uuids[0] ?? '';
    const bulkOutput = useMemo(() => uuids.join('\n'), [uuids]);
    const copyUuidLabel = copiedTarget === 'single' ? 'Copied UUID' : 'Copy UUID';
    const copyJsonLabel = copiedTarget === 'single-json' ? 'Copied JSON' : 'Copy JSON';

    usePageMeta({
        title: `${details.title} | Lesi`,
        description: details.subtitle,
    });

    useEffect(() => {
        void generate(1);
    }, [version]);

    async function generate(count: number) {
        setError('');

        try {
            const nextUuids = await generateUuids(version, count);
            setUuids(nextUuids);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Unable to generate UUIDs.');
        }
    }

    async function copyText(text: string, target: NonNullable<typeof copiedTarget>) {
        if (!text) {
            return;
        }

        if (navigator.clipboard) {
            await navigator.clipboard.writeText(text);
        } else {
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

        setCopiedTarget(target);
        window.setTimeout(() => setCopiedTarget(null), 1400);
    }

    function downloadUuids() {
        if (!bulkOutput) {
            return;
        }

        const blob = new Blob([bulkOutput], { type: 'text/plain;charset=utf-8' });
        const url = URL.createObjectURL(blob);
        downloadFile(`uuid-v${version}.txt`, url);
        URL.revokeObjectURL(url);
    }

    function updateQuantity(value: string) {
        const nextQuantity = Math.max(1, Math.min(500, Number(value) || 1));
        setQuantity(nextQuantity);
    }

    function uuidJson(values: string[]) {
        return JSON.stringify(values, null, 2);
    }

    const singleUuidActions = [
        {
            key: 'copy-uuid',
            title: copyUuidLabel,
            ariaLabel: copyUuidLabel,
            icon: copiedTarget === 'single' ? Icons.Check : Icons.Copy,
            positionClass: canGenerateSingle ? 'sm:right-20' : 'sm:right-10',
            onClick: () => copyText(currentUuid, 'single'),
        },
        {
            key: 'copy-json',
            title: copyJsonLabel,
            ariaLabel: copyJsonLabel,
            icon: copiedTarget === 'single-json' ? Icons.Check : Icons.CopyJson,
            positionClass: canGenerateSingle ? 'sm:right-10' : 'sm:right-0',
            onClick: () => copyText(uuidJson(currentUuid ? [currentUuid] : []), 'single-json'),
        },
        ...(canGenerateSingle
            ? [
                  {
                      key: 'generate',
                      title: 'Generate UUID',
                      ariaLabel: 'Generate UUID',
                      icon: Icons.Refresh,
                      positionClass: 'sm:right-0',
                      onClick: () => generate(1),
                  },
              ]
            : []),
    ];

    const bulkActions = [
        {
            key: 'generate',
            icon: Icons.Refresh,
            label: 'Generate',
            onClick: () => generate(quantity),
        },
        {
            key: 'copy-all',
            icon: Icons.Copy,
            label: copiedTarget === 'bulk' ? 'Copied' : 'Copy All',
            onClick: () => copyText(bulkOutput, 'bulk'),
        },
        {
            key: 'copy-json',
            icon: Icons.CopyJson,
            label: copiedTarget === 'bulk-json' ? 'Copied' : 'Copy JSON',
            onClick: () => copyText(uuidJson(uuids), 'bulk-json'),
        },
        {
            key: 'download',
            icon: Icons.Download,
            label: 'Download',
            onClick: downloadUuids,
        },
    ];

    return (
        <>
            <Description title={details.title} subtitle={details.subtitle} />

            <div className='mb-6 grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5'>
                {UUID_ROUTE_VERSIONS.map((item) => (
                    <Button key={item} asChild variant={item === version ? 'secondary' : 'outline'} size='sm'>
                        <Link to={item === 'nil' ? '/uuid-generator/nil' : `/uuid-generator/v${item}`}>
                            {item === 'nil' ? 'Nil' : `v${item}`}
                        </Link>
                    </Button>
                ))}
            </div>

            <div className='flex flex-col'>
                <div className='uuidField relative w-full pb-4 sm:pb-8'>
                    <div
                        className={
                            canGenerateSingle ? 'grid gap-2 sm:block' : 'grid gap-2 sm:grid-cols-[1fr_auto_auto]'
                        }>
                        <Input
                            id='uuid-output'
                            value={currentUuid}
                            readOnly
                            variant='underline'
                            spellCheck={false}
                            className='min-w-0 w-full px-2 text-center font-mono text-xs sm:px-32 sm:text-sm md:text-base'
                        />
                        <div
                            className={
                                canGenerateSingle
                                    ? 'grid grid-cols-3 gap-2 sm:block'
                                    : 'grid grid-cols-2 gap-2 sm:contents'
                            }>
                            {singleUuidActions.map((action) => {
                                const Icon = action.icon;

                                return (
                                    <button
                                        key={action.key}
                                        type='button'
                                        className={cn(
                                            'flex h-10 w-full items-center justify-center rounded-lg bg-secondary text-secondary-content hover:bg-secondary/80 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-neutral-950 sm:absolute sm:top-1 sm:h-8 sm:w-8 sm:rounded-full sm:bg-transparent sm:text-base-content sm:hover:bg-accent/10',
                                            action.positionClass
                                        )}
                                        title={action.title}
                                        aria-label={action.ariaLabel}
                                        onClick={action.onClick}>
                                        <Icon className='h-5 w-5' />
                                    </button>
                                );
                            })}
                        </div>
                    </div>
                </div>

                {supportsBulk ? (
                    <section className='flex flex-col gap-3 rounded-lg border border-base-content/20 bg-base-200 p-4'>
                        <div className='grid gap-4 md:grid-cols-[160px_1fr] md:items-end'>
                            <label htmlFor='uuid-quantity' className='flex flex-col gap-2 text-sm text-base-content/70'>
                                Bulk quantity
                                <Input
                                    id='uuid-quantity'
                                    type='number'
                                    min={1}
                                    max={500}
                                    value={quantity}
                                    variant='outline'
                                    className='w-full !text-left'
                                    onChange={(event) => updateQuantity(event.target.value)}
                                />
                            </label>

                            <div className='grid gap-2 sm:grid-cols-4'>
                                {bulkActions.map((action) => {
                                    const Icon = action.icon;

                                    return (
                                        <Button
                                            key={action.key}
                                            type='button'
                                            variant='secondary'
                                            onClick={action.onClick}>
                                            <Icon className='mr-2 h-4 w-4' />
                                            {action.label}
                                        </Button>
                                    );
                                })}
                            </div>
                        </div>

                        <TextArea
                            id='uuid-bulk-output'
                            value={bulkOutput}
                            readOnly
                            spellCheck={false}
                            variant='outline'
                            className='min-h-[220px] w-full resize-y font-mono text-xs sm:text-sm'
                        />

                        {error ? <p className='text-sm text-error'>{error}</p> : null}
                    </section>
                ) : error ? (
                    <p className='text-sm text-error'>{error}</p>
                ) : null}
            </div>
        </>
    );
}
