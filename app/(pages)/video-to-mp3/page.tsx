import Description from '@/app/_components/description';

export default function VideoToMp3() {
    return (
        <>
            <Description
                title="Video To Mp3"
                subtitle={
                    <>
                        Drag and drop a video file to convert it to MP3 and download it
                        <br />
                        Press the play button once loaded to play the audio in your browser
                    </>
                }
            />
        </>
    );
}
