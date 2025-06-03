"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

import { LoginForm } from "./LoginForm";
import { RegisterForm } from "./RegisterForm";
import { ForgotPasswordForm } from "./ForgotPasswordForm";

interface AuthContainerProps {
  initialMode?: "login" | "register" | "forgot_password";
  onSuccess?: () => void;
  redirectTo?: string;
}

export const AuthContainer = ({
  initialMode = "login",
  onSuccess,
  redirectTo = "/admin/dashboard",
}: AuthContainerProps) => {
  const [mode, setMode] = useState<"login" | "register" | "forgot_password">(
    initialMode
  );
  const [mounting, setMounting] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // 添加渐入动画
    setMounting(false);
  }, []);

  const handleAuthSuccess = () => {
    if (onSuccess) {
      onSuccess();
    } else if (redirectTo) {
      router.push(redirectTo);
    }
  };

  const handleSwitchMode = (
    newMode: "login" | "register" | "forgot_password"
  ) => {
    // 切换表单时添加过渡动画
    setMounting(true);
    setTimeout(() => {
      setMode(newMode);
      setMounting(false);
    }, 300);
  };

  return (
    <div className="flex flex-col justify-center items-center min-h-[100vh] w-full relative overflow-hidden bg-gradient-to-b from-white via-blue-50 to-white dark:from-gray-950 dark:via-gray-900 dark:to-gray-950">
      {/* AuthJS风格的背景 */}
      <div className="absolute inset-0 overflow-hidden z-0">
        {/* 网格背景 */}
        <div
          className="absolute inset-0 opacity-[0.03] dark:opacity-[0.05]"
          style={{
            backgroundImage: `linear-gradient(to right, #ddd 1px, transparent 1px), linear-gradient(to bottom, #ddd 1px, transparent 1px)`,
            backgroundSize: "40px 40px",
          }}
        />

        {/* 渐变圆形 - 左上 */}
        <div className="absolute -top-20 -left-20 w-[30rem] h-[30rem] rounded-full bg-gradient-to-br from-blue-400/20 via-indigo-400/20 to-purple-400/20 dark:from-blue-500/20 dark:via-indigo-500/20 dark:to-purple-500/20 blur-[6rem] opacity-70" />

        {/* 渐变圆形 - 右下 */}
        <div className="absolute -bottom-32 -right-32 w-[40rem] h-[40rem] rounded-full bg-gradient-to-br from-blue-400/20 via-indigo-400/20 to-purple-400/20 dark:from-blue-500/20 dark:via-indigo-500/20 dark:to-purple-500/20 blur-[7rem] opacity-70" />

        {/* 中间装饰 */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[25rem] h-[25rem] rounded-full bg-gradient-to-br from-cyan-300/10 via-blue-300/10 to-indigo-300/10 dark:from-cyan-500/10 dark:via-blue-500/10 dark:to-indigo-500/10 blur-[5rem] opacity-60" />
      </div>

      {/* 品牌标志 */}
      <div className="relative z-10 mb-6 text-center">
        <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-400 text-transparent bg-clip-text">
          DDPay
        </h1>
        <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
          安全支付，便捷生活
        </p>
      </div>

      {/* 表单容器 */}
      <div
        className={`relative z-10 transition-all duration-300 transform ${
          mounting ? "opacity-0 scale-95" : "opacity-100 scale-100"
        } w-full max-w-3xl px-4`}
      >
        {mode === "login" ? (
          <LoginForm
            onSuccess={handleAuthSuccess}
            onRegisterClick={() => handleSwitchMode("register")}
            onForgotPasswordClick={() => handleSwitchMode("forgot_password")}
          />
        ) : mode === "register" ? (
          <RegisterForm
            onSuccess={handleAuthSuccess}
            onLoginClick={() => handleSwitchMode("login")}
          />
        ) : (
          <ForgotPasswordForm
            onSuccess={handleAuthSuccess}
            onBackToLogin={() => handleSwitchMode("login")}
          />
        )}
      </div>

      {/* 页脚 */}
      <div className="relative z-10 mt-6 text-xs text-gray-500 dark:text-gray-400">
        &copy; {new Date().getFullYear()} DDPay. 保留所有权利
      </div>
    </div>
  );
};
