'use client';

import Description from '@/app/_components/description';
import DropBox from '@/app/_components/dropbox';
import IllustrationImagePost from '@/app/_assets/_illustrations/image-post';

export default function IconConverter() {
    return (
        <>
            <Description
                title="Icon Converter"
                subtitle={
                    <>
                        Drag and drop an image to convert it to a .ico file
                        <br />
                        &nbsp;
                    </>
                }
            />
            <DropBox fileType="image/*" illustration={<IllustrationImagePost />} mp3url={null} />
        </>
    );
}
