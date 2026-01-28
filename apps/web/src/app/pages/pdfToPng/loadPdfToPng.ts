export function loadPdfToPng() {
    return import('./pdfToPng').then((module) => {
        return { default: module.PdfToPng };
    });
}
