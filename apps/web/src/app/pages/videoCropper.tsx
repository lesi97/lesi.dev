import { Description } from '@/components/ui';
import { useFfmpeg } from '@/context/FfmpegContext';
import { Dropbox } from '@/components/layout';
import { illustrations } from '@/components/ui';
import { useRef, useEffect } from 'react';
import { useVideoCropper } from '@/hooks';

export function VideoCropper() {
    const ffmpeg = useFfmpeg();
    const progressBarRef = useRef<SVGPathElement | null>(null);
    const progressBarTotalRef = useRef<SVGPathElement | null>(null);
    const videoCropper = useVideoCropper(ffmpeg);

    useEffect(() => {
        const totalLength = 350;
        const progressLength = (videoCropper.progress / 100) * totalLength;
        const startX = 128.83144;
        const endX = startX + progressLength;
        const progressPath = `M478.83144,356.80298H128.83144c-3.31,0-6,2.69-6,6s2.69,6,6,6H${endX}c3.31,0,6-2.69,6-6s-2.69-6-6-6Z`;

        if (progressBarRef.current) {
            progressBarRef.current.setAttribute('d', progressPath);
            if (videoCropper.progress === 0) {
                progressBarRef.current.classList.add('hidden');
            } else {
                progressBarRef.current.classList.remove('hidden');
            }
        }

        if (progressBarTotalRef.current && videoCropper.progress !== 0) {
            progressBarTotalRef.current.classList.remove('hidden');
        }
    }, [videoCropper.progress]);

    return (
        <>
            <Description
                title='Video Cropper'
                subtitle={
                    <>
                        Drag and drop a video file to crop it to a 9:16 aspect ratio
                        <br />
                        Perfect for YouTube shorts, TikToks or Instagram reels
                    </>
                }
            />
            {ffmpeg.ready && (
                <Dropbox
                    fileType='video/mp4'
                    illustration={
                        <illustrations.Influencer progressBarTotal={progressBarTotalRef} progressBar={progressBarRef} />
                    }
                    callback={(file) => {
                        videoCropper.setVideo(file);
                    }}
                    willLoad={true}
                    loading={videoCropper.loading}
                    progress={videoCropper.progress}
                />
            )}
            {videoCropper.video && (
                <video
                    className='hidden'
                    ref={videoCropper.videoRef}
                    muted
                    src={videoCropper.originalVideoBlobUrl}
                    onLoadedMetadata={() => console.log('Metadata loaded')}
                    onError={(e) => console.error('Error loading video:', e)}
                />
            )}
        </>
    );
}
