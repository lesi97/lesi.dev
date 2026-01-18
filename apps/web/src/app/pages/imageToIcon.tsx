import { Description, illustrations } from '@/components/ui';
import { Dropbox } from '@/components/layout';
import { usePageMeta } from '@/hooks';
import { ConvertApp } from '@/lib/iconConverter';
import { useCallback } from 'react';

export function ImageToIcon() {
    usePageMeta({
        title: 'Icon Converter | Lesi',
        description: 'Convert an image to a .ico file',
    });

    const convertFile = useCallback((file: File) => {
        if (!file) {
            return;
        }
        ConvertApp.convert(file);
    }, []);

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
