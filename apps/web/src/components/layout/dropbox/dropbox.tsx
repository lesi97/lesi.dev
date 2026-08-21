import { useState, useEffect, useRef, Fragment, ReactNode } from 'react';
import { checkFileTypeValidity, uploadBoxDropOverOrEnter, removeDropZone, uploadBoxOnDrop } from './dropbox-helpers';
import parseMessage from './message-parser';
import { type AcceptedFileTypes } from '@/schema';
import { highlightText } from './dropbox-helpers';
import { cn } from '@/utils';
import { Waveform } from './waveform';

export function Dropbox({
    fileType,
    illustration,
    url,
    callback,
    willLoad,
    loading,
    progress,
}: {
    fileType: AcceptedFileTypes;
    illustration: ReactNode;
    url?: string | null;
    callback: (file: File) => void;
    willLoad?: boolean;
    loading?: boolean;
    progress?: number;
}) {
    const [isBrowser, setIsBrowser] = useState(false);
    const [message, setMessage] = useState<ReactNode>('');
    const [fileName, setFileName] = useState<string | undefined>();
    const fileInputRef = useRef<HTMLInputElement>(null);
    const labelMessageRef = useRef<HTMLLabelElement>(null);
    const hiddenDropAreaRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setIsBrowser(true);
        setMessage(
            window.innerWidth < 767
                ? parseMessage(`${highlightText('Click here')} to browse\n for a document to upload`)
                : parseMessage(
                      `Drag and drop your file here\n\nOr ${highlightText('click here')} to browse\nyour PC for documents to upload`
                  )
        );
    }, []);

    useEffect(() => {
        if (loading && progress) {
            const message = `${highlightText(fileName || 'Your File', 'info')} is currently being converted, please wait a moment\n\n${parseFloat(progress.toString()).toFixed(2)}% Complete`;
            setMessage(parseMessage(message));
        }
        if (progress && progress > 99) {
            const message = `${highlightText(fileName || 'Your File', 'success')} has now finished converting\n\nCheck your downloads folder for the converted file`;
            setMessage(parseMessage(message));
        }
    }, [loading, progress]);

    if (!isBrowser) return null;
    const inputAccept = Array.isArray(fileType) ? fileType.join(',') : fileType;

    function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0];
        if (file) {
            const { valid, message } = checkFileTypeValidity(fileType, file);
            if (valid) {
                setFileName(file.name);
                const parsedMessage = parseMessage(message);
                if (!willLoad) {
                    setMessage(parsedMessage);
                }
                callback(file);
            } else {
                const parsedMessage = parseMessage(message);
                setMessage(parsedMessage);
            }
        }
    }

    return (
        <>
            <div className='relative z-10 mx-auto flex h-fit w-full flex-row items-center rounded-lg bg-base-200'>
                <div className='flex h-full w-full flex-col-reverse items-center justify-start pt-4 md:flex-row md:pt-0'>
                    <div className={cn('relative flex h-full w-full rounded-lg bg-inherit md:w-[70%]')}>
                        <div
                            className={cn(
                                'm-5 max-h-full w-full cursor-pointer rounded-lg bg-base-100 p-4 lg:p-[60px_80px]',
                                url ? '!pb-[143px]' : ''
                            )}
                            onClick={() => {
                                if (!url) fileInputRef.current?.click();
                            }}>
                            {illustration}
                        </div>
                        {/* {url && <Waveform url={url} />} */}
                    </div>
                    <label
                        htmlFor='fileInput'
                        id='uploadText'
                        className='font-size-1-25 flex h-full w-full cursor-pointer flex-col items-center justify-center text-pretty break-words pr-1 text-left md:w-[27%]'
                        ref={labelMessageRef}
                        onClick={() => {
                            fileInputRef.current?.click();
                        }}>
                        {message}
                        <input
                            type='file'
                            className='hidden'
                            ref={fileInputRef}
                            accept={inputAccept}
                            onChange={(e) => handleFileChange(e)}
                        />
                    </label>
                </div>
            </div>
            <div
                className='absolute left-0 top-0 z-0 inline-block h-full w-full cursor-default rounded-lg opacity-0 outline-dashed outline-4 outline-accent focus:opacity-50 focus-visible:z-40 mb-20'
                tabIndex={0}
                role='button'
                ref={hiddenDropAreaRef}
                onKeyDown={(e) => {
                    if (e.key === 'Enter') fileInputRef.current?.click();
                }}
                onMouseDown={(e) => e.preventDefault()}
                onDragOver={(e) => uploadBoxDropOverOrEnter(e, hiddenDropAreaRef)}
                onDragEnter={(e) => uploadBoxDropOverOrEnter(e, hiddenDropAreaRef)}
                onDragLeave={(e) => removeDropZone(e, hiddenDropAreaRef)}
                onDrop={(e) => uploadBoxOnDrop(e, hiddenDropAreaRef, fileInputRef)}></div>
        </>
    );
}
