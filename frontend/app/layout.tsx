import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Shadow IAM - 身份认证系统",
  description: "基于 Next.js 的现代化身份认证系统",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}

