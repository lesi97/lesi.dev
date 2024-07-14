import "./footer.scss";
import "./footer-mobile.scss";
import { YouTube, TikTok, Discord, XCorp, Instagram, Threads, Twitch, Email, BuyMeACoffee } from "../icons";

export function Footer() {
    const footerData = [
        {
            Title: "Twitch",
            Image: <Twitch />,
            Link: "https://twitch.tv/c_lesi",
            Active: true,
            ID: 0,
        },
        {
            Title: "YouTube",
            Image: <YouTube />,
            Link: "https://www.youtube.com/@C_Lesi",
            Active: true,
            ID: 1,
        },
        {
            Title: "TikTok",
            Image: <TikTok />,
            Link: "https://www.tiktok.com/@c_lesi",
            Active: true,
            ID: 2,
        },
        {
            Title: "Discord",
            Image: <Discord />,
            Link: "https://discord.gg/RUDhkXT",
            Active: true,
            ID: 3,
        },
        {
            Title: "X (Formerly Twitter)",
            Image: <XCorp />,
            Link: "https://x.com/Chris_Lesi",
            Active: true,
            ID: 4,
        },
        {
            Title: "Instagram",
            Image: <Instagram />,
            Link: "https://www.instagram.com/christian_lesi",
            Active: true,
            ID: 5,
        },
        {
            Title: "Threads",
            Image: <Threads />,
            Link: "https://www.threads.net/christian_lesi",
            Active: false,
            ID: 6,
        },
        {
            Title: "Buy Me A Coffee",
            Image: <BuyMeACoffee />,
            Link: "https://www.buymeacoffee.com/lesi",
            Active: true,
            ID: 7,
        },
        {
            Title: "Email",
            Image: <Email />,
            Link: "mailto:chris.lesi001@gmail.com",
            Active: true,
            ID: 8,
        },
    ];

    return (
        <footer>
            <ul>
                {footerData.map((socials) => {
                    if (socials.Active) {
                        return (
                            <li key={socials.ID}>
                                <a href={socials.Link} target="_blank" rel="noreferrer">
                                    {socials.Image}
                                </a>
                            </li>
                        );
                    }
                    return null;
                })}
            </ul>
        </footer>
    );
}
