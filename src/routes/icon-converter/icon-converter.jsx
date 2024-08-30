import "./icon-converter.scss";
import { useEffect, useState } from "react";
import { DropBox } from "../../components/dropbox/dropbox";
import { ConvertApp } from "./iconConverter";

export const IconConverter = ({ setError }) => {
    const [image, setImage] = useState(null);

    useEffect(() => {
        document.title = `Lesi | Icon Converter`;
        setError(false);
    }, []);

    useEffect(() => {
        if (image) {
            ConvertApp.convert(image);
        }
    }, [image]);

    const loadImage = () => {
        const fileInput = document.getElementById("fileInput");
        const file = fileInput?.files?.[0];
        const label = document.getElementById("uploadText");

        if (!file?.type.startsWith("image/")) {
            label.innerHTML = `
                <span class="wrongFileType">File type not valid!</span>
                <span>
                    Please use a valid <span class="wrongFileType">image</span> file<br /><br />
                    Or click here to browse<br /> 
                    your PC for a <span class="wrongFileType">image</span> file to upload
                </span>
            `;
            label.classList.add("invalid-animation");
            setTimeout(() => {
                label.classList.remove("invalid-animation");
            }, 450);
            return;
        }
        label.innerHTML = file.name;

        setImage(file);
    };



    return (
        <main>
            <div className="iconConverter">
                <div className="description">
                    <h1>Icon Converter</h1>
                    <h2>
                        Drag and drop an image to convert it to a .ico file<br />&nbsp;
                    </h2>
                </div>

                <DropBox type="imgToIco" fn={loadImage} />

            </div>
        </main>
    )
}