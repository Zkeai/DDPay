"use client";

import React from "react";
import {
  ArrowUpIcon,
  ArrowDownIcon,
  UserGroupIcon,
  ShoppingBagIcon,
  CreditCardIcon,
  BuildingStorefrontIcon,
  BellIcon,
  ChartBarIcon,
  ClipboardDocumentListIcon,
} from "@heroicons/react/24/outline";

const DashboardPage = () => {
  // 模拟数据
  const stats = [
    {
      title: "总用户数",
      value: "12,345",
      change: "+12%",
      isIncrease: true,
      icon: <UserGroupIcon className="h-6 w-6" />,
    },
    {
      title: "总订单数",
      value: "8,642",
      change: "+8%",
      isIncrease: true,
      icon: <ShoppingBagIcon className="h-6 w-6" />,
    },
    {
      title: "总收入",
      value: "¥128,456",
      change: "+15%",
      isIncrease: true,
      icon: <CreditCardIcon className="h-6 w-6" />,
    },
    {
      title: "活跃商店",
      value: "126",
      change: "-3%",
      isIncrease: false,
      icon: <BuildingStorefrontIcon className="h-6 w-6" />,
    },
  ];

  // 最近订单数据
  const recentOrders = [
    {
      id: "ORD-001",
      customer: "张三",
      amount: "¥256.00",
      status: "已完成",
      date: "2023-10-12",
    },
    {
      id: "ORD-002",
      customer: "李四",
      amount: "¥1,200.00",
      status: "处理中",
      date: "2023-10-11",
    },
    {
      id: "ORD-003",
      customer: "王五",
      amount: "¥658.50",
      status: "已完成",
      date: "2023-10-10",
    },
    {
      id: "ORD-004",
      customer: "赵六",
      amount: "¥99.99",
      status: "已取消",
      date: "2023-10-09",
    },
    {
      id: "ORD-005",
      customer: "孙七",
      amount: "¥3,250.00",
      status: "已完成",
      date: "2023-10-08",
    },
  ];

  // 公告数据 - JSON格式
  const announcements = [
    {
      id: 1,
      type: "官方",
      content:
        "免备案香港/美国大宽带服务器，建站首选，售后无忧直接找官方，支持TG/Discord/QQ三个渠道售后！无脑直接入！",
      timestamp: "2024-11-03 12:28:42",
      color: "blue",
    },
    {
      id: 2,
      type: "广告",
      content:
        "[白舟智选]: 0代理费，1元起充，支持电商自动发货。企业级管理平台，提供多样数字权益和优质虚拟货源，安全高效，一站式服务！",
      timestamp: "2024-10-30 23:07:39",
      color: "green",
    },
    {
      id: 3,
      type: "通知",
      content:
        "DDPay系统将于2023年11月15日凌晨2点进行例行维护，预计维护时间2小时。请提前做好相关准备。",
      timestamp: "2024-10-28 09:15:22",
      color: "yellow",
    },
  ];

  return (
    // 响应式布局容器
    <div className="flex flex-col lg:flex-row gap-6 p-4 md:p-6">
      {/* 系统公告 - 移动端顶部/桌面端左侧 */}
      <div className="w-full lg:w-1/3 lg:sticky lg:top-6 lg:self-start">
        <div className="rounded-lg bg-white p-4 md:p-6 shadow-sm dark:bg-gray-800 border border-gray-100 dark:border-gray-700 h-full">
          <h2 className="mb-4 text-lg font-semibold text-gray-900 dark:text-white flex items-center">
            <BellIcon className="w-5 h-5 mr-2 text-blue-500" />
            系统公告
          </h2>
          <div className="space-y-4">
            {announcements.map((announcement) => (
              <div
                key={announcement.id}
                className={`rounded-lg border border-${announcement.color}-100 bg-${announcement.color}-50 p-4 dark:border-${announcement.color}-900/30 dark:bg-${announcement.color}-900/20`}
              >
                <div className="flex items-start">
                  <span
                    className={`inline-flex px-2 py-0.5 rounded-md text-xs font-semibold mr-2 bg-${announcement.color}-200 text-${announcement.color}-800 dark:bg-${announcement.color}-900/40 dark:text-${announcement.color}-300`}
                  >
                    {announcement.type}
                  </span>
                  <p
                    className={`text-sm text-${announcement.color}-700 dark:text-${announcement.color}-300 flex-1`}
                  >
                    {announcement.content}
                  </p>
                </div>
                <div
                  className={`mt-2 text-xs text-${announcement.color}-600 dark:text-${announcement.color}-400 text-right`}
                >
                  {announcement.timestamp}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* 右侧内容区域 */}
      <div className="w-full lg:w-2/3 space-y-6">
        {/* 统计卡片 */}
        <div className="rounded-lg bg-white p-4 md:p-6 shadow-sm dark:bg-gray-800 border border-gray-100 dark:border-gray-700">
          <h2 className="mb-4 text-lg font-semibold text-gray-900 dark:text-white flex items-center">
            <ChartBarIcon className="w-5 h-5 mr-2 text-blue-500" />
            数据概览
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
            {stats.map((stat, index) => (
              <div
                key={index}
                className="rounded-lg bg-white p-4 shadow-sm dark:bg-gray-800 border border-gray-100 dark:border-gray-700"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-xs sm:text-sm font-medium text-gray-600 dark:text-gray-400">
                      {stat.title}
                    </p>
                    <p className="text-xl sm:text-2xl font-semibold text-gray-900 dark:text-white">
                      {stat.value}
                    </p>
                  </div>
                  <div className="rounded-full bg-blue-50 p-2 sm:p-3 dark:bg-blue-900/20">
                    {stat.icon}
                  </div>
                </div>
                <div className="mt-3 flex items-center">
                  {stat.isIncrease ? (
                    <ArrowUpIcon className="h-3 w-3 sm:h-4 sm:w-4 text-green-500" />
                  ) : (
                    <ArrowDownIcon className="h-3 w-3 sm:h-4 sm:w-4 text-red-500" />
                  )}
                  <span
                    className={`ml-1 text-xs sm:text-sm font-medium ${
                      stat.isIncrease ? "text-green-500" : "text-red-500"
                    }`}
                  >
                    {stat.change}
                  </span>
                  <span className="ml-2 text-xs sm:text-sm text-gray-600 dark:text-gray-400">
                    与上月相比
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 最近订单 */}
        <div className="rounded-lg bg-white p-4 md:p-6 shadow-sm dark:bg-gray-800 border border-gray-100 dark:border-gray-700">
          <h2 className="mb-4 text-lg font-semibold text-gray-900 dark:text-white flex items-center">
            <ClipboardDocumentListIcon className="w-5 h-5 mr-2 text-blue-500" />
            最近订单
          </h2>
          <div className="overflow-x-auto">
            <div className="min-w-full rounded-lg border border-gray-200 dark:border-gray-700">
              <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead className="bg-gray-50 dark:bg-gray-700">
                  <tr>
                    <th
                      scope="col"
                      className="px-3 py-2 md:px-6 md:py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    >
                      订单号
                    </th>
                    <th
                      scope="col"
                      className="px-3 py-2 md:px-6 md:py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    >
                      客户
                    </th>
                    <th
                      scope="col"
                      className="px-3 py-2 md:px-6 md:py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    >
                      金额
                    </th>
                    <th
                      scope="col"
                      className="px-3 py-2 md:px-6 md:py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    >
                      状态
                    </th>
                    <th
                      scope="col"
                      className="px-3 py-2 md:px-6 md:py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
                    >
                      日期
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 bg-white dark:divide-gray-700 dark:bg-gray-800">
                  {recentOrders.map((order) => (
                    <tr key={order.id}>
                      <td className="whitespace-nowrap px-3 py-2 md:px-6 md:py-4">
                        <div className="text-xs md:text-sm font-medium text-blue-600 dark:text-blue-400">
                          {order.id}
                        </div>
                      </td>
                      <td className="whitespace-nowrap px-3 py-2 md:px-6 md:py-4">
                        <div className="text-xs md:text-sm text-gray-900 dark:text-white">
                          {order.customer}
                        </div>
                      </td>
                      <td className="whitespace-nowrap px-3 py-2 md:px-6 md:py-4">
                        <div className="text-xs md:text-sm text-gray-900 dark:text-white">
                          {order.amount}
                        </div>
                      </td>
                      <td className="whitespace-nowrap px-3 py-2 md:px-6 md:py-4">
                        <span
                          className={`inline-flex rounded-full px-2 text-xs font-semibold leading-5 ${
                            order.status === "已完成"
                              ? "bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400"
                              : order.status === "处理中"
                                ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400"
                                : "bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400"
                          }`}
                        >
                          {order.status}
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-3 py-2 md:px-6 md:py-4 text-xs md:text-sm text-gray-500 dark:text-gray-400">
                        {order.date}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
