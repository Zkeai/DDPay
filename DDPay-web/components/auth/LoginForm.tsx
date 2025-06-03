"use client";

import { useState } from "react";
import {
  ArrowRightIcon,
  LockClosedIcon,
  EnvelopeIcon,
  ShieldCheckIcon,
  CheckCircleIcon,
  ExclamationCircleIcon,
  ArrowPathIcon,
} from "@heroicons/react/24/outline";

import { useAuthStore } from "@/store/auth";
import { login } from "@/lib/afetch";

interface LoginFormProps {
  onSuccess?: () => void;
  onRegisterClick?: () => void;
  onForgotPasswordClick?: () => void;
}

export const LoginForm = ({
  onSuccess,
  onRegisterClick,
  onForgotPasswordClick,
}: LoginFormProps) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const { login: storeLogin } = useAuthStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await login(email, password);

      if (response.code === 200) {
        const { user, access_token, refresh_token, expires_in } = response.data;

        storeLogin(user, access_token, refresh_token, expires_in);

        if (onSuccess) {
          onSuccess();
        }
      } else {
        setError(response.msg || "登录失败");
      }
    } catch (err: any) {
      setError(err.message || "登录失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full max-w-3xl bg-white/70 dark:bg-gray-800/70 rounded-2xl shadow-xl overflow-hidden backdrop-blur-sm border border-gray-200/50 dark:border-gray-700/50">
      <div className="flex flex-col lg:flex-row min-h-[480px]">
        {/* 左侧彩色部分 - 仅在大屏幕及以上显示 */}
        <div className="hidden lg:block lg:w-5/12 relative overflow-hidden">
          <div className="absolute inset-0 bg-gradient-to-br from-blue-600 via-blue-600 to-indigo-700">
            <div
              className="absolute inset-0 opacity-20"
              style={{
                backgroundImage:
                  "url('data:image/svg+xml,%3Csvg width='100' height='100' viewBox='0 0 100 100' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M11 18c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm48 25c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm-43-7c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm63 31c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM34 90c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm56-76c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM12 86c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm28-65c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm23-11c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-6 60c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm29 22c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zM32 63c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm57-13c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-9-21c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM60 91c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM35 41c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM12 60c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2z' fill='%23ffffff' fill-opacity='1' fill-rule='evenodd'/%3E%3C/svg%3E')",
              }}
            />
          </div>
          <div className="relative h-full flex flex-col justify-between p-6 z-10">
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-white/20 rounded-xl backdrop-blur-sm">
                <ShieldCheckIcon className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-white text-lg font-bold">DDPay 登录</h3>
            </div>

            <div className="space-y-4">
              <div>
                <h2 className="text-white text-2xl font-bold leading-tight">
                  欢迎回到您的账户
                </h2>
                <p className="text-blue-100 mt-2 text-sm">
                  登录DDPay，继续您的安全支付体验
                </p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <div className="flex-shrink-0 w-6 h-6 bg-blue-500/30 rounded-full flex items-center justify-center">
                    <ShieldCheckIcon className="w-4 h-4 text-white" />
                  </div>
                  <p className="text-white text-xs">安全加密登录</p>
                </div>
                <div className="flex items-center space-x-2">
                  <div className="flex-shrink-0 w-6 h-6 bg-blue-500/30 rounded-full flex items-center justify-center">
                    <CheckCircleIcon className="w-4 h-4 text-white" />
                  </div>
                  <p className="text-white text-xs">快速支付处理</p>
                </div>
                <div className="flex items-center space-x-2">
                  <div className="flex-shrink-0 w-6 h-6 bg-blue-500/30 rounded-full flex items-center justify-center">
                    <ArrowPathIcon className="w-4 h-4 text-white" />
                  </div>
                  <p className="text-white text-xs">多平台同步支持</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 右侧表单部分 */}
        <div className="w-full lg:w-7/12 flex items-center justify-center p-4 lg:p-6">
          <div className="w-full max-w-md">
            <div className="lg:hidden mb-6 text-center">
              <div className="inline-flex items-center justify-center p-2 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-xl mb-3">
                <ShieldCheckIcon className="w-6 h-6 text-white" />
              </div>
              <h2 className="text-xl font-bold text-gray-800 dark:text-white">
                欢迎回来
              </h2>
              <p className="text-gray-600 dark:text-gray-300 mt-1 text-sm">
                请登录您的账户以继续
              </p>
            </div>

            <div className="hidden lg:block mb-4">
              <h2 className="text-xl font-bold text-gray-800 dark:text-white">
                欢迎回来
              </h2>
              <p className="text-gray-600 dark:text-gray-300 mt-1 text-sm">
                请登录您的账户以继续
              </p>
            </div>

            {error && (
              <div className="mb-4 bg-red-50 dark:bg-red-900/30 p-3 rounded-xl border border-red-100 dark:border-red-800">
                <div className="flex items-start">
                  <div className="flex-shrink-0">
                    <ExclamationCircleIcon className="w-4 h-4 text-red-500" />
                  </div>
                  <div className="ml-2">
                    <p className="text-xs font-medium text-red-800 dark:text-red-300">
                      {error}
                    </p>
                  </div>
                </div>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-3">
              <div>
                <label
                  className="block text-gray-700 dark:text-gray-300 text-xs font-medium mb-1"
                  htmlFor="email"
                >
                  邮箱地址
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <EnvelopeIcon className="h-4 w-4 text-gray-400" />
                  </div>
                  <input
                    id="email"
                    type="email"
                    placeholder="请输入邮箱"
                    className="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700/50 dark:text-white bg-white/70"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    disabled={loading}
                  />
                </div>
              </div>

              <div>
                <div className="flex justify-between items-center mb-1">
                  <label
                    className="block text-gray-700 dark:text-gray-300 text-xs font-medium"
                    htmlFor="password"
                  >
                    密码
                  </label>
                  <button
                    type="button"
                    className="text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400 font-medium"
                    onClick={onForgotPasswordClick}
                  >
                    忘记密码?
                  </button>
                </div>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <LockClosedIcon className="h-4 w-4 text-gray-400" />
                  </div>
                  <input
                    id="password"
                    type="password"
                    placeholder="请输入密码"
                    className="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700/50 dark:text-white bg-white/70"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    disabled={loading}
                  />
                </div>
              </div>

              <div className="flex items-center">
                <input
                  id="remember"
                  type="checkbox"
                  className="h-3 w-3 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label
                  htmlFor="remember"
                  className="ml-2 block text-xs text-gray-700 dark:text-gray-300"
                >
                  记住我
                </label>
              </div>

              <button
                type="submit"
                className="w-full bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-medium py-2 px-4 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 transition-all duration-300 flex items-center justify-center disabled:opacity-70 shadow-md hover:shadow-lg text-sm mt-2"
                disabled={loading}
              >
                {loading ? (
                  <div className="flex items-center">
                    <svg
                      className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      />
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      />
                    </svg>
                    登录中...
                  </div>
                ) : (
                  <div className="flex items-center">
                    登录
                    <ArrowRightIcon className="ml-1 h-4 w-4" />
                  </div>
                )}
              </button>
            </form>

            <div className="mt-4 text-center">
              <p className="text-xs text-gray-600 dark:text-gray-400">
                还没有账号？
                <button
                  onClick={onRegisterClick}
                  className="ml-1 text-blue-600 hover:text-blue-700 font-medium focus:outline-none transition-colors"
                  disabled={loading}
                >
                  立即注册
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
