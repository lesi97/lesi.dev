import { Description, illustrations } from '@/components/ui';
import { Dropbox } from '@/components/layout';
import { useImageToIcon } from '@/hooks/useImageToIcon';
import { usePageMeta } from '@/hooks';

export function ImageToIcon() {
    usePageMeta({
        title: 'Icon Converter | Lesi',
        description: 'Convert an image to a .ico file',
    });
    const { convertFile } = useImageToIcon();

    return (
        <>
            <Description
                title='Icon Converter'
                subtitle={
                    <>
                        Drag and drop an image to convert it to a .ico file
                        <br />
                        &nbsp;
                    </>
                }
            />
            <Dropbox
                fileType='image/*'
                illustration={<illustrations.ImagePost />}
                url={null}
                callback={(file) => convertFile(file)}
            />
        </>
    );
}
