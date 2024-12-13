'use client';

import { useState, useEffect, useRef, Fragment } from 'react';
import { checkFileTypeValidity, uploadBoxDropOverOrEnter, removeDropZone, uploadBoxOnDrop } from './dropbox-helpers';

export default function DropBox({
    fileType,
    illustration,
    mp3url,
    callback,
}: {
    fileType: string;
    illustration: React.ReactNode;
    mp3url: string | null;
    callback: () => void;
}) {
    const [isBrowser, setIsBrowser] = useState(false);
    const [message, setMessage] = useState('');
    const fileInputRef = useRef<HTMLInputElement>(null);
    const labelMessageRef = useRef<HTMLLabelElement>(null);
    const hiddenDropAreaRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setIsBrowser(true);
        setMessage(
            window.innerWidth < 767
                ? 'Click here to browse\n for a document to upload'
                : 'Drag and drop your file here\n\nOr click here to browse\nyour PC for documents to upload'
        );
    }, []);

    if (!isBrowser) return null;

    function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
        console.log('file change', e.target.files);
        const file = e.target.files?.[0];
        if (file) {
            const { valid, message } = checkFileTypeValidity(fileType, file.type);
            if (valid) {
                console.log('File type is valid');
            } else {
                setMessage(message);
                callback();
            }
        }
    }

    return (
        <>
            <div className="flex flex-col items-center justify-center relative bg-base-100 rounded-lg w-fit mx-auto">
                <div id="dropContainer" className="flex flex-row items-center h-full">
                    <div
                        className={`flex w-[70%] h-full relative bg-base-100 rounded-lg  ${mp3url ? 'p-bottom-[128px]' : ''}`}
                    >
                        <div className="m-5 w-full h-fullrounded-[10px] bg-[var(--main-background-colour)] p-[60px_80px] w-full cursor-pointer rounded-lg">
                            {illustration}
                        </div>
                    </div>
                    <label
                        htmlFor="fileInput"
                        id="uploadText"
                        className="text-pretty break-words h-full cursor-pointer flex flex-col justify-center cursor-pointer pr-1 w-[27%] font-size-1-25 text-left"
                        ref={labelMessageRef}
                    >
                        {message.split('\n').map((line, index) => (
                            <Fragment key={index}>
                                {line}
                                {index < message.split('\n').length - 1 && <br />}
                            </Fragment>
                        ))}
                    </label>
                    <input
                        type="file"
                        className="hidden"
                        ref={fileInputRef}
                        accept={fileType}
                        onChange={(e) => handleFileChange(e)}
                    />
                </div>
            </div>
            <div
                className="absolute inline-block left-0 top-0 w-full h-full rounded-lg border-dashed border-4 opacity-0 z-0 focus-visible:z-40"
                tabIndex={0}
                role="button"
                ref={hiddenDropAreaRef}
                onDragOver={(e) => uploadBoxDropOverOrEnter(e, hiddenDropAreaRef)}
                onDragEnter={(e) => uploadBoxDropOverOrEnter(e, hiddenDropAreaRef)}
                onDragLeave={(e) => removeDropZone(e, hiddenDropAreaRef)}
                onDrop={(e) => uploadBoxOnDrop(e, hiddenDropAreaRef, fileInputRef)}
            ></div>
        </>
    );
}
