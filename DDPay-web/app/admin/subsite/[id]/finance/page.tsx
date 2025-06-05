"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import {
  ArrowLeftIcon,
  ArrowPathIcon,
  ArrowDownTrayIcon,
} from "@heroicons/react/24/outline";
import { toast } from "react-hot-toast";

import { afetch } from "@/lib/afetch";
import { useAuthStore } from "@/store/auth";

// 分站财务信息类型定义
interface SubsiteFinance {
  balance: number;
  total_commission: number;
  total_sales: number;
  total_withdrawals: number;
  pending_withdrawals: number;
}

// 佣金记录类型定义
interface CommissionRecord {
  id: number;
  order_id: string;
  amount: number;
  status: number;
  remark: string;
  created_at: string;
}

// 提现记录类型定义
interface WithdrawalRecord {
  id: number;
  amount: number;
  status: number;
  payment_method: string;
  payment_account: string;
  remark: string;
  created_at: string;
  updated_at: string;
}

// 分站信息类型定义
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

const SubsiteFinancePage = () => {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuthStore();
  const [subsiteInfo, setSubsiteInfo] = useState<SubsiteInfo | null>(null);
  const [financeInfo, setFinanceInfo] = useState<SubsiteFinance | null>(null);
  const [commissions, setCommissions] = useState<CommissionRecord[]>([]);
  const [withdrawals, setWithdrawals] = useState<WithdrawalRecord[]>([]);
  const [activeTab, setActiveTab] = useState<"commissions" | "withdrawals">(
    "commissions"
  );
  const [loading, setLoading] = useState(true);
  const [withdrawalLoading, setWithdrawalLoading] = useState(false);
  const [error, setError] = useState("");
  const [withdrawalForm, setWithdrawalForm] = useState({
    amount: 0,
    payment_method: "alipay",
    payment_account: "",
    remark: "",
  });
  const [showWithdrawalForm, setShowWithdrawalForm] = useState(false);
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

        // 获取分站财务信息
        const financeResponse = await afetch<{
          code: number;
          msg: string;
          data: {
            finance: SubsiteFinance;
          };
        }>(`/api/v1/subsite/finance?id=${subsiteId}`);

        if (financeResponse.code === 200) {
          setFinanceInfo(financeResponse.data.finance);
        }

        // 获取佣金记录
        await fetchCommissions();

        // 获取提现记录
        await fetchWithdrawals();
      } else {
        setError(infoResponse.msg || "获取分站信息失败");
      }
    } catch (error: any) {
      setError(error.message || "获取分站信息失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  };

  const fetchCommissions = async () => {
    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: {
          commissions: CommissionRecord[];
        };
      }>(`/api/v1/subsite/commissions?id=${subsiteId}`);

      if (response.code === 200) {
        setCommissions(response.data.commissions);
      }
    } catch (error) {}
  };

  const fetchWithdrawals = async () => {
    try {
      const response = await afetch<{
        code: number;
        msg: string;
        data: {
          withdrawals: WithdrawalRecord[];
        };
      }>(`/api/v1/subsite/withdrawals?id=${subsiteId}`);

      if (response.code === 200) {
        setWithdrawals(response.data.withdrawals);
      }
    } catch (error) {}
  };

  const handleWithdrawalChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement
    >
  ) => {
    const { name, value } = e.target;

    // 对金额进行特殊处理
    if (name === "amount") {
      let numValue = parseFloat(value);

      if (isNaN(numValue)) numValue = 0;

      // 不能超过当前余额
      if (financeInfo && numValue > financeInfo.balance) {
        numValue = financeInfo.balance;
      }

      // 不能小于0
      if (numValue < 0) numValue = 0;

      setWithdrawalForm((prev) => ({
        ...prev,
        [name]: numValue,
      }));
    } else {
      setWithdrawalForm((prev) => ({
        ...prev,
        [name]: value,
      }));
    }
  };

  const submitWithdrawal = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!financeInfo || withdrawalForm.amount <= 0) return;

    setWithdrawalLoading(true);

    try {
      const response = await afetch<{
        code: number;
        msg: string;
      }>("/api/v1/subsite/withdraw", {
        method: "POST",
        body: JSON.stringify({
          subsite_id: subsiteId,
          ...withdrawalForm,
        }),
      });

      if (response.code === 200) {
        toast.success("提现申请已提交");
        setShowWithdrawalForm(false);
        // 重置表单
        setWithdrawalForm({
          amount: 0,
          payment_method: "alipay",
          payment_account: "",
          remark: "",
        });
        // 刷新财务信息和提现记录
        await fetchSubsiteInfo();
      } else {
        toast.error(response.msg || "提现申请失败");
      }
    } catch (error: any) {
      toast.error(error.message || "提现申请失败，请稍后重试");
    } finally {
      setWithdrawalLoading(false);
    }
  };

  // 格式化状态显示
  const formatStatus = (status: number) => {
    switch (status) {
      case 0:
        return {
          label: "待处理",
          color:
            "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
        };
      case 1:
        return {
          label: "已完成",
          color:
            "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
        };
      case 2:
        return {
          label: "已拒绝",
          color: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
        };
      default:
        return {
          label: "未知状态",
          color:
            "bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400",
        };
    }
  };

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

  // 格式化支付方式
  const formatPaymentMethod = (method: string) => {
    switch (method) {
      case "alipay":
        return "支付宝";
      case "wechat":
        return "微信支付";
      case "bank":
        return "银行卡";
      case "usdt":
        return "USDT";
      default:
        return method;
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto p-4 flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500" />
      </div>
    );
  }

  if (error || !subsiteInfo || !financeInfo) {
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
          {subsiteInfo.subsite.name} - 财务管理
        </h1>

        {/* 财务概览卡片 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
          <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
            <div className="text-sm text-blue-600 dark:text-blue-400">
              当前余额
            </div>
            <div className="text-2xl font-bold text-blue-800 dark:text-blue-300">
              ${financeInfo.balance.toFixed(2)}
            </div>
          </div>
          <div className="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
            <div className="text-sm text-green-600 dark:text-green-400">
              累计佣金
            </div>
            <div className="text-2xl font-bold text-green-800 dark:text-green-300">
              ${financeInfo.total_commission.toFixed(2)}
            </div>
          </div>
          <div className="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-4">
            <div className="text-sm text-purple-600 dark:text-purple-400">
              销售总额
            </div>
            <div className="text-2xl font-bold text-purple-800 dark:text-purple-300">
              ${financeInfo.total_sales.toFixed(2)}
            </div>
          </div>
          <div className="bg-amber-50 dark:bg-amber-900/20 rounded-lg p-4">
            <div className="text-sm text-amber-600 dark:text-amber-400">
              已提现金额
            </div>
            <div className="text-2xl font-bold text-amber-800 dark:text-amber-300">
              ${financeInfo.total_withdrawals.toFixed(2)}
            </div>
          </div>
        </div>

        {/* 提现按钮 */}
        <div className="flex justify-between items-center mb-6">
          <div className="flex space-x-4">
            <button
              className={`px-4 py-2 font-medium rounded-lg ${
                activeTab === "commissions"
                  ? "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
                  : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              }`}
              onClick={() => setActiveTab("commissions")}
            >
              佣金记录
            </button>
            <button
              className={`px-4 py-2 font-medium rounded-lg ${
                activeTab === "withdrawals"
                  ? "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
                  : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              }`}
              onClick={() => setActiveTab("withdrawals")}
            >
              提现记录
            </button>
          </div>

          <button
            className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            onClick={() => setShowWithdrawalForm(true)}
            disabled={financeInfo.balance <= 0}
          >
            <ArrowDownTrayIcon className="h-4 w-4 mr-1" />
            申请提现
          </button>
        </div>

        {/* 提现表单 */}
        {showWithdrawalForm && (
          <div className="mb-6 bg-blue-50 dark:bg-blue-900/10 p-4 rounded-lg">
            <h3 className="text-lg font-semibold mb-3">申请提现</h3>
            <form onSubmit={submitWithdrawal}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label
                    htmlFor="amount"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    提现金额 <span className="text-red-600">*</span>
                  </label>
                  <input
                    type="number"
                    id="amount"
                    name="amount"
                    className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    placeholder="请输入提现金额"
                    min="0"
                    step="0.01"
                    max={financeInfo.balance}
                    value={withdrawalForm.amount}
                    onChange={handleWithdrawalChange}
                    required
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    当前可提现余额: ${financeInfo.balance.toFixed(2)}
                  </p>
                </div>
                <div>
                  <label
                    htmlFor="payment_method"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    提现方式 <span className="text-red-600">*</span>
                  </label>
                  <select
                    id="payment_method"
                    name="payment_method"
                    className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    value={withdrawalForm.payment_method}
                    onChange={handleWithdrawalChange}
                    required
                  >
                    <option value="alipay">支付宝</option>
                    <option value="wechat">微信支付</option>
                    <option value="bank">银行卡</option>
                    <option value="usdt">USDT</option>
                  </select>
                </div>
                <div className="md:col-span-2">
                  <label
                    htmlFor="payment_account"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    收款账号 <span className="text-red-600">*</span>
                  </label>
                  <input
                    type="text"
                    id="payment_account"
                    name="payment_account"
                    className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    placeholder="请输入收款账号"
                    value={withdrawalForm.payment_account}
                    onChange={handleWithdrawalChange}
                    required
                  />
                </div>
                <div className="md:col-span-2">
                  <label
                    htmlFor="remark"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                  >
                    备注
                  </label>
                  <textarea
                    id="remark"
                    name="remark"
                    className="w-full rounded-lg border border-gray-300 p-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    placeholder="可选备注信息"
                    rows={2}
                    value={withdrawalForm.remark}
                    onChange={handleWithdrawalChange}
                  />
                </div>
              </div>
              <div className="flex justify-end space-x-2">
                <button
                  type="button"
                  className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 dark:text-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700"
                  onClick={() => setShowWithdrawalForm(false)}
                >
                  取消
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  disabled={
                    withdrawalLoading ||
                    withdrawalForm.amount <= 0 ||
                    !withdrawalForm.payment_account
                  }
                >
                  {withdrawalLoading ? "提交中..." : "提交申请"}
                </button>
              </div>
            </form>
          </div>
        )}

        {/* 佣金记录 */}
        {activeTab === "commissions" && (
          <div>
            <div className="flex justify-between items-center mb-3">
              <h3 className="text-lg font-semibold">佣金记录</h3>
              <button
                onClick={fetchCommissions}
                className="flex items-center text-blue-600 hover:text-blue-800"
                title="刷新"
              >
                <ArrowPathIcon className="h-4 w-4" />
              </button>
            </div>

            {commissions.length === 0 ? (
              <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                暂无佣金记录
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
                        订单号
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                      >
                        金额
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
                        备注
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                      >
                        时间
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
                    {commissions.map((commission) => (
                      <tr
                        key={commission.id}
                        className="hover:bg-gray-50 dark:hover:bg-gray-800"
                      >
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                          {commission.order_id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-green-600 dark:text-green-400">
                          +${commission.amount.toFixed(2)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span
                            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${formatStatus(commission.status).color}`}
                          >
                            {formatStatus(commission.status).label}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                          {commission.remark || "-"}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                          {formatDate(commission.created_at)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* 提现记录 */}
        {activeTab === "withdrawals" && (
          <div>
            <div className="flex justify-between items-center mb-3">
              <h3 className="text-lg font-semibold">提现记录</h3>
              <button
                onClick={fetchWithdrawals}
                className="flex items-center text-blue-600 hover:text-blue-800"
                title="刷新"
              >
                <ArrowPathIcon className="h-4 w-4" />
              </button>
            </div>

            {withdrawals.length === 0 ? (
              <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                暂无提现记录
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
                        金额
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
                        提现方式
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                      >
                        账号
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                      >
                        备注
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                      >
                        申请时间
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
                    {withdrawals.map((withdrawal) => (
                      <tr
                        key={withdrawal.id}
                        className="hover:bg-gray-50 dark:hover:bg-gray-800"
                      >
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-amber-600 dark:text-amber-400">
                          ${withdrawal.amount.toFixed(2)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span
                            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${formatStatus(withdrawal.status).color}`}
                          >
                            {formatStatus(withdrawal.status).label}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                          {formatPaymentMethod(withdrawal.payment_method)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                          {withdrawal.payment_account}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                          {withdrawal.remark || "-"}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                          {formatDate(withdrawal.created_at)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default SubsiteFinancePage;
