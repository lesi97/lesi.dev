import { Dropbox, Description } from '@/components/layout';
import { usePageMeta, useVideoToMp3 } from '@/hooks';
import { illustrations } from '@/components/ui';
import { useMusic } from '@/context/MusicContext';

export function VideoToMp3() {
    usePageMeta({
        title: 'Video To MP3 | Lesi',
        description: 'Convert a video to MP3 file format',
    });
    const { isPlaying, setIsPlaying, ffmpeg } = useMusic();
    const videoToMp3 = useVideoToMp3(ffmpeg);

    return (
        <>
            <Description
                title='Video To Mp3'
                subtitle={
                    <>
                        Drag and drop a video file to convert it to MP3 and download it
                        <br />
                        Press the play button once loaded to play the audio in your browser
                    </>
                }
            />
            {ffmpeg.ready && (
                <Dropbox
                    fileType='video/mp4'
                    illustration={
                        <illustrations.Music
                            playPause={() => {
                                if (videoToMp3.mp3) {
                                    setIsPlaying((prev) => !prev);
                                }
                            }}
                            isPlaying={isPlaying}
                            progressBar={videoToMp3.progress}
                        />
                    }
                    url={videoToMp3.mp3}
                    callback={(file) => {
                        videoToMp3.setVideo(file);
                    }}
                    willLoad={true}
                    loading={videoToMp3.loading}
                    progress={videoToMp3.progress}
                />
            )}
            {videoToMp3.video && (
                <video
                    className='hidden'
                    ref={videoToMp3.videoRef}
                    muted
                    src={videoToMp3.originalVideoBlobUrl}
                    onLoadedMetadata={() => console.log('Metadata loaded')}
                    onError={(e) => console.error('Error loading video:', e)}
                />
            )}
        </>
    );
}
