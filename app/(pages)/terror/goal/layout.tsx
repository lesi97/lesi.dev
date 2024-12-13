export const metadata = {
    title: 'Terror',
    description: 'Terror',
};

export default async function RootLayout({ children }: { children: React.ReactNode }) {
    return (
        <div>{children}</div>
        //<html lang="en">
        //    <body className={inter.className}>{children}</body>
        //</html>
    );
}
