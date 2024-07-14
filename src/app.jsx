import "./app.scss";
import { createRoot } from "react-dom/client";
import { useState, useEffect, useRef } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { Nav } from "./components/nav/nav";
import { Footer } from "./components/footer/footer";
import ltoe from "./assets/images/ltoe.png";
import { Home, Error404, AspectRatio, VideoToMp3, CropToShort, Settings, VideoEditor, PdfToPng, PasswordGenerator, WeightConverter, IconConverter, Minifier } from "./routes";

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: Infinity,
            cacheTime: Infinity,
        },
    },
});


const App = () => {
    const [isNightMode, setIsNightMode] = useState();
    const memeRef = useRef(null);


    useEffect(() => {
        const nightModePreference = localStorage.getItem("nightmode");
        if (nightModePreference) {
            setIsNightMode(nightModePreference === "true");
        } else {
            const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
            localStorage.setItem("nightmode", mediaQuery.matches);
            setIsNightMode(mediaQuery.matches);
        }
    }, []);

    useEffect(() => {
        if (isNightMode) {
            document.body.classList.add("night");
        } else {
            document.body.classList.remove("night");
        }
    }, [isNightMode]);

    function toggleNightMode() {
        const newMode = !isNightMode;
        setIsNightMode(newMode);
        localStorage.setItem("nightmode", String(newMode));
    };

    const randomPosition = () => {
        const windowHeight = window.innerHeight - memeRef.current.offsetHeight;
        const windowWidth = window.innerWidth - memeRef.current.offsetWidth;
        const randomTop = Math.floor(Math.random() * windowHeight);
        const randomLeft = Math.floor(Math.random() * windowWidth);

        return { top: randomTop, left: randomLeft };
    };

    useEffect(() => {
        const memeElement = memeRef.current;
        const position = randomPosition();
        memeElement.style.top = `${position.top}px`;
        memeElement.style.left = `${position.left}px`;

        const handleMouseEnter = () => {
            memeElement.style.transition = 'all 0.2s';
            memeElement.style.opacity = '1';
            memeElement.style.width = '800px';
            memeElement.style.height = '800px';
            setTimeout(() => {
                memeElement.style.transition = 'all 1s';
                memeElement.style.width = '650px';
                memeElement.style.height = '650px';
            }, 200);
        };

        const handleMouseLeave = () => {
            memeElement.style.transition = 'all 0.5s';
            memeElement.style.width = '1px';
            memeElement.style.height = '1px';
        };

        memeElement.addEventListener('mouseenter', handleMouseEnter);
        memeElement.addEventListener('mouseleave', handleMouseLeave);

        return () => {
            memeElement.removeEventListener('mouseenter', handleMouseEnter);
            memeElement.removeEventListener('mouseleave', handleMouseLeave);
        };
    }, []);

    return (
        <>
            <BrowserRouter>
                <QueryClientProvider client={queryClient}>
                    <Nav toggleNightMode={toggleNightMode} nigthModeState={isNightMode} />
                    {/* <main> */}
                    <Routes>
                        <Route path="" element={<Home />} />
                        <Route path="*" element={<Error404 />} />
                        <Route path="/aspect-ratio-calculator" element={<AspectRatio />} />
                        <Route path="/pdf-to-png" element={<PdfToPng />} />
                        <Route path="/ico-converter" element={<IconConverter />} />
                        <Route path="/password-generator" element={<PasswordGenerator />} />
                        <Route path="/svg-converter" element={null} />
                        <Route path="/weight-converter" element={<WeightConverter />} />
                        <Route path="/video-editor" element={<VideoEditor />} />
                        <Route path="/video-to-mp3" element={<VideoToMp3 />} />
                        <Route path="/video-cropper" element={<CropToShort />} />
                        <Route path="/minifier" element={<Minifier />} />
                        <Route path="/settings" element={<Settings toggleNightMode={toggleNightMode} nigthModeState={isNightMode} />} />
                    </Routes>
                    {/* </main> */}
                    <Footer />
                </QueryClientProvider>
            </BrowserRouter>
            <a target="_blank" rel="noreferrer" href="https://twitch.tv/liight2k">
                {/*eslint-disable-next-line jsx-a11y/img-redundant-alt*/}
                <img ref={memeRef} className="forTheMeme" src={ltoe} alt="The sacred image is missing, plz forgiv" />
            </a>
        </>
    )

}


const root = createRoot(document.getElementById("root"));
root.render(<App />);
