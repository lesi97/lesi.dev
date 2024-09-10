import "./footer.scss";
import "./footer-mobile.scss";
import {
    YouTube,
    TikTok,
    Discord,
    XCorp,
    Instagram,
    //Threads, 
    Twitch,
    Email,
    BuyMeACoffee
} from "../icons";

export function Footer() {
    return (
        <footer>
            <ul>
                <li><a href="https://twitch.tv/c_lesi" target="_blank" rel="noreferrer"><Twitch /></a></li>
                <li><a href="https://www.youtube.com/@C_Lesi" target="_blank" rel="noreferrer"><YouTube /></a></li>
                <li><a href="https://www.tiktok.com/@c_lesi" target="_blank" rel="noreferrer"><TikTok /></a></li>
                <li><a href="https://discord.gg/RUDhkXT" target="_blank" rel="noreferrer"><Discord /></a></li>
                <li><a href="https://x.com/Chris_Lesi" target="_blank" rel="noreferrer"><XCorp /></a></li>
                <li><a href="https://www.instagram.com/christian_lesi" target="_blank" rel="noreferrer"><Instagram /></a></li>
                {/* <li><a href="https://www.threads.net/christian_lesi" target="_blank" rel="noreferrer"><Threads /></a></li> */}
                <li><a href="https://www.buymeacoffee.com/lesi" target="_blank" rel="noreferrer"><BuyMeACoffee /></a></li>
                <li><a href="mailto:chris.lesi001@gmail.com" target="_blank" rel="noreferrer"><Email /></a></li>
            </ul>
        </footer>
    );
}
