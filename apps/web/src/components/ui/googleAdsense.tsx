export function GoogleAdsense() {
    const id = process.env.GOOGLE_ADSENSE_ID;
    return (
        <script
            async
            src={`https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-${id}`}
            crossOrigin='anonymous'
        />
    );
}
