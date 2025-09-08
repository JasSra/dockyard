export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html>
      <body style={{ fontFamily: 'sans-serif', margin: 20 }}>
        <h1>Dockyard</h1>
        {children}
      </body>
    </html>
  );
}
