function checkFileTypeValidity(acceptableFileType: string, file: File) {
    const fileType = file.type;
    if (acceptableFileType.endsWith('/*')) {
        const generalType = acceptableFileType.split('/')[0];
        const actualType = fileType.split('/')[0];
        if (generalType === actualType) {
            const message = `${highlightText(file.name, 'success')} has now been converted\n\nCheck your downloads folder for the converted file`;
            return { valid: true, message };
        }
    } else if (acceptableFileType === fileType) {
        const message = `${highlightText(file.name, 'success')} has now been converted\n\nCheck your downloads folder for the converted file`;
        return { valid: true, message };
    }

    const validFileType = checkAcceptableFileType(acceptableFileType);
    const message = `${highlightText('File type not valid!', 'error')}\n\nPlease use a valid ${highlightText(validFileType, 'error')}</span> file\n\nOr click here to browse\n your PC for a ${highlightText(validFileType, 'error')} file to upload`;
    return { valid: false, message };
}

function checkAcceptableFileType(fileType: string) {
    switch (fileType) {
        case 'image/*':
            return 'image';
        case 'video/*':
            return 'video';
        case 'video/mp4':
            return 'video';
        case 'application/pdf':
            return 'PDF';
        case 'audio/*':
            return 'audio';
        default:
            return '';
    }
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
