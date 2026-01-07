import { getDocument, GlobalWorkerOptions } from 'pdfjs-dist';
import workerUrl from 'pdfjs-dist/build/pdf.worker.min.mjs?url';
import { useEffect, useRef, useCallback } from 'react';

export function usePdfToPng() {
    const readyRef = useRef(false);

    useEffect(() => {
        if (readyRef.current) {
            return;
        }

        GlobalWorkerOptions.workerSrc = workerUrl;
        readyRef.current = true;
    }, []);

    const convertPdfToPng = useCallback(async (pdfFile: File) => {
        if (!pdfFile || !readyRef.current) {
            return;
        }

        try {
            const pdfData = new Uint8Array(await pdfFile.arrayBuffer());
            const pdf = await getDocument(pdfData).promise;
            const numPages = pdf.numPages;

            for (let pageNumber = 1; pageNumber <= numPages; pageNumber++) {
                const page = await pdf.getPage(pageNumber);
                const viewport = page.getViewport({ scale: 5.0 });
                const canvas = document.createElement('canvas');
                canvas.width = viewport.width;
                canvas.height = viewport.height;
                const context = canvas.getContext('2d');
                if (!context) {
                    return;
                }
                await page.render({
                    canvasContext: context,
                    viewport: viewport,
                }).promise;

                const pngBlob = await new Promise((resolve) => canvas.toBlob(resolve, 'image/png'));
                const downloadLink = document.createElement('a');
                const objectUrl = URL.createObjectURL(pngBlob as Blob);
                downloadLink.href = objectUrl;
                const baseName = pdfFile.name.replace(/\.pdf$/i, '');
                if (numPages > 1) {
                    downloadLink.download = `${baseName}_page_${pageNumber}.png`;
                } else {
                    downloadLink.download = `${baseName}.png`;
                }
                downloadLink.click();
                URL.revokeObjectURL(objectUrl);
            }
        } catch (error) {
            alert(`PDF to PNG conversion error: ${error}`);
            console.error('PDF to PNG conversion error:', error);
        }
    }, []);

    return { convertPdfToPng };
}
