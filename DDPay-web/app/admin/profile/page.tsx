"use client";

import Link from "next/link";
import {
  ClockIcon,
  UserIcon,
  EnvelopeIcon,
  CalendarIcon,
} from "@heroicons/react/24/outline";
import { Image } from "@heroui/image";

import { useAuthStore } from "@/store/auth";

export default function ProfilePage() {
  const { user } = useAuthStore();

  if (!user) {
    return (
      <div className="flex items-center justify-center h-[300px]">
        <div className="text-center">
          <p className="text-gray-500 dark:text-gray-400">请先登录</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-6 text-gray-800 dark:text-white">
        个人资料
      </h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* 个人信息卡片 */}
        <div className="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <div className="flex flex-col md:flex-row md:items-center">
            <div className="flex-shrink-0 mb-4 md:mb-0 md:mr-6">
              <div className="w-20 h-20 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center">
                {user.avatar ? (
                  <Image
                    src={user.avatar}
                    alt={user.username}
                    className="w-full h-full rounded-full object-cover"
                  />
                ) : (
                  <UserIcon className="w-10 h-10 text-gray-500 dark:text-gray-400" />
                )}
              </div>
            </div>

            <div className="flex-grow">
              <h2 className="text-xl font-bold text-gray-800 dark:text-white">
                {user.username}
              </h2>
              <div className="mt-4 space-y-2">
                <div className="flex items-center text-sm text-gray-600 dark:text-gray-300">
                  <EnvelopeIcon className="h-4 w-4 mr-2" />
                  <span>{user.email}</span>
                </div>

                <div className="flex items-center text-sm text-gray-600 dark:text-gray-300">
                  <CalendarIcon className="h-4 w-4 mr-2" />
                  <span>账户ID: {user.id}</span>
                </div>

                <div className="flex items-center text-sm text-gray-600 dark:text-gray-300">
                  <ClockIcon className="h-4 w-4 mr-2" />
                  <span>角色: {user.role || "普通用户"}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 快捷链接 */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <h3 className="font-medium text-gray-900 dark:text-white mb-4">
            账户管理
          </h3>
          <div className="space-y-2">
            <Link
              href="/admin/profile/logs"
              className="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              登录日志
            </Link>
            {/* 可以添加更多相关链接 */}
          </div>
        </div>
      </div>
    </div>
  );
}
