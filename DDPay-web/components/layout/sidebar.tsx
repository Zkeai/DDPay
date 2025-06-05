"use client";

import type { SidebarItem, SidebarSubItem } from "@/types";

import { Image } from "@heroui/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState, useEffect } from "react";
import { useTheme } from "next-themes";
import {
  ChevronDownIcon,
  ChevronRightIcon,
  ArrowRightOnRectangleIcon,
} from "@heroicons/react/24/outline";

interface SidebarCategory {
  name: string;
  items: SidebarItem[];
}

interface SidebarProps {
  items: SidebarItem[];
  categories?: SidebarCategory[];
  collapsed: boolean;
  userName?: string;
  avatarUrl?: string;
  onLogout?: () => void;
  openMenus: { [key: string]: boolean };
  onToggleMenu: (title: string) => void;
}

export default function Sidebar({
  items,
  categories = [],
  collapsed,
  userName = "用户",
  avatarUrl = "https://via.placeholder.com/40",
  onLogout = () => {},
  openMenus,
  onToggleMenu,
}: SidebarProps) {
  const pathname = usePathname();
  const [mounted, setMounted] = useState(false);
  const { theme } = useTheme();

  // 在客户端渲染后设置mounted为true
  useEffect(() => {
    setMounted(true);
  }, []);

  const isActive = (href: string) => {
    return pathname === href || pathname.startsWith(href + "/");
  };

  // 处理未分类的菜单项
  const uncategorizedItems =
    categories.length > 0
      ? items.filter(
          (item) => !categories.some((cat) => cat.items.includes(item))
        )
      : items;

  // 渲染菜单项
  const renderMenuItem = (item: SidebarItem) => {
    const active = item.href
      ? isActive(item.href)
      : item.children?.some((sub) => isActive(sub.href)) || false;

    if (item.href) {
      return (
        <Link
          href={item.href}
          className={`flex items-center px-3 py-2 rounded-md ${
            active && mounted
              ? "bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
              : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800"
          }`}
        >
          <span className="flex items-center gap-3">
            <span className="text-lg w-5 h-5 flex items-center justify-center">
              {item.icon}
            </span>
            {!collapsed && <span>{item.title}</span>}
          </span>
        </Link>
      );
    }

    return (
      <>
        <button
          className={`flex items-center justify-between w-full px-3 py-2 text-left rounded-md ${
            active && mounted
              ? "bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
              : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800"
          }`}
          onClick={() => onToggleMenu(item.title)}
          aria-expanded={mounted ? !!openMenus[item.title] : false}
        >
          <span className="flex items-center gap-3">
            <span className="text-lg w-5 h-5 flex items-center justify-center">
              {item.icon}
            </span>
            {!collapsed && <span>{item.title}</span>}
          </span>
          {!collapsed &&
            item.children &&
            (mounted && openMenus[item.title] ? (
              <ChevronDownIcon className="w-4 h-4 text-gray-400" />
            ) : (
              <ChevronRightIcon className="w-4 h-4 text-gray-400" />
            ))}
        </button>
        {!collapsed && item.children && mounted && openMenus[item.title] && (
          <ul className="ml-7 mt-1 space-y-1 text-sm">
            {item.children.map((subItem: SidebarSubItem) => (
              <li key={subItem.href}>
                <Link
                  className={`flex items-center px-3 py-1.5 rounded-md ${
                    isActive(subItem.href) && mounted
                      ? "bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
                      : "text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
                  }`}
                  href={subItem.href}
                >
                  {subItem.icon && (
                    <span className="w-4 h-4 mr-2 flex items-center justify-center text-xs">
                      {subItem.icon}
                    </span>
                  )}
                  <span>{subItem.title}</span>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </>
    );
  };

  // 初始渲染时使用统一样式以避免hydration不匹配
  if (!mounted) {
    return (
      <aside
        className={`h-full dark:bg-gray-900 bg-white dark:border-gray-800 border-gray-200 border-r transition-[width] duration-300 flex flex-col ${
          collapsed ? "w-0 overflow-hidden" : "w-64"
        }`}
      >
        <div className="sticky top-0 z-10">
          <div className="flex justify-center items-center py-4 dark:bg-gray-950 bg-gray-50">
            <h1
              className={`font-bold ${collapsed ? "text-sm" : "text-xl"} dark:text-blue-400 text-blue-600`}
            >
              {collapsed ? "DD" : "DDPay"}
            </h1>
          </div>
          <div className="flex flex-col items-center py-4 dark:bg-gray-900 bg-white border-b dark:border-gray-800 border-gray-200">
            <div
              className={`${collapsed ? "w-10 h-10" : "w-16 h-16"} rounded-full overflow-hidden mb-3`}
            >
              <Image
                alt="用户头像"
                className="w-full h-full object-cover"
                src={avatarUrl}
              />
            </div>
            {!collapsed && (
              <div className="flex items-center justify-center w-full px-2">
                <span className="text-sm font-medium dark:text-white text-gray-900">
                  {userName}
                </span>
              </div>
            )}
          </div>
        </div>
        <nav className="px-3 py-4 space-y-6 flex-1 overflow-y-auto">
          {/* 初始渲染时的基本结构 */}
        </nav>
      </aside>
    );
  }

  return (
    <aside
      className={`h-full dark:bg-gray-900 bg-white dark:border-gray-800 border-gray-200 border-r transition-[width] duration-300 flex flex-col ${
        collapsed ? "w-0 overflow-hidden" : "w-64"
      }`}
    >
      {/* 固定头部区域 */}
      <div className="sticky top-0 z-10">
        {/* 顶部标题 */}
        <div className="flex justify-center items-center py-4 dark:bg-gray-950 bg-gray-50">
          <h1
            className={`font-bold ${collapsed ? "text-sm" : "text-xl"} dark:text-blue-400 text-blue-600`}
          >
            {collapsed ? "DD" : "DDPay"}
          </h1>
        </div>

        {/* 用户头像模块 */}
        <div className="flex flex-col items-center py-4 dark:bg-gray-900 bg-white border-b dark:border-gray-800 border-gray-200">
          <div
            className={`${collapsed ? "w-10 h-10" : "w-16 h-16"} rounded-full overflow-hidden mb-3`}
          >
            <Image
              alt="用户头像"
              className="w-full h-full object-cover"
              src={avatarUrl}
            />
          </div>
          {!collapsed && (
            <div className="flex items-center justify-center w-full px-2">
              <span className="text-sm font-medium dark:text-white text-gray-900">
                {userName}
              </span>
              <button
                className="dark:text-gray-400 text-gray-500 dark:hover:text-red-400 hover:text-red-500 ml-2"
                title="登出"
                onClick={onLogout}
                aria-label="登出"
              >
                <ArrowRightOnRectangleIcon className="w-5 h-5" />
              </button>
            </div>
          )}
        </div>
      </div>

      {/* 菜单导航 - 可滚动部分 */}
      <nav className="px-3 py-4 space-y-6 flex-1 overflow-y-auto">
        {/* 分类菜单 */}
        {categories.length > 0 &&
          categories.map((category) => (
            <div key={category.name} className="space-y-2">
              {!collapsed && (
                <div className="pl-2 mb-2">
                  <span className="text-xs font-semibold tracking-wider uppercase dark:text-gray-500 text-gray-400">
                    {category.name}
                  </span>
                </div>
              )}

              {category.items.map((item) => (
                <div key={item.title}>{renderMenuItem(item)}</div>
              ))}
            </div>
          ))}

        {/* 未分类菜单项 */}
        {uncategorizedItems.length > 0 && (
          <div className="space-y-2">
            {categories.length > 0 && !collapsed && (
              <div className="pl-2 mb-2">
                <span className="text-xs font-semibold tracking-wider uppercase dark:text-gray-500 text-gray-400">
                  其他
                </span>
              </div>
            )}

            {uncategorizedItems.map((item) => (
              <div key={item.title}>{renderMenuItem(item)}</div>
            ))}
          </div>
        )}
      </nav>
    </aside>
  );
}
