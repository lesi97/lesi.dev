import "./app.scss";
import { createRoot } from "react-dom/client";
import { useState, useEffect, useRef, StrictMode } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Nav } from "./components/nav/nav";
import { Footer } from "./components/footer/footer";
import ltoe from "./assets/images/ltoe.png";
import snowflakeImg1 from "./assets/images/snowflake-1.webp";
import snowflakeImg2 from "./assets/images/snowflake-2.webp";
import snowflakeImg3 from "./assets/images/snowflake-3.webp";
import Snowfall from 'react-snowfall'
import {
    Home,
    Error404,
    AspectRatio,
    VideoToMp3,
    CropToShort,
    Settings,
    VideoEditor,
    PdfToPng,
    PasswordGenerator,
    WeightConverter,
    IconConverter,
    Minifier,
    Privacy,
} from "./routes";
import { Alien } from "./components/icons";

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
    const [isErrorPage, setIsErrorPage] = useState(false);
    const [isChristmas,] = useState(true);
    const [isSnowing, setIsSnowing] = useState();
    const memeRef = useRef(null);

    useEffect(() => {
        const nightModePreference = localStorage.getItem("nightmode");
        if (nightModePreference) {
            setIsNightMode(nightModePreference === "true");
        } else {
            const mediaQuery = window.matchMedia(
                "(prefers-color-scheme: dark)"
            );
            localStorage.setItem("nightmode", mediaQuery.matches);
            setIsNightMode(mediaQuery.matches);
        }
    }, []);

    useEffect(() => {
        if (!isChristmas) return;
        const snowPreference = localStorage.getItem("snow");
        if (snowPreference) {
            setIsSnowing(snowPreference === "true");
        } else {
            localStorage.setItem("snow", 'true');
            setIsSnowing(true);
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
    }

    function toggleSnow() {
        const newMode = !isSnowing;
        setIsSnowing(newMode);
        localStorage.setItem("snow", String(newMode));
    }

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
            memeElement.style.transition = "all 0.2s";
            memeElement.style.opacity = "1";
            memeElement.style.width = "800px";
            memeElement.style.height = "800px";
            setTimeout(() => {
                memeElement.style.transition = "all 1s";
                memeElement.style.width = "650px";
                memeElement.style.height = "650px";
            }, 200);
        };

        const handleMouseLeave = () => {
            memeElement.style.transition = "all 0.5s";
            memeElement.style.width = "1px";
            memeElement.style.height = "1px";
        };

        memeElement.addEventListener("mouseenter", handleMouseEnter);
        memeElement.addEventListener("mouseleave", handleMouseLeave);

        return () => {
            memeElement.removeEventListener("mouseenter", handleMouseEnter);
            memeElement.removeEventListener("mouseleave", handleMouseLeave);
        };
    }, []);

    const snowflake1 = document.createElement('img')
    snowflake1.src = snowflakeImg1;
    const snowflake2 = document.createElement('img')
    snowflake2.src = snowflakeImg2;
    const snowflake3 = document.createElement('img')
    snowflake3.src = snowflakeImg3;
    const images = [snowflake1, snowflake2, snowflake3]

    return (
        <>
            <BrowserRouter>
                <QueryClientProvider client={queryClient}>
                    <Nav
                        toggleNightMode={toggleNightMode}
                        nigthModeState={isNightMode}
                        christmas={isChristmas}
                    />
                    {(isChristmas && isSnowing) ? (
                        <Snowfall
                            color={isNightMode ? "#dee4fd" : "#B1F4E7"}
                            snowflakeCount={isNightMode ? 150 : 500}
                            images={images}
                            radius={[0.25, 15]}
                            style={{
                                top: "47px",
                                height: "calc(100vh - 47px)",
                                position: "fixed",
                                pointerEvents: "none",
                            }} />
                    ) : null}
                    <Routes>
                        <Route
                            path=""
                            element={
                                <Home setError={(e) => setIsErrorPage(e)} />
                            }
                        />
                        <Route
                            path="*"
                            element={
                                <Error404 setError={(e) => setIsErrorPage(e)} />
                            }
                        />
                        <Route
                            path="/aspect-ratio-calculator"
                            element={
                                <AspectRatio
                                    setError={(e) => setIsErrorPage(e)}
                                />
                            }
                        />
                        <Route
                            path="/pdf-to-png"
                            element={
                                <PdfToPng setError={(e) => setIsErrorPage(e)} />
                            }
                        />
                        <Route
                            path="/ico-converter"
                            element={
                                <IconConverter
                                    setError={(e) => setIsErrorPage(e)}
                                />
                            }
                        />
                        <Route
                            path="/password-generator"
                            element={
                                <PasswordGenerator
                                    setError={(e) => setIsErrorPage(e)}
                                />
                            }
                        />

                        {/* <Route path="/svg-converter" element={null} setError={(e) => setIsErrorPage(e)} /> */}

                        <Route
                            path="/weight-converter"
                            element={
                                <WeightConverter
                                    setError={(e) => setIsErrorPage(e)}
                                />
                            }
                        />
                        <Route
                            path="/video-editor"
                            element={
                                <VideoEditor
                                    setError={(e) => setIsErrorPage(e)}
                                />
                            }
                        />
                        <Route
                            path="/video-to-mp3"
                            element={
                                <VideoToMp3
                                    setError={(e) => setIsErrorPage(e)}
                                />
                            }
                        />
                        <Route
                            path="/video-cropper"
                            element={
                                <CropToShort
                                    setError={(e) => setIsErrorPage(e)}
                                />
                            }
                        />
                        <Route
                            path="/minifier"
                            element={
                                <Minifier setError={(e) => setIsErrorPage(e)} />
                            }
                        />
                        <Route
                            path="/privacy"
                            element={
                                <Privacy setError={(e) => setIsErrorPage(e)} />
                            }
                        />
                        <Route
                            path="/settings"
                            element={
                                <Settings
                                    toggleNightMode={toggleNightMode}
                                    nigthModeState={isNightMode}
                                    setError={(e) => setIsErrorPage(e)}
                                    christmas={isChristmas}
                                    toggleSnow={toggleSnow}
                                    snowState={isSnowing}
                                />
                            }
                        />
                    </Routes>
                    {!isErrorPage && <Footer />}
                </QueryClientProvider>
            </BrowserRouter>
            <a
                target="_blank"
                rel="noreferrer"
                href="https://twitch.tv/liight2k"
                tabIndex="-1"
            >
                {/*eslint-disable-next-line jsx-a11y/img-redundant-alt*/}
                <img
                    ref={memeRef}
                    className="forTheMeme"
                    src={ltoe}
                    alt="The sacred image is missing, plz forgiv"
                />
            </a>
        </>
    );
};

const root = createRoot(document.getElementById("root"));
root.render(
    <StrictMode>
        <App />
    </StrictMode>
);
