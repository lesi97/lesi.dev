import { Description, illustrations } from '@/components/ui';
import { Dropbox } from '@/components/layout';
import { useImageToIcon } from '@/hooks/useImageToIcon';

export function ImageToIcon() {
    const { convertFile } = useImageToIcon();

    return (
        <div>
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
        </div>
    );
}
