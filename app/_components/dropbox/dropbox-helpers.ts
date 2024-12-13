function checkFileTypeValidity(acceptableFileType: string, fileType: string) {
    if (acceptableFileType.endsWith('/*')) {
        const generalType = acceptableFileType.split('/')[0];
        const actualType = fileType.split('/')[0];
        if (generalType === actualType) {
            const message = '';
            return { valid: true, message };
        }
    } else if (acceptableFileType === fileType) {
        const message = '';
        return { valid: true, message };
    }

    const message = `File type not valid!\n\nPlease use a valid ${checkAcceptableFileType(acceptableFileType)} file\n\nOr click here to browse\n your PC for a ${checkAcceptableFileType(acceptableFileType)} file to upload`;
    return { valid: false, message };
}

function checkAcceptableFileType(fileType: string) {
    switch (fileType) {
        case 'image/*':
            return 'image';
        case 'video/*':
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
    e: React.DragEvent<HTMLDivElement>,
    hiddenDropAreaRef: React.RefObject<HTMLDivElement>
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
        dropArea.classList.add('border-red');
        dropArea.classList.add('z-40');
    }
}

function removeDropZone(e: React.DragEvent<HTMLDivElement>, hiddenDropAreaRef: React.RefObject<HTMLDivElement>) {
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
    hiddenDropAreaRef: React.RefObject<HTMLDivElement>,
    fileInputRef: React.RefObject<HTMLInputElement>
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

export { checkFileTypeValidity, uploadBoxDropOverOrEnter, removeDropZone, uploadBoxOnDrop };
