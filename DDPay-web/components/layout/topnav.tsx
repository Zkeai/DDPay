"use client";

import React, { useEffect, useState, useRef } from "react";
import {
  ChevronDoubleLeftIcon,
  ChevronDoubleRightIcon,
  GlobeAltIcon,
  SunIcon,
  MoonIcon,
  ChevronDownIcon,
  KeyIcon,
  ClockIcon,
  ArrowRightOnRectangleIcon,
  UserIcon,
} from "@heroicons/react/24/outline";
import { useRouter } from "next/navigation";
import { Image } from "@heroui/image";

import { useTitle } from "@/components/TitleContext";
import { logout } from "@/lib/afetch";

interface TopNavProps {
  collapsed: boolean;
  onToggleCollapse: (collapsed: boolean) => void;
  currentRoute?: string;
  userName?: string;
  avatarUrl?: string;
  onLogout?: () => void;
}

export default function TopNav({
  collapsed,
  onToggleCollapse,
  userName = "管理员",
  avatarUrl,
  onLogout = async () => {
    await logout();

    localStorage.removeItem("access_token");
    localStorage.removeItem("token");
    localStorage.removeItem("user");

    window.location.reload();
  },
}: TopNavProps) {
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const userMenuRef = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const { title } = useTitle();

  useEffect(() => {
    // 检查系统偏好或本地存储的主题设置
    const isDark =
      localStorage.getItem("theme") === "dark" ||
      (!localStorage.getItem("theme") &&
        window.matchMedia("(prefers-color-scheme: dark)").matches);

    setIsDarkMode(isDark);

    if (isDark) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }, []);

  // 点击外部关闭用户菜单
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        userMenuRef.current &&
        !userMenuRef.current.contains(event.target as Node)
      ) {
        setUserMenuOpen(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const toggleSidebar = () => {
    onToggleCollapse(!collapsed);
  };

  const toggleTheme = () => {
    // 添加类禁用过渡动画
    document.documentElement.classList.add("disable-transitions");

    const newTheme = !isDarkMode;

    setIsDarkMode(newTheme);

    if (newTheme) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("theme", "dark");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("theme", "light");
    }

    // 短暂延迟后移除禁用过渡的类
    setTimeout(() => {
      document.documentElement.classList.remove("disable-transitions");
    }, 100);
  };

  const toggleUserMenu = () => {
    setUserMenuOpen(!userMenuOpen);
  };

  const handleLoginLog = () => {
    router.push("/admin/profile/logs");
    setUserMenuOpen(false);
  };

  const handleSecuritySettings = () => {
    router.push("/admin/profile/security");
    setUserMenuOpen(false);
  };

  const handleLogout = () => {
    onLogout();
    setUserMenuOpen(false);
  };

  return (
    <div className="flex flex-col w-full">
      {/* 顶部导航栏 */}
      <div className="h-16 dark:bg-gray-900 bg-white dark:border-gray-800 border-gray-200 border-b flex items-center justify-between px-6 w-full">
        <div className="flex items-center">
          <button
            onClick={toggleSidebar}
            className="dark:text-gray-400 text-gray-500 dark:hover:text-gray-300 hover:text-gray-700 dark:bg-gray-800 bg-gray-100 p-2 rounded-full mr-4 z-10"
            title={collapsed ? "展开侧边栏" : "折叠侧边栏"}
          >
            {collapsed ? (
              <ChevronDoubleRightIcon className="w-5 h-5" />
            ) : (
              <ChevronDoubleLeftIcon className="w-5 h-5" />
            )}
          </button>
        </div>

        <div className="flex items-center gap-4">
          <button
            onClick={toggleTheme}
            className="dark:text-gray-400 text-gray-500 dark:hover:text-gray-300 hover:text-gray-700 dark:bg-gray-800 bg-gray-100 p-2 rounded-full"
            title={isDarkMode ? "切换到亮色模式" : "切换到暗色模式"}
          >
            {isDarkMode ? (
              <SunIcon className="w-5 h-5" />
            ) : (
              <MoonIcon className="w-5 h-5" />
            )}
          </button>

          {/* 用户下拉菜单 */}
          <div className="relative" ref={userMenuRef}>
            <button
              className="flex items-center rounded-md px-3 py-2 dark:bg-gray-800 bg-gray-100 dark:text-gray-300 text-gray-700 hover:bg-gray-200 dark:hover:bg-gray-700"
              onClick={toggleUserMenu}
              aria-expanded={userMenuOpen}
              aria-haspopup="true"
            >
              {avatarUrl ? (
                <div className="w-8 h-8 rounded-full overflow-hidden mr-2">
                  <Image
                    src={avatarUrl}
                    alt="User Avatar"
                    className="w-full h-full object-cover"
                  />
                </div>
              ) : (
                <UserIcon className="w-5 h-5 mr-2" />
              )}
              <span className="font-medium">{userName}</span>
              <ChevronDownIcon className="w-4 h-4 ml-2" />
            </button>

            {userMenuOpen && (
              <div className="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 z-20">
                <div className="py-1" role="menu" aria-orientation="vertical">
                  <button
                    className="flex items-center w-full px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                    onClick={handleLoginLog}
                    role="menuitem"
                  >
                    <ClockIcon className="w-5 h-5 mr-3 text-gray-500 dark:text-gray-400" />
                    <span>登录日志</span>
                  </button>
                  <button
                    className="flex items-center w-full px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                    onClick={handleSecuritySettings}
                    role="menuitem"
                  >
                    <KeyIcon className="w-5 h-5 mr-3 text-gray-500 dark:text-gray-400" />
                    <span>安全设置</span>
                  </button>
                  <div className="border-t border-gray-200 dark:border-gray-700 my-1" />
                  <button
                    className="flex items-center w-full px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700"
                    onClick={handleLogout}
                    role="menuitem"
                  >
                    <ArrowRightOnRectangleIcon className="w-5 h-5 mr-3" />
                    <span>注销登录</span>
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* 路径导航区域 */}
      <div className="py-4 px-6 dark:bg-gray-800 bg-gray-50 dark:border-gray-700 border-gray-200 border-b flex items-center">
        <div className="flex items-center gap-2 dark:text-gray-300 text-gray-700">
          <GlobeAltIcon className="w-5 h-5" />
          <span className="font-medium">{title}</span>
        </div>
      </div>
    </div>
  );
}
