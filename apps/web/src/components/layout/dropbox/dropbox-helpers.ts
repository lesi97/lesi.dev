import type { AcceptedFileTypes } from '@/schema';

const knownExtensionsByMimeType: Record<string, string[]> = {
    'application/pdf': ['.pdf'],
    'video/mp4': ['.mp4', '.m4v'],
    'video/quicktime': ['.mov', '.qt'],
};

function checkFileTypeValidity(acceptableFileType: AcceptedFileTypes, file: File) {
    const acceptableFileTypes = normaliseAcceptableFileTypes(acceptableFileType);
    if (acceptableFileTypes.some((fileType) => fileMatchesAcceptableType(fileType, file))) {
        const message = `${highlightText(file.name, 'success')} has now been converted\n\nCheck your downloads folder for the converted file`;
        return { valid: true, message };
    }

    const validFileType = checkAcceptableFileType(acceptableFileType);
    const message = `${highlightText('File type not valid!', 'error')}\n\nPlease use a valid ${highlightText(validFileType, 'error')}</span> file\n\nOr click here to browse\n your PC for a ${highlightText(validFileType, 'error')} file to upload`;
    return { valid: false, message };
}

function normaliseAcceptableFileTypes(acceptableFileType: AcceptedFileTypes) {
    const acceptableFileTypes = Array.isArray(acceptableFileType) ? acceptableFileType : [acceptableFileType];
    return acceptableFileTypes.map((fileType) => fileType.toLowerCase().trim()).filter(Boolean);
}

function fileMatchesAcceptableType(acceptableFileType: string, file: File) {
    const fileType = file.type.toLowerCase();
    const fileName = file.name.toLowerCase();

    if (acceptableFileType.startsWith('.')) {
        return fileName.endsWith(acceptableFileType);
    }

    if (acceptableFileType.endsWith('/*')) {
        const generalType = acceptableFileType.split('/')[0];
        const actualType = fileType.split('/')[0];
        return generalType === actualType;
    }

    if (acceptableFileType === fileType) {
        return true;
    }

    return knownExtensionsByMimeType[acceptableFileType]?.some((extension) => fileName.endsWith(extension)) ?? false;
}

function checkAcceptableFileType(fileType: AcceptedFileTypes) {
    const fileTypes = normaliseAcceptableFileTypes(fileType);
    if (fileTypes.some((type) => type === 'image/*' || type.startsWith('image/'))) {
        return 'image';
    }
    if (fileTypes.some((type) => type === 'video/*' || type.startsWith('video/'))) {
        return 'video';
    }
    if (fileTypes.includes('application/pdf')) {
        return 'PDF';
    }
    if (fileTypes.some((type) => type === 'audio/*' || type.startsWith('audio/'))) {
        return 'audio';
    }

    return fileTypes.find((type) => type.startsWith('.'))?.toUpperCase() ?? '';
}

function uploadBoxDropOverOrEnter(
    e: React.DragEvent<HTMLDivElement> | React.FocusEvent<HTMLDivElement, Element>,
    hiddenDropAreaRef: React.RefObject<HTMLDivElement | null>
) {
    e.preventDefault();
    e.stopPropagation();
    if (hiddenDropAreaRef.current) {
        const dropArea = hiddenDropAreaRef.current;
        dropArea.classList.remove('opacity-0');
        dropArea.classList.remove('z-0');
        dropArea.classList.add('bg-opacity-30');
        dropArea.classList.add('bg-slate-500');
        dropArea.classList.add('opacity-100');
        dropArea.classList.add('z-40');
    }
}

function removeDropZone(e: React.DragEvent<HTMLDivElement>, hiddenDropAreaRef: React.RefObject<HTMLDivElement | null>) {
    e.preventDefault();
    e.stopPropagation();
    if (hiddenDropAreaRef.current) {
        const dropArea = hiddenDropAreaRef.current;
        dropArea.classList.add('opacity-0');
        dropArea.classList.add('z-0');
        dropArea.classList.remove('bg-opacity-30');
        dropArea.classList.remove('bg-slate-500');
        dropArea.classList.remove('opacity-100');
        dropArea.classList.remove('border-red');
        dropArea.classList.remove('z-40');
    }
}

function uploadBoxOnDrop(
    e: React.DragEvent<HTMLDivElement>,
    hiddenDropAreaRef: React.RefObject<HTMLDivElement | null>,
    fileInputRef: React.RefObject<HTMLInputElement | null>
) {
    e.preventDefault();
    e.stopPropagation();
    const fileInput = fileInputRef.current;
    if (fileInput && e.dataTransfer.files.length > 0) {
        const dataTransfer = new DataTransfer();
        dataTransfer.items.add(e.dataTransfer.files[0]);
        fileInput.files = dataTransfer.files;
        fileInput.dispatchEvent(new Event('change', { bubbles: true }));
    }
    removeDropZone(e, hiddenDropAreaRef);
}

export function highlightText(
    message: string,
    highlightType?: 'success' | 'info' | 'warning' | 'error' | 'primary' | 'secondary' | 'accent'
): string {
    return `<span class="text-${highlightType || 'accent'} contents font-bold">${message}</span>`;
}

export { checkFileTypeValidity, uploadBoxDropOverOrEnter, removeDropZone, uploadBoxOnDrop };
