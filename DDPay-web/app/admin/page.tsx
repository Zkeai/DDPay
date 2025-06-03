"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { AuthGuard } from "@/components/auth/AuthGuard";
import { useAuthStore } from "@/store/auth";

const AdminPage = () => {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();

  // 登录成功后直接跳转到dashboard页面
  useEffect(() => {
    if (isAuthenticated) {
      router.push("/admin/dashboard");
    }
  }, [isAuthenticated, router]);

  return (
    <AuthGuard redirectTo="/admin/dashboard">
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-center">
          <p className="text-lg mb-4">正在跳转到管理控制台...</p>
        </div>
      </div>
    </AuthGuard>
  );
};

export default AdminPage;
