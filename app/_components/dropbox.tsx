'use client';

import { useState, useEffect, useRef } from 'react';

//accept={
//    props.type === 'videoCrop' || props.type === 'videoToMp3'
//        ? 'video/*'
//        : props.type === 'pdfToPng'
//          ? 'application/pdf'
//          : props.type === 'imageUploader' || props.type === 'imgToIco'
//            ? 'image/*'
//            : undefined
//}

export default function DropBox({ fileType, illustration, mp3url }) {
    const [isBrowser, setIsBrowser] = useState(false);
    const [isMobile, setIsMobile] = useState(true);

    useEffect(() => {
        setIsBrowser(true);
        setIsMobile(window.innerWidth < 767);
    }, []);

    if (!isBrowser) return null;

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

                    {/* eslint-disable-next-line jsx-a11y/no-noninteractive-element-to-interactive-role, jsx-a11y/no-noninteractive-element-interactions */}
                    <label
                        htmlFor="fileInput"
                        id="uploadText"
                        className="text-pretty break-words h-full cursor-pointer flex flex-col justify-center cursor-pointer pr-1 w-[27%] font-size-1-25 text-left"
                    >
                        {!isMobile ? (
                            <>
                                Drag and drop your file here
                                <br />
                                <br /> Or click here to browse
                                <br /> your PC for documents to upload
                            </>
                        ) : (
                            <>
                                Click here to browse
                                <br /> for a document to upload
                            </>
                        )}
                    </label>
                    <input type="file" className="hidden" id="fileInput" accept={fileType} />
                </div>
            </div>
            <div
                className="absolute inline-block left-0 top-0 w-full h-full rounded border-dashed focus-visible:opacity-100 focus-visible:z-40 focus-visible:border-2 focus-visible:border-red"
                id="hiddenDropArea"
                tabIndex="0"
                role="button"
            ></div>
        </>
    );
}
