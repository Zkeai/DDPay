"use client";

import { useState, useEffect } from "react";
import {
  MegaphoneIcon,
  ArrowUpIcon,
  UserGroupIcon,
  CurrencyDollarIcon,
  ShoppingCartIcon,
} from "@heroicons/react/24/outline";

// 公告类型定义
interface Notice {
  id: number;
  title: string;
  content: string;
  created_at: string;
}

// 统计数据类型定义
interface Stats {
  total_users: number;
  total_orders: number;
  today_orders: number;
  today_amount: number;
  yesterday_amount: number;
  growth_rate: number;
}

// 格式化日期
const formatDate = (dateString: string) => {
  try {
    const date = new Date(dateString);

    return date.toLocaleString("zh-CN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch (error) {
    return dateString;
  }
};

// 获取系统公告
const fetchNotices = async (): Promise<Notice[]> => {
  try {
    // 这里应该调用实际的API，目前使用模拟数据
    // const response = await afetch<{code: number; msg: string; data: Notice[]}>('/api/v1/notices', {
    //   method: 'GET',
    // });
    // return response.data || [];

    // 模拟数据
    return [
      {
        id: 1,
        title: "﹤官方﹥免备案香港/美国大宽带服务器",
        content:
          "免备案香港/美国大宽带服务器，建站首选，售后无忧直接找官方，支持TG/Discord/QQ三个渠道售后！无脑直接入！",
        created_at: "2024-11-03 12:28:42",
      },
      {
        id: 2,
        title: "系统更新通知",
        content:
          "系统将于今晚22:00-23:00进行例行维护，期间可能会出现短暂服务中断，请做好相关准备。",
        created_at: "2024-11-02 15:30:00",
      },
      {
        id: 3,
        title: "新功能上线公告",
        content:
          "我们很高兴地宣布，全新的多链支付功能已经上线，现在支持ETH、BSC、TRON等多条链的支付。",
        created_at: "2024-11-01 10:15:00",
      },
    ];
  } catch (error) {
    return [];
  }
};

// 获取统计数据
const fetchStats = async (): Promise<Stats> => {
  try {
    // 实际开发中应调用API
    // const response = await afetch<{code: number; msg: string; data: Stats}>('/api/v1/dashboard/stats', {
    //   method: 'GET',
    // });
    // return response.data;

    // 模拟数据
    return {
      total_users: 1286,
      total_orders: 5932,
      today_orders: 142,
      today_amount: 12683.25,
      yesterday_amount: 10245.8,
      growth_rate: 23.79,
    };
  } catch (error) {
    return {
      total_users: 0,
      total_orders: 0,
      today_orders: 0,
      today_amount: 0,
      yesterday_amount: 0,
      growth_rate: 0,
    };
  }
};

const DashboardPage = () => {
  const [notices, setNotices] = useState<Notice[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const [noticesData, statsData] = await Promise.all([
          fetchNotices(),
          fetchStats(),
        ]);

        setNotices(noticesData);
        setStats(statsData);
      } catch (error) {
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  return (
    <div className="container mx-auto p-2 md:p-6">
      <h1 className="text-xl md:text-2xl font-bold mb-4 md:mb-6 text-gray-800 dark:text-white">
        管理仪表盘
      </h1>

      {loading ? (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500" />
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-6">
          {/* 左侧公告部分 */}
          <div className="lg:col-span-1">
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 mb-4">
              <div className="flex items-center mb-4">
                <MegaphoneIcon className="h-6 w-6 text-blue-500 mr-2" />
                <h2 className="text-lg font-semibold text-gray-800 dark:text-white">
                  系统公告
                </h2>
              </div>

              <div className="space-y-4">
                {notices.length === 0 ? (
                  <p className="text-gray-500 dark:text-gray-400 text-center py-6">
                    暂无公告
                  </p>
                ) : (
                  notices.map((notice) => (
                    <div
                      key={notice.id}
                      className="border-b border-gray-200 dark:border-gray-700 pb-4 last:border-0 last:pb-0"
                    >
                      <h3 className="font-medium text-gray-800 dark:text-white mb-1">
                        {notice.title}
                      </h3>
                      <p className="text-gray-600 dark:text-gray-300 text-sm mb-2">
                        {notice.content}
                      </p>
                      <p className="text-xs text-gray-400 dark:text-gray-500">
                        {formatDate(notice.created_at)}
                      </p>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>

          {/* 右侧统计数据部分 */}
          <div className="lg:col-span-2">
            {stats && (
              <>
                {/* 概览卡片 */}
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 mb-4">
                  <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-4">
                    数据概览
                  </h2>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* 今日交易额 */}
                    <div className="bg-blue-50 dark:bg-blue-900/30 rounded-lg p-4">
                      <div className="flex justify-between items-start">
                        <div>
                          <p className="text-sm text-gray-500 dark:text-gray-400">
                            今日交易额
                          </p>
                          <p className="text-2xl font-bold text-gray-800 dark:text-white">
                            ${stats.today_amount.toFixed(2)}
                          </p>
                          <div className="flex items-center mt-1">
                            <ArrowUpIcon className="h-3 w-3 text-green-500 mr-1" />
                            <span className="text-xs text-green-500">
                              {stats.growth_rate.toFixed(2)}%
                            </span>
                            <span className="text-xs text-gray-400 dark:text-gray-500 ml-1">
                              较昨日
                            </span>
                          </div>
                        </div>
                        <div className="bg-blue-500 rounded-full p-2">
                          <CurrencyDollarIcon className="h-6 w-6 text-white" />
                        </div>
                      </div>
                    </div>

                    {/* 今日订单数 */}
                    <div className="bg-purple-50 dark:bg-purple-900/30 rounded-lg p-4">
                      <div className="flex justify-between items-start">
                        <div>
                          <p className="text-sm text-gray-500 dark:text-gray-400">
                            今日订单数
                          </p>
                          <p className="text-2xl font-bold text-gray-800 dark:text-white">
                            {stats.today_orders}
                          </p>
                        </div>
                        <div className="bg-purple-500 rounded-full p-2">
                          <ShoppingCartIcon className="h-6 w-6 text-white" />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* 总计数据 */}
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
                  <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-4">
                    总计数据
                  </h2>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* 总用户数 */}
                    <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                      <div className="flex items-center">
                        <div className="bg-gray-100 dark:bg-gray-700 rounded-full p-2 mr-3">
                          <UserGroupIcon className="h-5 w-5 text-gray-500 dark:text-gray-400" />
                        </div>
                        <div>
                          <p className="text-sm text-gray-500 dark:text-gray-400">
                            总用户数
                          </p>
                          <p className="text-xl font-semibold text-gray-800 dark:text-white">
                            {stats.total_users.toLocaleString()}
                          </p>
                        </div>
                      </div>
                    </div>

                    {/* 总订单数 */}
                    <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                      <div className="flex items-center">
                        <div className="bg-gray-100 dark:bg-gray-700 rounded-full p-2 mr-3">
                          <ShoppingCartIcon className="h-5 w-5 text-gray-500 dark:text-gray-400" />
                        </div>
                        <div>
                          <p className="text-sm text-gray-500 dark:text-gray-400">
                            总订单数
                          </p>
                          <p className="text-xl font-semibold text-gray-800 dark:text-white">
                            {stats.total_orders.toLocaleString()}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default DashboardPage;
