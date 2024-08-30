import "./404.scss";
import "./whitebeard.scss";
import "./404-mobile.scss";
import { Link } from "react-router-dom";
import { useEffect } from "react";

export const Error404 = ({ setError }) => {

    useEffect(() => {
        document.title = `Lesi | 404`;
        setError(true);
    }, []);

    return (
        <Link to="/">
            <div className="error404">
                <div className="background-container">
                    <div className="clouds"></div>
                    {/* <h1>4<span className="hideMobile">&emsp;</span><span className="hideDesktop">o</span>4</h1> */}
                    <h2>
                        The One Piece may be real<br />
                        but this page isn&#39;t
                    </h2>
                </div>
            </div>
        </Link>
    )
}