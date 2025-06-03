"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { AuthContainer } from "./AuthContainer";

import { useAuthStore } from "@/store/auth";

interface AuthGuardProps {
  children: React.ReactNode;
  redirectTo?: string;
}

export const AuthGuard = ({
  children,
  redirectTo = "/admin/dashboard",
}: AuthGuardProps) => {
  const [isClient, setIsClient] = useState(false);
  const router = useRouter();
  const { isAuthenticated, isTokenExpired } = useAuthStore();

  // 确保在客户端渲染
  useEffect(() => {
    setIsClient(true);
  }, []);

  // 如果不在客户端，显示加载状态
  if (!isClient) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        加载中...
      </div>
    );
  }

  // 如果未登录或令牌已过期，显示登录表单
  if (!isAuthenticated || isTokenExpired()) {
    return <AuthContainer initialMode="login" redirectTo={redirectTo} />;
  }

  // 已登录且令牌有效，显示子组件
  return <>{children}</>;
};
