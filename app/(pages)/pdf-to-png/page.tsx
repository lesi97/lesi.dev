'use client';

import Description from '@/app/_components/description';
import DropBox from '@/app/_components/dropbox';
import IllustrationPersonalFile from '@/app/_assets/_illustrations/personal-file';

export default function PdfToPng() {
    function doStuff() {
        console.log('do stuff');
    }

    return (
        <>
            <Description
                title="PDF To PNG"
                subtitle={
                    <>
                        Drag and drop a PDF file to convert it to a PNG and download it
                        <br />
                        Each page will be an individual PNG file
                    </>
                }
            />
            <DropBox
                fileType="application/pdf"
                illustration={<IllustrationPersonalFile />}
                mp3url={null}
                callback={() => doStuff()}
            />
        </>
    );
}
