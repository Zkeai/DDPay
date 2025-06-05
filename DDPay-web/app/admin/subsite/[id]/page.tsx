"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Image } from "@heroui/image";
import Link from "next/link";
import {
  ArrowLeftIcon,
  PencilIcon,
  ShoppingBagIcon,
  ClipboardDocumentListIcon,
  Cog6ToothIcon,
  CurrencyDollarIcon,
} from "@heroicons/react/24/outline";
import { toast } from "react-hot-toast";

import { afetch } from "@/lib/afetch";
import { useAuthStore } from "@/store/auth";

// 分站信息类型定义
interface SubsiteInfo {
  subsite: {
    id: number;
    name: string;
    domain: string;
    subdomain: string;
    logo: string;
    description: string;
    theme: string;
    status: number;
    commission_rate: number;
    created_at: string;
    updated_at: string;
    owner_id: number;
  };
  owner: {
    id: number;
    username: string;
    email: string;
    avatar: string;
  };
  product_count: number;
  order_count: number;
  balance: number;
}

const SubsiteDetailPage = () => {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuthStore();
  const [subsiteInfo, setSubsiteInfo] = useState<SubsiteInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const subsiteId = params.id;

  // 检查当前用户是否为分站所有者或管理员
  const isOwnerOrAdmin = () => {
    if (!subsiteInfo || !user) return false;

    return user.role === "admin" || subsiteInfo.owner.id === user.id;
  };

  useEffect(() => {
    fetchSubsiteInfo();
  }, [subsiteId]);

  const fetchSubsiteInfo = async () => {
    setLoading(true);
    setError("");

    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: {
          subsite_info: SubsiteInfo;
        };
      }>(`/api/v1/subsite/info?id=${subsiteId}`);

      if (response.code === 200) {
        setSubsiteInfo(response.data.subsite_info);
      } else {
        setError(response.msg || "获取分站信息失败");
      }
    } catch (error: any) {
      setError(error.message || "获取分站信息失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  };

  const toggleSubsiteStatus = async () => {
    if (!subsiteInfo) return;

    try {
      const newStatus = subsiteInfo.subsite.status === 1 ? 0 : 1;
      const response = await afetch<{
        code: number;
        msg: string;
      }>(`/api/v1/subsite/update`, {
        method: "PUT",
        body: JSON.stringify({
          id: subsiteInfo.subsite.id,
          status: newStatus,
          name: subsiteInfo.subsite.name,
          subdomain: subsiteInfo.subsite.subdomain,
          commission_rate: subsiteInfo.subsite.commission_rate,
        }),
      });

      if (response.code === 200) {
        toast.success(newStatus === 1 ? "分站已启用" : "分站已禁用");
        // 更新本地状态
        setSubsiteInfo((prev) => {
          if (!prev) return null;

          return {
            ...prev,
            subsite: {
              ...prev.subsite,
              status: newStatus,
            },
          };
        });
      } else {
        toast.error(response.msg || "操作失败");
      }
    } catch (error: any) {
      toast.error(error.message || "操作失败，请稍后重试");
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto p-4 flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500" />
      </div>
    );
  }

  if (error || !subsiteInfo) {
    return (
      <div className="container mx-auto p-4">
        <div className="bg-red-50 dark:bg-red-900/20 p-4 rounded-lg text-red-600 dark:text-red-400">
          <p>{error || "无法获取分站信息"}</p>
          <button
            onClick={() => router.push("/admin/subsite")}
            className="mt-2 text-blue-600 dark:text-blue-400 hover:underline"
          >
            返回分站列表
          </button>
        </div>
      </div>
    );
  }

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

  return (
    <div className="container mx-auto p-4">
      <div className="mb-6 flex items-center justify-between">
        <Link
          href="/admin/subsite"
          className="flex items-center text-blue-600 hover:text-blue-800"
        >
          <ArrowLeftIcon className="h-4 w-4 mr-1" />
          返回分站列表
        </Link>

        {/* 只有管理员和分站所有者可以看到编辑按钮 */}
        {isOwnerOrAdmin() && (
          <Link href={`/admin/subsite/${subsiteId}/edit`}>
            <button className="flex items-center px-3 py-1.5 rounded-lg bg-blue-100 text-blue-600 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400 dark:hover:bg-blue-900/50">
              <PencilIcon className="h-4 w-4 mr-1" />
              编辑分站
            </button>
          </Link>
        )}
      </div>

      {/* 分站基本信息卡片 */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between mb-4">
          <div className="flex items-center">
            {subsiteInfo.subsite.logo ? (
              <Image
                src={subsiteInfo.subsite.logo}
                alt={subsiteInfo.subsite.name}
                className="w-16 h-16 rounded-lg mr-4 object-cover"
              />
            ) : (
              <div className="w-16 h-16 rounded-lg mr-4 bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-gray-500 dark:text-gray-400">
                {(subsiteInfo.subsite.name || "").substring(0, 2).toUpperCase()}
              </div>
            )}
            <div>
              <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
                {subsiteInfo.subsite.name}
              </h1>
              <div className="flex items-center mt-1">
                <span
                  className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    subsiteInfo.subsite.status === 1
                      ? "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400"
                      : "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400"
                  }`}
                >
                  {subsiteInfo.subsite.status === 1 ? "正常" : "已禁用"}
                </span>
                {isOwnerOrAdmin() && (
                  <button
                    onClick={toggleSubsiteStatus}
                    className="ml-2 text-xs text-blue-600 dark:text-blue-400 hover:underline"
                  >
                    {subsiteInfo.subsite.status === 1 ? "禁用" : "启用"}
                  </button>
                )}
              </div>
            </div>
          </div>

          <div className="mt-4 md:mt-0 flex flex-wrap gap-2">
            <a
              href={`https://${subsiteInfo.subsite.subdomain}.ddpay.com`}
              target="_blank"
              rel="noopener noreferrer"
              className="px-3 py-1.5 rounded-lg bg-indigo-100 text-indigo-600 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400 dark:hover:bg-indigo-900/50 text-sm"
            >
              访问分站
            </a>
            {isOwnerOrAdmin() && (
              <>
                <Link
                  href={`/admin/subsite/${subsiteId}/products`}
                  className="px-3 py-1.5 rounded-lg bg-blue-100 text-blue-600 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400 dark:hover:bg-blue-900/50 text-sm flex items-center"
                >
                  <ShoppingBagIcon className="h-4 w-4 mr-1" />
                  商品管理
                </Link>
                <Link
                  href={`/admin/subsite/${subsiteId}/orders`}
                  className="px-3 py-1.5 rounded-lg bg-green-100 text-green-600 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400 dark:hover:bg-green-900/50 text-sm flex items-center"
                >
                  <ClipboardDocumentListIcon className="h-4 w-4 mr-1" />
                  订单管理
                </Link>
                <Link
                  href={`/admin/subsite/${subsiteId}/settings`}
                  className="px-3 py-1.5 rounded-lg bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-400 dark:hover:bg-gray-600 text-sm flex items-center"
                >
                  <Cog6ToothIcon className="h-4 w-4 mr-1" />
                  设置
                </Link>
              </>
            )}
          </div>
        </div>

        <div className="mt-6 border-t border-gray-200 dark:border-gray-700 pt-4">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-3">
            分站信息
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">子域名</p>
              <p className="font-medium">
                {subsiteInfo.subsite.subdomain}.ddpay.com
              </p>
            </div>
            {subsiteInfo.subsite.domain && (
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  自定义域名
                </p>
                <p className="font-medium">{subsiteInfo.subsite.domain}</p>
              </div>
            )}
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">所有者</p>
              <div className="flex items-center">
                {subsiteInfo.owner.avatar ? (
                  <Image
                    src={subsiteInfo.owner.avatar}
                    alt={subsiteInfo.owner.username}
                    className="w-5 h-5 rounded-full mr-2"
                  />
                ) : (
                  <div className="w-5 h-5 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-xs text-gray-500 dark:text-gray-400 mr-2">
                    {(subsiteInfo.owner.username || "")
                      .substring(0, 1)
                      .toUpperCase()}
                  </div>
                )}
                <span className="font-medium">
                  {subsiteInfo.owner.username}
                </span>
              </div>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                佣金比例
              </p>
              <p className="font-medium">
                {subsiteInfo.subsite.commission_rate}%
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">主题</p>
              <p className="font-medium capitalize">
                {subsiteInfo.subsite.theme === "default"
                  ? "默认主题"
                  : subsiteInfo.subsite.theme === "dark"
                    ? "暗黑主题"
                    : subsiteInfo.subsite.theme === "light"
                      ? "明亮主题"
                      : subsiteInfo.subsite.theme === "anime"
                        ? "二次元主题"
                        : subsiteInfo.subsite.theme}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                创建时间
              </p>
              <p className="font-medium">
                {formatDate(subsiteInfo.subsite.created_at)}
              </p>
            </div>
          </div>

          {subsiteInfo.subsite.description && (
            <div className="mt-4">
              <p className="text-sm text-gray-500 dark:text-gray-400">描述</p>
              <p className="mt-1">{subsiteInfo.subsite.description}</p>
            </div>
          )}
        </div>
      </div>

      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 flex items-center">
          <div className="rounded-full bg-blue-100 dark:bg-blue-900/30 p-3 mr-4">
            <ShoppingBagIcon className="h-6 w-6 text-blue-600 dark:text-blue-400" />
          </div>
          <div>
            <p className="text-sm text-gray-500 dark:text-gray-400">商品数量</p>
            <p className="text-2xl font-bold">{subsiteInfo.product_count}</p>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 flex items-center">
          <div className="rounded-full bg-green-100 dark:bg-green-900/30 p-3 mr-4">
            <ClipboardDocumentListIcon className="h-6 w-6 text-green-600 dark:text-green-400" />
          </div>
          <div>
            <p className="text-sm text-gray-500 dark:text-gray-400">订单数量</p>
            <p className="text-2xl font-bold">{subsiteInfo.order_count}</p>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 flex items-center">
          <div className="rounded-full bg-yellow-100 dark:bg-yellow-900/30 p-3 mr-4">
            <CurrencyDollarIcon className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
          </div>
          <div>
            <p className="text-sm text-gray-500 dark:text-gray-400">余额</p>
            <p className="text-2xl font-bold">
              ${subsiteInfo.balance.toFixed(2)}
            </p>
            {isOwnerOrAdmin() && (
              <Link
                href={`/admin/subsite/${subsiteId}/finance`}
                className="text-xs text-blue-600 dark:text-blue-400 hover:underline"
              >
                查看财务
              </Link>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SubsiteDetailPage;
