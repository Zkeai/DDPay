"use client";

import { usePathname } from "next/navigation";

import AdminLayout from "@/components/layout/AdminLayout";
import { TitleProvider } from "@/components/TitleContext";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const isLoginPage = pathname === "/admin";

  if (isLoginPage) {
    return <>{children}</>;
  }

  return (
    <TitleProvider>
      <AdminLayout>{children}</AdminLayout>
    </TitleProvider>
  );
}
