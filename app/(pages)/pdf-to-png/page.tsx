import Description from '@/app/_components/description';
import DropBox from '@/app/_components/dropbox';
import IllustrationPersonalFile from '@/app/_assets/_illustrations/personal-file';
export default function PdfToPng() {
    return (
        <>
            <Description
                title="Pdf To Png"
                subtitle={
                    <>
                        Drag and drop a PDF file to convert it to a PNG and download it
                        <br />
                        Each page will be an individual PNG file
                    </>
                }
            />
            <DropBox fileType="application/pdf" illustration={<IllustrationPersonalFile />} mp3url={null} />
        </>
    );
}
