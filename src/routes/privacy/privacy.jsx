import "./privacy.scss"
import { useEffect } from "react";

export const Privacy = ({ setError }) => {

    useEffect(() => {
        document.title = `Lesi | Privacy`;
        setError(false);
    }, []);

    return (
        <main>
            <div className="privacy">
                <div>
                    <h1>Privacy Policy for lesi.dev</h1>
                    <p>
                        This Privacy Policy document contains types of information that is collected and recorded by lesi.dev and how it is used.
                    </p>

                    <h2>1. General Information</h2>
                    <p>
                        We take your privacy seriously and are committed to complying with the General Data Protection Regulation (GDPR) of the European Union. This Privacy Policy explains how we collect, use, and protect your personal data when you visit our website.
                    </p>

                    <h2>2. Data Collection</h2>
                    <p>
                        We collect personal data from you in the following ways:
                        <ul>
                            <li><strong>Cookies:</strong> Our website uses cookies to improve your browsing experience. Cookies are small text files stored on your device that help us recognize you and remember your preferences.</li>
                        </ul>
                    </p>

                    <h2>3. Use of Data</h2>
                    <p>
                        The data we collect may be used for the following purposes:
                        <ul>
                            <li>To provide and maintain our website</li>
                            <li>To improve user experience</li>
                            <li>To comply with legal obligations</li>
                        </ul>
                    </p>

                    <h2>4. Google AdSense</h2>
                    <p>
                        We use Google AdSense to display advertisements on our website. Google AdSense uses cookies to serve ads based on your prior visits to lesi.dev or other websites. Google's use of advertising cookies enables it and its partners to serve ads to you based on your visit to lesi.dev and/or other sites on the Internet.
                    </p>
                    <p>
                        <strong>Cookies and Data Use:</strong> Google AdSense may collect personal data such as IP addresses and may use cookies and similar technologies to provide you with relevant advertisements. For more information on how Google uses data when you use our site, please visit <a href="https://policies.google.com/technologies/partner-sites">Google's Privacy &amp; Terms site</a>.
                    </p>
                    <h2>5. Data Security</h2>
                    <p>
                        We implement appropriate technical and organizational measures to ensure a level of security appropriate to the risk of processing your personal data. However, please be aware that no method of transmission over the Internet, or method of electronic storage, is 100% secure.
                    </p>

                    <h2>6. Changes to This Privacy Policy</h2>
                    <p>
                        We may update our Privacy Policy from time to time. We encourage you to review this page periodically for any changes. Any updates will be posted on this page with an updated effective date.
                    </p>
                </div>

            </div>
        </main>
    )
}