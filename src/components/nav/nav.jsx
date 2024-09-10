import "./nav.scss";
import "./nav-mobile.scss";
import { Link } from "react-router-dom";
import { BurgerMenu, Settings } from "../icons";
import emote56 from "../../assets/images/emote56.webp";
import { useState } from "react";

export function Nav() {
    const [currentPage, setCurrentPage] = useState(window.location.pathname);
    const [isMobile] = useState(window.innerWidth < 767);

    function handleClick(e) {
        setCurrentPage(e.target.attributes[1].nodeValue);
        isMobile ? mobileMenu() : null;
    }

    function handleButtonCLick(e) {
        if (e.code === "Enter") return handleClick(e);
    }

    function mobileMenu() {
        const nav = document.getElementById("myLinks");
        if (nav.style.display === "block") {
            nav.style.display = "none";
        } else {
            nav.style.display = "block";
        }
    }

    return (
        <header>
            <nav>
                <Link id="homeLink" to="/" tabIndex="-1">
                    <img src={emote56} height={47} width={47} alt="Lesi" />
                </Link>

                <ul id="myLinks" className="myLinks">
                    <li>
                        <Link
                            to=""
                            className={`${"/" === currentPage || "0" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                            tabIndex={0}
                        >
                            Home
                        </Link>
                    </li>
                    <li>
                        <Link
                            to="/aspect-ratio-calculator"
                            className={`${"/aspect-ratio-calculator" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            Aspect Ratio Calculator
                        </Link>
                    </li>
                    <li>
                        <Link
                            to="/pdf-to-png"
                            className={`${"/pdf-to-png" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            PDF Converter
                        </Link>
                    </li>
                    <li>
                        <Link
                            to="/ico-converter"
                            className={`${"/ico-converter" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            Icon Converter
                        </Link>
                    </li>
                    <li>
                        <Link
                            to="/password-generator"
                            className={`${"/password-generator" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            Password Generator
                        </Link>
                    </li>
                    {/* <li>
                        <Link to="/svg-converter"
                            className={`${"/svg-converter" === currentPage ? " active" : ""}`}>
                            SVG Converter
                        </Link>
                    </li> */}
                    <li>
                        <Link
                            to="/weight-converter"
                            className={`${"/weight-converter" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            Weight Converter
                        </Link>
                    </li>
                    {/* <li>
                        <Link to="/video-editor"
                            className={`${"/video-editor" === currentPage ? " active" : ""}`}>
                            Video Editor
                        </Link>
                    </li> */}
                    <li>
                        <Link
                            to="/video-to-mp3"
                            className={`${"/video-to-mp3" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            Video To MP3
                        </Link>
                    </li>
                    <li>
                        <Link
                            to="/video-cropper"
                            className={`${"/video-cropper" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            Video Cropper
                        </Link>
                    </li>
                    <li>
                        <Link
                            to="/minifier"
                            className={`${"/minifier" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            Minifier
                        </Link>
                    </li>
                    {/* <li>
                        <a href="https://www.etsy.com/uk/shop/ShopLesii" target="blank" rel="noreferrer">
                            Store
                        </a>
                    </li> */}
                    <li>
                        <Link
                            to="/settings"
                            className={`${"/settings" === currentPage || "0 0 512 512" === currentPage ? " active" : ""}`}
                            onClick={(e) => handleClick(e)}
                            onKeyDown={(e) => handleButtonCLick(e)}
                        >
                            {!isMobile ? <Settings /> : "Settings"}
                        </Link>
                    </li>
                </ul>

                {isMobile && (
                    <button className="burger-menu" onClick={mobileMenu}>
                        <BurgerMenu />
                    </button>
                )}
            </nav>
        </header>
    );
}
