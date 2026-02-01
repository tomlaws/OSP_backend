import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Online Survey Platform",
  description: "Create surveys, collect responses, and get AI-powered insights.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased min-h-screen bg-white text-black selection:bg-[#D80000] selection:text-white">
        {children}
      </body>
    </html>
  );
}
