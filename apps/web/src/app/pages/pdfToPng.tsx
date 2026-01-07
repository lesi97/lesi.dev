'use client';

import { useCallback } from 'react';
import { Dropbox } from '@/components/layout';
import { Description, illustrations } from '@/components/ui';
import { usePdfToPng } from '@/hooks/usePdfToPng';
import { usePageMeta } from '@/hooks';

export function PdfToPng() {
    usePageMeta({
        title: 'PDF To PNG | Lesi',
        description: 'Convert a PDF to a PNG',
    });
    const { convertPdfToPng } = usePdfToPng();

    const handleFileDrop = useCallback((file: File) => {
        convertPdfToPng(file);
    }, []);

    return (
        <>
            <Description
                title='PDF To PNG'
                subtitle={
                    <>
                        Drag and drop a PDF file to convert it to a PNG and download it
                        <br />
                        Each page will be an individual PNG file
                    </>
                }
            />
            <Dropbox
                fileType='application/pdf'
                illustration={<illustrations.PersonalFile />}
                url={null}
                callback={handleFileDrop}
            />
        </>
    );
}
