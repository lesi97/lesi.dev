import "./dropbox.scss";
import "./dropbox-mobile.scss";
import { IllustrationMusic, IllustrationInfluencer, IllustrationPersonalFile, IllustrationImagePost } from "../illustrations";
import { Waveform } from "../../routes/video-editor/video-mp3/waveform";
import { useState, useRef } from "react";

export const DropBox = (props) => {
    const [isPlaying, setIsPlaying] = useState(false);

    const isMobile = window.innerWidth < 767;

    const waveformRef = useRef(null);

    function uploadBoxDropOverOrEnter(evt) {
        evt.preventDefault();
        evt.stopPropagation();
        document.getElementById("hiddenDropArea").classList.add("dropAreaFileHover");
    }

    function uploadBoxOnDrop(evt) {
        evt.preventDefault();
        evt.stopPropagation();
        const fileInput = document.getElementById("fileInput");
        fileInput.files = evt?.dataTransfer?.files;
        removeDropZone();
        if (props.type === "videoToMp3" || props.type === "videoCrop") return props.loadVideo();

        if (typeof props?.fn === "function") {
            props.fn();
        }
    }

    // function uploadBoxOnFileInputChange(evt) {
    //     const fileInput = document.getElementById("fileInput");
    //     let file = fileInput?.files?.[0];
    //     if (!file)
    //         file = !isMobile ? "Drag and drop your file here<br /><br /> Or click here to browse<br /> your PC for documents to upload" :
    //             "Click here to browse<br /> for a document to upload";
    //     const label = fileInput.previousElementSibling;
    //     label.innerHTML = file;
    //     if (props.type === "video") return props.loadVideo();

    //     if (typeof props?.fn === "function") {
    //         props.fn();
    //     }
    // }

    function clickLabel() {
        document.getElementById("fileInput").click();
    }

    function removeDropZone() {
        document
            .getElementById("hiddenDropArea")
            .classList.remove("dropAreaFileHover");
    }

    const handlePlayPause = () => {
        if (!props?.mp3url) return;
        waveformRef.current.playPause();
        setIsPlaying(!isPlaying);
    };

    const handleMp3ChangeFile = (e) => {
        if (props.type === "videoToMp3" && props?.mp3url) {
            return;
        }
        return e.preventDefault();
    }

    return (
        <>
            <div
                className="dropContainer"
                onDragOver={uploadBoxDropOverOrEnter}
                onDragEnter={uploadBoxDropOverOrEnter}
                onDrop={uploadBoxOnDrop}
                onClick={!props?.mp3url ? clickLabel : null}
                onKeyDown={!props?.mp3url ? clickLabel : null}
            >
                <div id="dropContainer" className="dropZone">
                    <div className={`svgContainer ${props?.mp3url ? "additionalPadding" : ""}`} >
                        {props.type === "videoToMp3" && <IllustrationMusic playPause={handlePlayPause} isPlaying={isPlaying} progressBar={props.progressBar} />}
                        {props.type === "videoCrop" && <IllustrationInfluencer progressBar={props.progressBar} progressBarTotal={props.progressBarTotal} />}
                        {props.type === "pdfToPng" && <IllustrationPersonalFile />}
                        {props.type === "imgToIco" && <IllustrationImagePost />}
                        {props?.mp3url && <Waveform url={props?.mp3url} ref={waveformRef} setPlaying={(e) => setIsPlaying(e)} onClick={(e) => { console.log(e) }} />}
                    </div>

                    { /* eslint-disable-next-line jsx-a11y/no-noninteractive-element-to-interactive-role, jsx-a11y/no-noninteractive-element-interactions */}
                    <label htmlFor="fileInput" id="uploadText" onClick={(e) => handleMp3ChangeFile(e)} onKeyDown={(e) => handleMp3ChangeFile(e)}>
                        {!isMobile ? (
                            <>
                                Drag and drop your file here
                                <br />
                                <br /> Or click here to browse
                                <br /> your PC for documents to upload
                            </>
                        ) : (
                            <>
                                Click here to browse<br /> for a document to upload
                            </>
                        )}
                    </label>
                    <input
                        type="file"
                        className="fileInput"
                        id="fileInput"
                        onChange={(e) => uploadBoxOnDrop(e)}
                        accept={
                            (props.type === "videoCrop" || props.type === "videoToMp3") ? "video/*" :
                                props.type === "pdfToPng" ? "application/pdf" :
                                    (props.type === "imageUploader" || props.type === "imgToIco") ? "image/*" :
                                        undefined}
                    />
                </div>
            </div >
            <div
                className="hiddenDropArea"
                id="hiddenDropArea"
                tabIndex="0"
                role="button"
                onKeyDown={(e) => { if (e.code === "Enter") return clickLabel() }}
                onDragOver={uploadBoxDropOverOrEnter}
                onDragEnter={uploadBoxDropOverOrEnter}
                onDragLeave={removeDropZone}
                onDrop={uploadBoxOnDrop}
            ></div>
        </>
    );
}