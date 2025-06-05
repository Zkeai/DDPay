"use client";

import type { CalendarDate } from "@internationalized/date";

import { useState, useEffect } from "react";
import {
  ArrowLeftIcon,
  ArrowRightIcon,
  MagnifyingGlassIcon,
  XMarkIcon,
} from "@heroicons/react/24/outline";
import { DateRangePicker } from "@heroui/date-picker";

import { useAuthStore } from "@/store/auth";
import { getLoginLogs } from "@/lib/afetch";

// 自定义日期范围类型
type DateRange = {
  start: CalendarDate;
  end: CalendarDate;
} | null;

// 状态标签组件
const StatusBadge = ({ status }: { status: number }) => {
  return status === 1 ? (
    <span className="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
      成功
    </span>
  ) : (
    <span className="px-2 py-1 text-xs rounded-full bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200">
      失败
    </span>
  );
};

// 日期格式化，在移动端显示更简洁的格式
const formatDate = (dateString: string) => {
  try {
    const date = new Date(dateString);

    // 在小屏幕设备上使用更简洁的日期格式
    const isMobile = window.innerWidth < 768;

    if (isMobile) {
      return date.toLocaleString("zh-CN", {
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
      });
    }

    return date.toLocaleString("zh-CN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  } catch (error) {
    return dateString;
  }
};

// 格式化IP地址，将::1替换为127.0.0.1
const formatIP = (ip: string) => {
  return ip === "::1" ? "127.0.0.1" : ip;
};

// 格式化User Agent，显示完整的UA
const formatUserAgent = (ua: string) => {
  if (!ua) return "-";

  return ua;
};

interface LoginLog {
  id: number;
  user_id: number;
  login_type: string;
  ip: string;
  user_agent: string;
  status: number;
  fail_reason: string;
  created_at: string;
}

export default function LoginLogsPage() {
  const { user } = useAuthStore();
  const [logs, setLogs] = useState<LoginLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(1);

  const [value, setValue] = useState<DateRange>(null);

  // 筛选条件
  const [ip, setIp] = useState("");
  const [statusFilter, setStatusFilter] = useState<number | undefined>(
    undefined
  );

  // 加载数据
  const loadLogs = async () => {
    if (!user) return;

    setLoading(true);
    try {
      const params: any = {
        user_id: user.id,
        page,
        page_size: pageSize,
      };

      if (ip) params.ip = ip;
      if (value?.start) {
        const startDate = new Date(
          value.start.year,
          value.start.month - 1,
          value.start.day
        );

        params.start_time = startDate.toISOString();
      }
      if (value?.end) {
        // 设置结束日期为当天的23:59:59，确保包含整个结束日期
        const endDate = new Date(
          value.end.year,
          value.end.month - 1,
          value.end.day
        );

        endDate.setHours(23, 59, 59, 999);
        params.end_time = endDate.toISOString();
      }
      if (statusFilter !== undefined) params.status = statusFilter;

      const data = await getLoginLogs(params);

      // 确保返回的数据有效
      if (data && data.logs) {
        setLogs(data.logs);
        setTotal(data.total || 0);
        setTotalPages(data.total_pages || 1);
      } else {
        // 如果数据无效，设置为默认值
        setLogs([]);
        setTotal(0);
        setTotalPages(1);
      }
    } catch (error) {
      // 出错时设置为默认值
      setLogs([]);
      setTotal(0);
      setTotalPages(1);
    } finally {
      setLoading(false);
    }
  };

  // 首次加载和筛选条件/分页变化时重新加载
  useEffect(() => {
    loadLogs();
  }, [user, page, pageSize]);

  // 应用筛选
  const applyFilters = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1); // 重置到第一页
    loadLogs();
  };

  // 重置筛选
  const resetFilters = () => {
    setIp("");
    setValue(null);
    setStatusFilter(undefined);
    setPage(1);
    // 通过下一个渲染周期触发loadLogs
    setTimeout(loadLogs, 0);
  };

  return (
    <div className="container mx-auto p-2 md:p-4">
      <h1 className="text-xl md:text-2xl font-bold mb-4 md:mb-6 text-gray-800 dark:text-white">
        登录日志
      </h1>

      {/* 筛选表单 */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-3 md:p-4 mb-4 md:mb-6">
        <form onSubmit={applyFilters} className="space-y-3 md:space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 md:gap-4">
            <div>
              <label
                htmlFor="ip-filter"
                className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              >
                IP地址
              </label>
              <div className="relative">
                <input
                  id="ip-filter"
                  type="text"
                  value={ip}
                  onChange={(e) => setIp(e.target.value)}
                  className="w-full rounded-md border border-gray-300 dark:border-gray-600 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white text-sm"
                  placeholder="搜索IP地址 (::1 等同于 127.0.0.1)"
                />
                <MagnifyingGlassIcon className="absolute right-3 top-2.5 h-4 w-4 text-gray-400" />
              </div>
            </div>

            <div>
              <DateRangePicker
                label="时间范围"
                value={value}
                onChange={(newValue) => {
                  setValue(newValue);
                }}
              />
            </div>
          </div>

          {/* 状态筛选和按钮 */}
          <div className="flex flex-col md:flex-row md:items-center gap-3 md:gap-2">
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                状态:
              </span>
              <button
                type="button"
                onClick={() => setStatusFilter(undefined)}
                className={`px-3 py-1 text-xs rounded-full ${
                  statusFilter === undefined
                    ? "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
                    : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"
                }`}
              >
                全部
              </button>
              <button
                type="button"
                onClick={() => setStatusFilter(1)}
                className={`px-3 py-1 text-xs rounded-full ${
                  statusFilter === 1
                    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                    : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"
                }`}
              >
                成功
              </button>
              <button
                type="button"
                onClick={() => setStatusFilter(0)}
                className={`px-3 py-1 text-xs rounded-full ${
                  statusFilter === 0
                    ? "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
                    : "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"
                }`}
              >
                失败
              </button>
            </div>

            <div className="flex-grow" />

            <div className="flex space-x-2 mt-2 md:mt-0">
              <button
                type="button"
                onClick={resetFilters}
                className="inline-flex items-center px-3 py-1.5 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
              >
                <XMarkIcon className="h-4 w-4 mr-1" />
                重置
              </button>
              <button
                type="submit"
                className="inline-flex items-center px-3 py-1.5 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <MagnifyingGlassIcon className="h-4 w-4 mr-1" />
                查询
              </button>
            </div>
          </div>
        </form>
      </div>

      {/* 数据表格 - 小屏幕使用卡片视图 */}
      <div className="hidden md:block bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th
                  scope="col"
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  时间
                </th>
                <th
                  scope="col"
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  登录类型
                </th>
                <th
                  scope="col"
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  用户代理
                </th>
                <th
                  scope="col"
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  IP地址
                </th>
                <th
                  scope="col"
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  状态
                </th>
                <th
                  scope="col"
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  原因
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {loading ? (
                <tr>
                  <td
                    colSpan={6}
                    className="px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400"
                  >
                    加载中...
                  </td>
                </tr>
              ) : !logs || logs.length === 0 ? (
                <tr>
                  <td
                    colSpan={6}
                    className="px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400"
                  >
                    暂无数据
                  </td>
                </tr>
              ) : (
                logs.map((log) => (
                  <tr
                    key={log.id}
                    className="hover:bg-gray-50 dark:hover:bg-gray-700"
                  >
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {formatDate(log.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {log.login_type}
                    </td>
                    <td className="px-6 py-4 text-xs text-gray-500 dark:text-gray-400 max-w-md break-words whitespace-normal">
                      {formatUserAgent(log.user_agent)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {formatIP(log.ip)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <StatusBadge status={log.status} />
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {log.fail_reason || "-"}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* 移动端卡片视图 */}
      <div className="md:hidden space-y-4">
        {loading ? (
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4 text-center text-gray-500 dark:text-gray-400">
            加载中...
          </div>
        ) : !logs || logs.length === 0 ? (
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-4 text-center text-gray-500 dark:text-gray-400">
            暂无数据
          </div>
        ) : (
          logs.map((log) => (
            <div
              key={log.id}
              className="bg-white dark:bg-gray-800 rounded-lg shadow p-4 space-y-2"
            >
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  {formatDate(log.created_at)}
                </span>
                <StatusBadge status={log.status} />
              </div>

              <div className="grid grid-cols-2 gap-2 text-sm">
                <div>
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    IP地址
                  </span>
                  <div className="font-medium">{formatIP(log.ip)}</div>
                </div>
                <div>
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    登录类型
                  </span>
                  <div className="font-medium">{log.login_type}</div>
                </div>
              </div>

              {log.fail_reason && (
                <div>
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    原因
                  </span>
                  <div className="text-sm font-medium text-red-500">
                    {log.fail_reason}
                  </div>
                </div>
              )}

              <div>
                <span className="text-xs text-gray-500 dark:text-gray-400">
                  用户代理
                </span>
                <div className="text-xs text-gray-600 dark:text-gray-400 break-words">
                  {formatUserAgent(log.user_agent)}
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {/* 分页 */}
      <div className="px-4 py-3 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 rounded-lg shadow-md mt-4">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between space-y-3 md:space-y-0">
          <div className="flex items-center space-x-2">
            <div className="text-sm text-gray-700 dark:text-gray-300">
              共 <span className="font-medium">{total}</span> 条记录
            </div>
            <div className="flex items-center ml-0 md:ml-4">
              <label
                htmlFor="page-size"
                className="text-sm text-gray-700 dark:text-gray-300 mr-2"
              >
                每页:
              </label>
              <select
                id="page-size"
                value={pageSize}
                onChange={(e) => {
                  setPageSize(Number(e.target.value));
                  setPage(1); // 重置到第一页
                }}
                className="rounded-md border border-gray-300 dark:border-gray-600 px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
              >
                <option value="10">10条</option>
                <option value="50">50条</option>
                <option value="100">100条</option>
                <option value="500">500条</option>
              </select>
            </div>
          </div>
          <div className="flex justify-between md:justify-end space-x-2">
            <button
              onClick={() => setPage(Math.max(1, page - 1))}
              disabled={page === 1}
              className={`inline-flex items-center px-3 py-1.5 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md ${
                page === 1
                  ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                  : "text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
              }`}
            >
              <ArrowLeftIcon className="h-4 w-4 mr-1" />
              上一页
            </button>
            <span className="inline-flex items-center px-3 py-1.5 text-sm text-gray-700 dark:text-gray-300">
              {page} / {totalPages}
            </span>
            <button
              onClick={() => setPage(Math.min(totalPages, page + 1))}
              disabled={page === totalPages}
              className={`inline-flex items-center px-3 py-1.5 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md ${
                page === totalPages
                  ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                  : "text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
              }`}
            >
              下一页
              <ArrowRightIcon className="h-4 w-4 ml-1" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
