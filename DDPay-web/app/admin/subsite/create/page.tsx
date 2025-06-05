"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeftIcon } from "@heroicons/react/24/outline";
import Link from "next/link";
import { toast } from "react-hot-toast";

import { afetch } from "@/lib/afetch";

const CreateSubsitePage = () => {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    subdomain: "",
    description: "",
    logo: "",
    domain: "",
    commission_rate: 10, // 默认佣金比例为10%
    theme: "default",
    status: 1, // 默认启用
  });

  // 处理表单变更
  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    let newValue: string | number = value;

    // 佣金比例字段转换为数字
    if (name === "commission_rate") {
      newValue = Number(value);
      if (isNaN(newValue) || newValue < 0) newValue = 0;
      if (newValue > 100) newValue = 100;
    }

    setFormData((prev) => ({
      ...prev,
      [name]: newValue,
    }));
  };

  // 处理表单提交
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: {
          subsite: {
            id: number;
          };
        };
      }>("/api/v1/subsite/create", {
        method: "POST",
        body: JSON.stringify(formData),
      });

      if (response.code === 200) {
        toast.success("分站创建成功");
        // 跳转到分站详情页
        router.push(`/admin/subsite/${response.data.subsite.id}`);
      } else {
        toast.error(response.msg || "创建失败");
      }
    } catch (error: any) {
      toast.error(error.message || "创建失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  };

  // 检查子域名格式
  const validateSubdomain = (subdomain: string) => {
    return /^[a-z0-9-]+$/.test(subdomain);
  };

  const isValidForm = () => {
    return (
      formData.name.trim().length > 0 &&
      formData.subdomain.trim().length > 0 &&
      validateSubdomain(formData.subdomain)
    );
  };

  return (
    <div className="container mx-auto p-4 max-w-3xl">
      <div className="mb-6 flex items-center">
        <Link
          href="/admin/subsite"
          className="flex items-center text-blue-600 hover:text-blue-800"
        >
          <ArrowLeftIcon className="h-4 w-4 mr-1" />
          返回分站列表
        </Link>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h1 className="text-2xl font-bold mb-6 text-gray-800 dark:text-white">
          创建分站
        </h1>

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="name"
            >
              分站名称 <span className="text-red-600">*</span>
            </label>
            <input
              type="text"
              id="name"
              name="name"
              className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              placeholder="请输入分站名称"
              value={formData.name}
              onChange={handleChange}
              required
            />
          </div>

          <div className="mb-4">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="subdomain"
            >
              子域名前缀 <span className="text-red-600">*</span>
            </label>
            <div className="flex items-center">
              <input
                type="text"
                id="subdomain"
                name="subdomain"
                className="flex-grow rounded-l-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                placeholder="请输入子域名前缀"
                value={formData.subdomain}
                onChange={handleChange}
                required
              />
              <span className="bg-gray-100 dark:bg-gray-600 px-3 py-2 rounded-r-lg border border-l-0 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-300">
                .ddpay.com
              </span>
            </div>
            <p className="text-xs text-gray-500 mt-1">
              只能使用小写字母、数字和连字符(-)
            </p>
            {formData.subdomain && !validateSubdomain(formData.subdomain) && (
              <p className="text-xs text-red-500 mt-1">
                子域名格式不正确，只能使用小写字母、数字和连字符
              </p>
            )}
          </div>

          <div className="mb-4">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="description"
            >
              分站描述
            </label>
            <textarea
              id="description"
              name="description"
              className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              placeholder="请输入分站描述"
              value={formData.description}
              onChange={handleChange}
              rows={3}
            />
          </div>

          <div className="mb-4">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="logo"
            >
              分站Logo
            </label>
            <input
              type="text"
              id="logo"
              name="logo"
              className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              placeholder="请输入Logo图片URL"
              value={formData.logo}
              onChange={handleChange}
            />
            <p className="text-xs text-gray-500 mt-1">填入Logo图片URL地址</p>
          </div>

          <div className="mb-4">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="domain"
            >
              自定义域名
            </label>
            <input
              type="text"
              id="domain"
              name="domain"
              className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              placeholder="请输入自定义域名（可选）"
              value={formData.domain}
              onChange={handleChange}
            />
            <p className="text-xs text-gray-500 mt-1">
              如需使用自己的域名，请在此填写，如: shop.example.com
            </p>
          </div>

          <div className="mb-4">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="commission_rate"
            >
              佣金比例 (%) <span className="text-red-600">*</span>
            </label>
            <input
              type="number"
              id="commission_rate"
              name="commission_rate"
              className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              min="0"
              max="100"
              step="0.1"
              value={formData.commission_rate}
              onChange={handleChange}
              required
            />
            <p className="text-xs text-gray-500 mt-1">
              销售额中将有该比例作为佣金返还给分站所有者
            </p>
          </div>

          <div className="mb-6">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="theme"
            >
              主题风格
            </label>
            <select
              id="theme"
              name="theme"
              className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              value={formData.theme}
              onChange={handleChange}
            >
              <option value="default">默认主题</option>
              <option value="dark">暗黑主题</option>
              <option value="light">明亮主题</option>
              <option value="anime">二次元主题</option>
            </select>
          </div>

          <div className="flex justify-end mt-6">
            <Link href="/admin/subsite">
              <button
                type="button"
                className="px-4 py-2 rounded-lg border border-gray-300 text-gray-700 mr-2 hover:bg-gray-100 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
              >
                取消
              </button>
            </Link>
            <button
              type="submit"
              className="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={loading || !isValidForm()}
            >
              {loading ? "创建中..." : "创建分站"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreateSubsitePage;
