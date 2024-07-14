import "./nav.scss";
import "./nav-mobile.scss";
import { Link } from "react-router-dom";
import { BurgerMenu, Settings } from "../icons";
import emote56 from "../../assets/images/emote56.webp";
import { useState, useEffect } from "react";


// function mobileMenu(e) {
//     e.preventDefault();
//     const navList = document.getElementById("myLinks");
//     navList.classList.toggle("myLinks");
// }

function mobileMenu() {
    const nav = document.getElementById("myLinks");
    if (nav.style.display === "block") {
        nav.style.display = "none";
    } else {
        nav.style.display = "block";
    }
}

const NavButton = (props) => {

    function handleClick() {
        props.setPage(props.link);
        props.mobile ? mobileMenu() : null;
    }


    if (props.active && props.title !== "Mobile Menu") {
        return (
            <>
                <li key={props.id}
                    className={`${props.mobileHide ? "mobileHide" : ""}${props.desktopHide ? " desktopHide" : ""}`}
                    onClick={handleClick}
                    id={props.title === "Mobile Menu" ? "mobileMenu" : ""}
                    tabIndex="-1">
                    <Link to={props.link}
                        className={`${props.link === props.activePage ? " active" : ""}`}>
                        {props.title !== "Settings" && props.title}
                        {props.title === "Settings" && !props.mobile && (
                            <span className="mobileHide">
                                <Settings />
                            </span>
                        )}
                        {props.title === "Settings" && props.mobile && (
                            <span className="mobileShow">Settings</span>
                        )}
                    </Link>
                </li>

            </>
        );
    }
    return null;
};

export function Nav() {
    const [currentPage, setCurrentPage] = useState(window.location.pathname);
    const [isMobile, setIsMobile] = useState(window.innerWidth < 767);

    function setPage(e) {
        setCurrentPage(e)
    }

    const navbarData = [
        {
            ID: 0,
            Title: "Home",
            Image: emote56,
            Link: "/",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 1,
            Title: "Aspect Ratio Calculator",
            Link: "/aspect-ratio-calculator",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 2,
            Title: "PDF Converter",
            Link: "/pdf-to-png",
            MobileHide: true,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 3,
            Title: "Icon Converter",
            Link: "/ico-converter",
            MobileHide: true,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 4,
            Title: "Password Generator",
            Link: "/password-generator",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 5,
            Title: "SVG Converter",
            Link: "/svg-converter",
            MobileHide: true,
            DesktopHide: false,
            Active: false,
        },
        {
            ID: 6,
            Title: "Weight Converter",
            Link: "/weight-converter",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 7,
            Title: "Video Editor",
            Link: "/video-editor",
            MobileHide: false,
            DesktopHide: false,
            Active: false,
        },
        {
            ID: 8,
            Title: "Video To MP3",
            Link: "/video-to-mp3",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 9,
            Title: "Video Cropper",
            Link: "/video-cropper",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 10,
            Title: "Minifier",
            Link: "/minifier",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 11,
            Title: "Store",
            Link: "https://www.etsy.com/uk/shop/ShopLesii",
            MobileHide: false,
            DesktopHide: false,
            Active: false,
        },
        {
            ID: 12,
            Title: "Mobile Menu",
            Link: "",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
        {
            ID: 13,
            Title: "Settings",
            Link: "/settings",
            MobileHide: false,
            DesktopHide: false,
            Active: true,
        },
    ];

    return (
        <header>
            <nav>
                <Link id="homeLink" to={navbarData[0].Link}>
                    <img
                        src={navbarData[0].Image}
                        height={47}
                        width={47}
                        alt="Lesi"
                    />
                </Link>

                <ul id="myLinks" className="myLinks">
                    {navbarData.map((nav) => {
                        return (
                            <NavButton
                                key={nav.ID}
                                id={nav.ID}
                                title={nav.Title}
                                link={nav.Link}
                                mobileHide={nav.MobileHide}
                                desktopHide={nav.DesktopHide}
                                active={nav.Active}
                                activePage={currentPage}
                                setPage={e => setPage(e)}
                                mobile={isMobile}
                            />
                        );
                    })}
                </ul>

                {isMobile && (
                    <div className="burger-menu" onClick={mobileMenu}><BurgerMenu /></div>
                )}
            </nav>
        </header>
    );
}
