"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeftIcon, TrashIcon } from "@heroicons/react/24/outline";
import Link from "next/link";
import { toast } from "react-hot-toast";
import { Button } from "@heroui/button";

import { afetch } from "@/lib/afetch";
import { useAuthStore } from "@/store/auth";

const EditSubsitePage = () => {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuthStore();
  const subsiteId = params.id;

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [formData, setFormData] = useState({
    id: Number(subsiteId),
    name: "",
    subdomain: "",
    description: "",
    logo: "",
    domain: "",
    commission_rate: 10,
    theme: "default",
    status: 1,
  });

  // 加载分站数据
  useEffect(() => {
    fetchSubsiteData();
  }, [subsiteId]);

  const fetchSubsiteData = async () => {
    setLoading(true);
    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: {
          subsite_info: {
            subsite: {
              id: number;
              name: string;
              subdomain: string;
              domain: string;
              logo: string;
              description: string;
              theme: string;
              status: number;
              commission_rate: number;
            };
          };
        };
      }>(`/api/v1/subsite/info?id=${subsiteId}`);

      if (response.code === 200 && response.data?.subsite_info?.subsite) {
        const subsite = response.data.subsite_info.subsite;

        setFormData({
          id: subsite.id,
          name: subsite.name || "",
          subdomain: subsite.subdomain || "",
          description: subsite.description || "",
          logo: subsite.logo || "",
          domain: subsite.domain || "",
          commission_rate: subsite.commission_rate || 10,
          theme: subsite.theme || "default",
          status: subsite.status,
        });
      } else {
        toast.error(response.msg || "获取分站信息失败");
        router.push(`/admin/subsite/${subsiteId}`);
      }
    } catch (error: any) {
      toast.error(error.message || "获取分站信息失败，请稍后重试");
      router.push(`/admin/subsite/${subsiteId}`);
    } finally {
      setLoading(false);
    }
  };

  // 处理表单变更
  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    let newValue: string | number | boolean = value;

    // 佣金比例字段转换为数字
    if (name === "commission_rate") {
      newValue = Number(value);
      if (isNaN(newValue as number) || newValue < 0) newValue = 0;
      if (newValue > 100) newValue = 100;
    }

    // 状态字段转换为数字
    if (name === "status") {
      newValue = Number(value);
    }

    setFormData((prev) => ({
      ...prev,
      [name]: newValue,
    }));
  };

  // 处理表单提交
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: any;
      }>("/api/v1/subsite/update", {
        method: "PUT",
        body: JSON.stringify(formData),
      });

      if (response.code === 200) {
        toast.success("分站更新成功");
        // 跳转到分站详情页
        router.push(`/admin/subsite/${subsiteId}`);
      } else {
        toast.error(response.msg || "更新失败");
      }
    } catch (error: any) {
      toast.error(error.message || "更新失败，请稍后重试");
    } finally {
      setSubmitting(false);
    }
  };

  // 处理删除分站
  const handleDelete = async () => {
    // 二次确认
    if (!window.confirm("确定要删除该分站吗？此操作不可恢复！")) {
      return;
    }

    setDeleting(true);
    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: any;
      }>(`/api/v1/subsite/delete?id=${subsiteId}`, {
        method: "DELETE",
      });

      if (response.code === 200) {
        toast.success("分站已成功删除");
        // 跳转到分站列表页
        router.push("/admin/subsite");
      } else {
        toast.error(response.msg || "删除失败");
      }
    } catch (error: any) {
      toast.error(error.message || "删除失败，请稍后重试");
    } finally {
      setDeleting(false);
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

  // 检查是否为管理员或分站所有者
  const isOwnerOrAdmin = () => {
    if (!user) return false;

    return user.role === "admin"; // 分站所有者检查在后端已实现
  };

  if (loading) {
    return (
      <div className="container mx-auto p-4 flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500" />
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4 max-w-3xl">
      <div className="mb-6 flex items-center justify-between">
        <Link
          href={`/admin/subsite/${subsiteId}`}
          className="flex items-center text-blue-600 hover:text-blue-800"
        >
          <ArrowLeftIcon className="h-4 w-4 mr-1" />
          返回分站详情
        </Link>

        {isOwnerOrAdmin() && (
          <Button
            size="sm"
            variant="solid"
            className="flex items-center bg-red-500 hover:bg-red-600 text-white"
            onClick={handleDelete}
            disabled={deleting}
          >
            <TrashIcon className="h-4 w-4 mr-1" />
            {deleting ? "删除中..." : "删除分站"}
          </Button>
        )}
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h1 className="text-2xl font-bold mb-6 text-gray-800 dark:text-white">
          编辑分站
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

          <div className="mb-4">
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
            <p className="text-xs text-gray-500 mt-1">选择分站的显示风格</p>
          </div>

          <div className="mb-6">
            <label
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              htmlFor="status"
            >
              分站状态
            </label>
            <select
              id="status"
              name="status"
              className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              value={formData.status}
              onChange={handleChange}
            >
              <option value={1}>启用</option>
              <option value={0}>禁用</option>
            </select>
            <p className="text-xs text-gray-500 mt-1">
              禁用状态下分站将无法访问
            </p>
          </div>

          <div className="flex justify-between">
            <Link href={`/admin/subsite/${subsiteId}`}>
              <Button size="md" variant="bordered">
                取消
              </Button>
            </Link>
            <Button
              type="submit"
              size="md"
              variant="solid"
              disabled={!isValidForm() || submitting}
            >
              {submitting ? "保存中..." : "保存修改"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditSubsitePage;
