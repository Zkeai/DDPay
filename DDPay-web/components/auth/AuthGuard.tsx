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
  const [isChecking, setIsChecking] = useState(true);
  const router = useRouter();
  const { isAuthenticated, isTokenExpired } = useAuthStore();

  // 确保在客户端渲染并检查认证状态
  useEffect(() => {
    setIsClient(true);

    // 在状态检查完成后设置isChecking为false
    const checkAuth = () => {
      if (isAuthenticated && !isTokenExpired() && redirectTo) {
        // 仅在显式要求时进行重定向
        const currentPath = window.location.pathname;

        if (currentPath === "/admin") {
          router.push(redirectTo);
        }
      }
      setIsChecking(false);
    };

    // 延迟检查以确保状态已更新
    const timer = setTimeout(checkAuth, 100);

    return () => clearTimeout(timer);
  }, [isAuthenticated, isTokenExpired, redirectTo, router]);

  // 如果不在客户端或正在检查状态，显示加载状态
  if (!isClient || isChecking) {
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
