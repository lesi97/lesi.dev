import Description from '@/app/_components/description';

export default function VideoCropper() {
    return (
        <>
            <Description
                title="Video Cropper"
                subtitle={
                    <>
                        Drag and drop a video file to crop it to a 9:16 aspect ratio
                        <br />
                        Perfect for YouTube shorts, TikToks or Instagram reels
                    </>
                }
            />
        </>
    );
}
