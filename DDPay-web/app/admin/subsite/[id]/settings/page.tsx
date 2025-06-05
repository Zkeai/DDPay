"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeftIcon } from "@heroicons/react/24/outline";
import { toast } from "react-hot-toast";
import { Button } from "@heroui/button";

import { afetch } from "@/lib/afetch";
import { useAuthStore } from "@/store/auth";

interface SubsiteSettings {
  [key: string]: string;
}

interface SubsiteInfo {
  subsite: {
    id: number;
    name: string;
    owner_id: number;
  };
  owner: {
    id: number;
  };
}

const SubsiteSettingsPage = () => {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuthStore();
  const [subsiteInfo, setSubsiteInfo] = useState<SubsiteInfo | null>(null);
  const [settings, setSettings] = useState<SubsiteSettings>({
    seo_title: "",
    seo_keywords: "",
    seo_description: "",
    contact_email: "",
    contact_qq: "",
    contact_telegram: "",
    contact_wechat: "",
    announcement: "",
    footer_text: "",
    custom_css: "",
    custom_js: "",
    enable_registration: "1", // 默认启用注册
    enable_guest_purchase: "1", // 默认允许游客购买
    enable_alipay: "1", // 默认启用支付宝
    enable_wechatpay: "1", // 默认启用微信支付
    enable_crypto: "0", // 默认禁用加密货币支付
    crypto_address: "", // 加密货币地址
    enable_paypal: "0", // 默认禁用PayPal
    paypal_account: "", // PayPal账号
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const subsiteId = params.id;

  useEffect(() => {
    fetchSubsiteInfo();
  }, [subsiteId]);

  const fetchSubsiteInfo = async () => {
    setLoading(true);
    setError("");

    try {
      // 获取分站基本信息
      const infoResponse = await afetch<{
        code: number;
        msg: string;
        data: {
          subsite_info: SubsiteInfo;
        };
      }>(`/api/v1/subsite/info?id=${subsiteId}`);

      if (infoResponse.code === 200) {
        setSubsiteInfo(infoResponse.data.subsite_info);

        // 检查权限
        if (
          user?.role !== "admin" &&
          user?.id !== infoResponse.data.subsite_info.owner.id
        ) {
          setError("您没有权限管理此分站");
          setLoading(false);

          return;
        }

        // 获取分站JSON配置
        const configResponse = await afetch<{
          code: number;
          msg: string;
          data: {
            config: any;
          };
        }>(`/api/v1/subsite/config?subsite_id=${subsiteId}`);

        if (configResponse.code === 200 && configResponse.data.config) {
          // 将配置转换为字符串键值对形式
          const configData = configResponse.data.config;
          const formattedSettings: SubsiteSettings = {};

          // 遍历所有配置项并转换为字符串
          Object.keys(configData).forEach((key) => {
            const value = configData[key];

            formattedSettings[key] = String(value);
          });

          setSettings((prev) => ({
            ...prev,
            ...formattedSettings,
          }));
        } else {
          // 初始化空配置
          setSettings({});
        }
      } else {
        setError(infoResponse.msg || "获取分站信息失败");
      }
    } catch (error: any) {
      setError(error.message || "获取分站信息失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;

    setSettings((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!subsiteInfo) return;

    setSaving(true);
    setError("");

    try {
      // 使用JSON配置API
      await afetch<{
        code: number;
        msg: string;
      }>("/api/v1/subsite/config", {
        method: "POST",
        body: JSON.stringify({
          subsite_id: Number(subsiteId),
          config: settings,
        }),
      });

      toast.success("设置已保存");
      // 保存成功后返回上一级页面
      router.push(`/admin/subsite/${subsiteId}`);
    } catch (error: any) {
      setError(error.message || "保存设置失败，请稍后重试");
      toast.error("保存设置失败");
    } finally {
      setSaving(false);
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

  return (
    <div className="container mx-auto p-4 max-w-4xl">
      <div className="mb-6 flex items-center">
        <Link
          href={`/admin/subsite/${subsiteId}`}
          className="flex items-center text-blue-600 hover:text-blue-800"
        >
          <ArrowLeftIcon className="h-4 w-4 mr-1" />
          返回分站详情
        </Link>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
        <h1 className="text-2xl font-bold mb-6 text-gray-800 dark:text-white">
          {subsiteInfo.subsite.name} - 分站设置
        </h1>

        <form onSubmit={handleSubmit}>
          {/* 分组：SEO设置 */}
          <div className="mb-8">
            <h2 className="text-lg font-semibold mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
              SEO设置
            </h2>
            <div className="grid grid-cols-1 gap-4">
              <div>
                <label
                  htmlFor="seo_title"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  网站标题
                </label>
                <input
                  type="text"
                  id="seo_title"
                  name="seo_title"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="网站标题，显示在浏览器标签页"
                  value={settings.seo_title || ""}
                  onChange={handleChange}
                />
              </div>
              <div>
                <label
                  htmlFor="seo_keywords"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  关键词
                </label>
                <input
                  type="text"
                  id="seo_keywords"
                  name="seo_keywords"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="网站关键词，用逗号分隔"
                  value={settings.seo_keywords || ""}
                  onChange={handleChange}
                />
              </div>
              <div>
                <label
                  htmlFor="seo_description"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  网站描述
                </label>
                <textarea
                  id="seo_description"
                  name="seo_description"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="网站描述，简要介绍您的网站"
                  rows={3}
                  value={settings.seo_description || ""}
                  onChange={handleChange}
                />
              </div>
            </div>
          </div>

          {/* 分组：联系方式 */}
          <div className="mb-8">
            <h2 className="text-lg font-semibold mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
              联系方式
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label
                  htmlFor="contact_email"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  联系邮箱
                </label>
                <input
                  type="email"
                  id="contact_email"
                  name="contact_email"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="您的联系邮箱"
                  value={settings.contact_email || ""}
                  onChange={handleChange}
                />
              </div>
              <div>
                <label
                  htmlFor="contact_qq"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  QQ
                </label>
                <input
                  type="text"
                  id="contact_qq"
                  name="contact_qq"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="您的QQ号码"
                  value={settings.contact_qq || ""}
                  onChange={handleChange}
                />
              </div>
              <div>
                <label
                  htmlFor="contact_telegram"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  Telegram
                </label>
                <input
                  type="text"
                  id="contact_telegram"
                  name="contact_telegram"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="您的Telegram用户名"
                  value={settings.contact_telegram || ""}
                  onChange={handleChange}
                />
              </div>
              <div>
                <label
                  htmlFor="contact_wechat"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  微信
                </label>
                <input
                  type="text"
                  id="contact_wechat"
                  name="contact_wechat"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="您的微信号"
                  value={settings.contact_wechat || ""}
                  onChange={handleChange}
                />
              </div>
            </div>
          </div>

          {/* 分组：支付通道设置 */}
          <div className="mb-8">
            <h2 className="text-lg font-semibold mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
              支付通道设置
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label
                  htmlFor="enable_alipay"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  支付宝
                </label>
                <select
                  id="enable_alipay"
                  name="enable_alipay"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  value={settings.enable_alipay || "1"}
                  onChange={handleChange}
                >
                  <option value="1">启用</option>
                  <option value="0">禁用</option>
                </select>
              </div>
              <div>
                <label
                  htmlFor="enable_wechatpay"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  微信支付
                </label>
                <select
                  id="enable_wechatpay"
                  name="enable_wechatpay"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  value={settings.enable_wechatpay || "1"}
                  onChange={handleChange}
                >
                  <option value="1">启用</option>
                  <option value="0">禁用</option>
                </select>
              </div>
              <div>
                <label
                  htmlFor="enable_crypto"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  加密货币支付
                </label>
                <select
                  id="enable_crypto"
                  name="enable_crypto"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  value={settings.enable_crypto || "0"}
                  onChange={handleChange}
                >
                  <option value="1">启用</option>
                  <option value="0">禁用</option>
                </select>
              </div>
              {settings.enable_crypto === "1" && (
                <div>
                  <label
                    htmlFor="crypto_address"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    加密货币收款地址
                  </label>
                  <input
                    type="text"
                    id="crypto_address"
                    name="crypto_address"
                    className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    placeholder="请输入USDT(TRC20)收款地址"
                    value={settings.crypto_address || ""}
                    onChange={handleChange}
                  />
                </div>
              )}
              <div>
                <label
                  htmlFor="enable_paypal"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  PayPal支付
                </label>
                <select
                  id="enable_paypal"
                  name="enable_paypal"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  value={settings.enable_paypal || "0"}
                  onChange={handleChange}
                >
                  <option value="1">启用</option>
                  <option value="0">禁用</option>
                </select>
              </div>
              {settings.enable_paypal === "1" && (
                <div>
                  <label
                    htmlFor="paypal_account"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    PayPal账号
                  </label>
                  <input
                    type="email"
                    id="paypal_account"
                    name="paypal_account"
                    className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    placeholder="您的PayPal账号邮箱"
                    value={settings.paypal_account || ""}
                    onChange={handleChange}
                  />
                </div>
              )}
            </div>
          </div>

          {/* 分组：网站设置 */}
          <div className="mb-8">
            <h2 className="text-lg font-semibold mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
              网站设置
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label
                  htmlFor="enable_registration"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  启用用户注册
                </label>
                <select
                  id="enable_registration"
                  name="enable_registration"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  value={settings.enable_registration || "1"}
                  onChange={handleChange}
                >
                  <option value="1">启用</option>
                  <option value="0">禁用</option>
                </select>
              </div>
              <div>
                <label
                  htmlFor="enable_guest_purchase"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  允许游客购买
                </label>
                <select
                  id="enable_guest_purchase"
                  name="enable_guest_purchase"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  value={settings.enable_guest_purchase || "1"}
                  onChange={handleChange}
                >
                  <option value="1">允许</option>
                  <option value="0">不允许（强制登录）</option>
                </select>
              </div>
            </div>
            <div>
              <label
                htmlFor="announcement"
                className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              >
                网站公告
              </label>
              <textarea
                id="announcement"
                name="announcement"
                className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                placeholder="网站公告，显示在首页顶部"
                rows={3}
                value={settings.announcement || ""}
                onChange={handleChange}
              />
            </div>
            <div className="mt-4">
              <label
                htmlFor="footer_text"
                className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              >
                页脚文本
              </label>
              <textarea
                id="footer_text"
                name="footer_text"
                className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                placeholder="页脚文本，支持HTML"
                rows={2}
                value={settings.footer_text || ""}
                onChange={handleChange}
              />
              <p className="text-xs text-gray-500 mt-1">
                支持HTML，可添加版权信息、备案号等
              </p>
            </div>
          </div>

          {/* 分组：高级设置 */}
          <div className="mb-8">
            <h2 className="text-lg font-semibold mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
              高级设置
            </h2>
            <div className="grid grid-cols-1 gap-4">
              <div>
                <label
                  htmlFor="custom_css"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  自定义CSS
                </label>
                <textarea
                  id="custom_css"
                  name="custom_css"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white font-mono"
                  placeholder="/* 自定义CSS样式 */"
                  rows={5}
                  value={settings.custom_css || ""}
                  onChange={handleChange}
                />
                <p className="text-xs text-gray-500 mt-1">
                  自定义CSS样式，将被添加到网站头部
                </p>
              </div>
              <div>
                <label
                  htmlFor="custom_js"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  自定义JavaScript
                </label>
                <textarea
                  id="custom_js"
                  name="custom_js"
                  className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white font-mono"
                  placeholder="// 自定义JavaScript代码"
                  rows={5}
                  value={settings.custom_js || ""}
                  onChange={handleChange}
                />
                <p className="text-xs text-gray-500 mt-1">
                  自定义JavaScript代码，将被添加到网站底部
                </p>
              </div>
            </div>
          </div>

          {/* 提交按钮 */}
          <div className="flex justify-end mt-6">
            <Link href={`/admin/subsite/${subsiteId}`}>
              <Button variant="bordered" className="mr-2">
                取消
              </Button>
            </Link>
            <Button type="submit" variant="solid" disabled={saving}>
              {saving ? "保存中..." : "保存设置"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default SubsiteSettingsPage;
