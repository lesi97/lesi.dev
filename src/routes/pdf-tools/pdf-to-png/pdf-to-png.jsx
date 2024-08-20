import "./pdf-to-png.scss";
import { useEffect, useState } from "react";
import { DropBox } from "../../../components/dropbox/dropbox";
import * as pdfjsLib from 'pdfjs-dist/build/pdf';

pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/3.0.279/pdf.worker.min.js'


export const PdfToPng = ({ setError }) => {
    const [pdf, setPdf] = useState(null);

    setError(false);

    useEffect(() => {
        document.title = `Lesi | PDF To PNG Converter`;
    }, []);

    useEffect(() => {
        pdfToPng(pdf);
    }, [pdf])

    const loadPdf = () => {
        const fileInput = document.getElementById("fileInput");
        const file = fileInput?.files?.[0];
        const label = document.getElementById("uploadText");
        if (!file?.type.startsWith("application/pdf")) {
            label.innerHTML = `
                <span class="wrongFileType">File type not valid!</span>
                <span>
                    Please use a valid <span class="wrongFileType">PDF</span> file<br /><br />
                    Or click here to browse<br /> 
                    your PC for a <span class="wrongFileType">PDF</span> file to upload
                </span>
            `;
            label.classList.add("invalid-animation");
            setTimeout(() => {
                label.classList.remove("invalid-animation");
            }, 450);
            return;
        }
        label.innerHTML = file.name;

        setPdf(file);
    };


    const pdfToPng = async (pdfFile) => {
        if (!pdfFile) return;

        try {
            const pdfData = new Uint8Array(await pdfFile.arrayBuffer());
            const pdf = await pdfjsLib.getDocument(pdfData).promise;
            const numPages = pdf.numPages;

            for (let pageNumber = 1; pageNumber <= numPages; pageNumber++) {
                const page = await pdf.getPage(pageNumber);
                const viewport = page.getViewport({ scale: 5.0 });
                const canvas = document.createElement("canvas");
                canvas.width = viewport.width;
                canvas.height = viewport.height;
                const context = canvas.getContext("2d");
                await page.render({ canvasContext: context, viewport: viewport }).promise;

                const pngBlob = await new Promise((resolve) => canvas.toBlob(resolve, "image/png"));
                const downloadLink = document.createElement("a");
                const objectUrl = URL.createObjectURL(pngBlob);
                downloadLink.href = objectUrl;
                const baseName = pdfFile.name.replace(/\.pdf$/i, "");
                if (numPages > 1) {
                    downloadLink.download = `${baseName}_page${pageNumber}.png`;
                } else {
                    downloadLink.download = `${baseName}.png`;
                }
                downloadLink.click();
                URL.revokeObjectURL(objectUrl);
            }
        } catch (error) {
            console.error("PDF to PNG conversion error:", error);
        }
    };


    return (
        <main>
            <div className="pdfToPng">
                <div className="description">
                    <h1>PDF To PNG Converter</h1>
                    <h2>
                        Drag and drop a PDF file to convert it to a PNG and download it<br />
                        Each page will be an individual PNG file
                    </h2>
                </div>
                <DropBox type="pdfToPng" fn={loadPdf} />
            </div>
        </main>
    );
};
