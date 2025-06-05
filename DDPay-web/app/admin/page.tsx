"use client";

import { AuthGuard } from "@/components/auth/AuthGuard";

const AdminPage = () => {
  return (
    <AuthGuard redirectTo="/admin/dashboard">
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-center">
          <p className="text-lg mb-4">正在加载管理控制台...</p>
        </div>
      </div>
    </AuthGuard>
  );
};

export default AdminPage;
