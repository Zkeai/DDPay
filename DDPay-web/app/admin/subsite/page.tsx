"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import {
  PlusCircleIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon,
} from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";

import { useAuthStore } from "@/store/auth";
import { afetch } from "@/lib/afetch";

// 分站列表类型定义
interface Subsite {
  id: number;
  name: string;
  subdomain: string;
  domain: string;
  status: number;
  commission_rate: number;
}

const SubsitePage = () => {
  const [subsites, setSubsites] = useState<Subsite[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const router = useRouter();
  const { user } = useAuthStore();

  // 判断是否为管理员
  const isAdmin = user?.role === "admin";

  // 使用useCallback包装fetchSubsites函数
  const fetchSubsites = useCallback(async () => {
    setLoading(true);
    try {
      let url = `/api/v1/subsite/list?page=${page}&page_size=${pageSize}`;

      if (searchTerm) {
        url += `&search=${encodeURIComponent(searchTerm)}`;
      }

      console.log("请求URL:", url);

      // 添加超时处理
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 10000); // 10秒超时

      const response = await afetch<{
        code: number;
        msg: string;
        data: {
          subsites: Subsite[];
          total: number;
        } | null;
      }>(url, {
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (response.code === 200) {
        // 处理API返回的数据
        if (response.data && response.data.subsites) {
          setSubsites(response.data.subsites);
          setTotal(response.data.total || 0);
        } else {
          // API返回成功但没有数据
          setSubsites([]);
          setTotal(0);
        }
      } else {
        // API返回错误
        setSubsites([]);
        setTotal(0);
      }
    } catch (error) {
      // 显示更详细的错误信息
      if (error instanceof Error) {
        if (error.name === "AbortError") {
        }
      }
      // 发生错误时设置空数组
      setSubsites([]);
      setTotal(0);

      // 添加错误提示，帮助用户了解问题
      alert(
        `获取分站列表失败: ${error instanceof Error ? error.message : "未知错误"}\n请检查后端API是否正常运行，以及数据库连接是否正常。`
      );
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, searchTerm]);

  // 使用useEffect控制API调用时机
  useEffect(() => {
    fetchSubsites();
  }, [fetchSubsites]);

  // 手动刷新列表
  const refreshList = () => {
    fetchSubsites();
  };

  // 处理搜索
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1); // 重置到第一页
  };

  // 跳转到创建分站页面
  const goToCreateSubsite = () => {
    router.push("/admin/subsite/create");
  };

  // 查看分站详情
  const viewSubsite = (id: number) => {
    router.push(`/admin/subsite/${id}`);
  };

  // 显示用户自己的分站(非管理员)
  const goToMySubsite = async () => {
    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: {
          subsite_info: {
            subsite: {
              id: number;
            };
          };
        };
      }>("/api/v1/subsite/info");

      if (response.code === 200 && response.data.subsite_info) {
        router.push(`/admin/subsite/${response.data.subsite_info.subsite.id}`);
      } else {
        // 用户没有分站，跳转到创建页面
        router.push("/admin/subsite/create");
      }
    } catch (error) {
      // 出错时也跳转到创建页面
      router.push("/admin/subsite/create");
    }
  };

  return (
    <div className="container mx-auto p-4">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-4 gap-4">
        <form
          onSubmit={handleSearch}
          className="flex items-center w-full md:w-auto"
        >
          <div className="relative w-full md:w-64">
            <input
              type="text"
              className="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              placeholder="搜索分站名称或域名..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
            <MagnifyingGlassIcon className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
          </div>
          <Button type="submit" variant="solid" size="md" className="ml-2">
            搜索
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="md"
            className="ml-1"
            onClick={refreshList}
            title="刷新"
          >
            <ArrowPathIcon className="h-5 w-5" />
          </Button>
        </form>

        {isAdmin ? (
          <Button
            variant="solid"
            size="md"
            onClick={goToCreateSubsite}
            className="w-full md:w-auto"
          >
            <PlusCircleIcon className="h-5 w-5 mr-1" /> 创建分站
          </Button>
        ) : (
          <Button
            variant="solid"
            size="md"
            onClick={goToMySubsite}
            className="w-full md:w-auto"
          >
            我的分站
          </Button>
        )}
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4">
        {loading ? (
          <div className="flex justify-center py-10">
            <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500" />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th
                    scope="col"
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                  >
                    ID
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                  >
                    名称
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                  >
                    子域名
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                  >
                    域名
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                  >
                    佣金比例
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                  >
                    状态
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                  >
                    操作
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
                {subsites.length > 0 ? (
                  subsites.map((row) => (
                    <tr
                      key={row.id}
                      className="hover:bg-gray-50 dark:hover:bg-gray-800"
                    >
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                        {row.id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="font-medium text-gray-900 dark:text-gray-100">
                          {row.name}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                        {row.subdomain}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                        {row.domain || "-"}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                        {row.commission_rate}%
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            row.status === 1
                              ? "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400"
                              : "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400"
                          }`}
                        >
                          {row.status === 1 ? "正常" : "禁用"}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <Button
                          size="sm"
                          variant="solid"
                          onClick={() => viewSubsite(row.id)}
                        >
                          查看
                        </Button>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={7} className="px-6 py-10 text-center">
                      <p className="text-lg text-gray-500 dark:text-gray-400 mb-4">
                        没有找到分站数据
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
                        {isAdmin
                          ? "您可以创建第一个分站来开始使用系统"
                          : "请联系管理员或创建您的第一个分站"}
                      </p>
                      <Button
                        variant="solid"
                        size="md"
                        onClick={isAdmin ? goToCreateSubsite : goToMySubsite}
                      >
                        {isAdmin ? "创建分站" : "我的分站"}
                      </Button>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}

        {!loading && total > 0 && (
          <div className="mt-4 flex justify-between items-center">
            <div className="text-sm text-gray-500">共 {total} 条记录</div>
            <div className="flex items-center space-x-2">
              <Button
                size="sm"
                variant="bordered"
                onClick={() => setPage(page > 1 ? page - 1 : 1)}
                disabled={page <= 1}
              >
                上一页
              </Button>
              <span className="text-sm">
                第 {page} 页，共 {Math.ceil(total / pageSize) || 1} 页
              </span>
              <Button
                size="sm"
                variant="bordered"
                onClick={() =>
                  setPage(page < Math.ceil(total / pageSize) ? page + 1 : page)
                }
                disabled={page >= Math.ceil(total / pageSize) || total === 0}
              >
                下一页
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default SubsitePage;
