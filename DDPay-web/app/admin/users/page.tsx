"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-hot-toast";
import { Button } from "@heroui/button";
import { Badge } from "@heroui/badge";
import { useTheme } from "next-themes";

import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { afetch } from "@/lib/afetch";
import { useAuthStore } from "@/store/auth";

interface User {
  id: number;
  email: string;
  username: string;
  avatar: string;
  role: string;
  level: number;
  status: number;
  email_verified: boolean;
  created_at: string;
  last_login_at: string;
  last_login_ip: string;
}

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [selectedLevel, setSelectedLevel] = useState<string>("1");
  const [levelDialogOpen, setLevelDialogOpen] = useState(false);
  const [statusDialogOpen, setStatusDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const router = useRouter();
  const { user, isAuthenticated, isTokenExpired } = useAuthStore();
  const { theme } = useTheme();

  // 检测设备是否为移动端
  useEffect(() => {
    const checkIfMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    // 初始检查
    checkIfMobile();

    // 添加窗口大小变化监听
    window.addEventListener("resize", checkIfMobile);

    // 清理监听器
    return () => window.removeEventListener("resize", checkIfMobile);
  }, []);

  useEffect(() => {
    // 检查用户是否已登录且是管理员
    if (!isAuthenticated || isTokenExpired() || user?.role !== "admin") {
      toast.error("请先登录管理员账号");
      router.replace("/admin", { scroll: false });

      return;
    }

    fetchUsers();
  }, [isAuthenticated, user, router]);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const res = await afetch<{
        code: number;
        msg: string;
        data: {
          users: User[];
          total: number;
          page: number;
          page_size: number;
          total_pages: number;
        };
      }>("/api/v1/user/admin/list?page=1&page_size=100", {
        headers: {
          Authorization: `Bearer ${useAuthStore.getState().accessToken}`,
        },
      });

      if (res.code === 401) {
        toast.error("登录已过期，请重新登录");
        router.replace("/admin", { scroll: false });

        return;
      }

      if (res.code !== 200) {
        toast.error(res.msg || "获取用户列表失败");

        return;
      }

      setUsers(res.data.users);
    } catch (error) {
      toast.error("获取用户列表失败");
    } finally {
      setLoading(false);
    }
  };

  const handleLevelChange = async () => {
    if (!selectedUser) return;

    try {
      setIsSubmitting(true);

      // 打印请求信息
      const requestData = {
        user_id: selectedUser.id,
        level: parseInt(selectedLevel),
      };

      const res = await afetch<{
        code: number;
        msg: string;
      }>("/api/v1/user/admin/update-level", {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${useAuthStore.getState().accessToken}`,
        },
        body: JSON.stringify(requestData),
      });

      // 打印响应信息
      console.log("收到更新用户等级响应:", res);

      if (res.code === 401) {
        toast.error("登录已过期，请重新登录");
        router.replace("/admin", { scroll: false });

        return;
      }

      if (res.code !== 200) {
        toast.error(res.msg || "更新用户等级失败");

        return;
      }

      toast.success("用户等级更新成功");
      setLevelDialogOpen(false);
      fetchUsers(); // 刷新用户列表
    } catch (error) {
      toast.error("更新用户等级失败");
      console.error("更新用户等级出错:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleStatusChange = async () => {
    if (!selectedUser) return;

    const newStatus = selectedUser.status === 1 ? 0 : 1;

    try {
      setIsSubmitting(true);
      const res = await afetch<{
        code: number;
        msg: string;
      }>("/api/v1/user/admin/update-status", {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${useAuthStore.getState().accessToken}`,
        },
        body: JSON.stringify({
          user_id: selectedUser.id,
          status: newStatus,
        }),
      });

      if (res.code === 401) {
        toast.error("登录已过期，请重新登录");
        router.replace("/login", { scroll: false });

        return;
      }

      if (res.code === 403) {
        toast.error("不允许禁用管理员账号");
        setStatusDialogOpen(false);

        return;
      }

      if (res.code !== 200) {
        toast.error(res.msg || "更新用户状态失败");

        return;
      }

      toast.success(`用户状态已${newStatus === 1 ? "启用" : "禁用"}`);
      setStatusDialogOpen(false);
      fetchUsers(); // 刷新用户列表
    } catch (error) {
      toast.error("更新用户状态失败");
      console.error("更新用户状态出错:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const openLevelDialog = (user: User) => {
    setSelectedUser(user);
    setSelectedLevel(user.level.toString());
    setLevelDialogOpen(true);
  };

  const openStatusDialog = (user: User) => {
    setSelectedUser(user);
    setStatusDialogOpen(true);
  };

  const getLevelName = (level: number) => {
    switch (level) {
      case 1:
        return "青铜会员";
      case 2:
        return "白银会员";
      case 3:
        return "黄金会员";
      case 4:
        return "钻石会员";
      default:
        return "未知";
    }
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return "从未登录";

    try {
      const date = new Date(dateString);
      // 检查日期是否有效
      if (isNaN(date.getTime())) {
        return "日期无效";
      }
      return date.toLocaleString();
    } catch (error) {
      return "日期错误";
    }
  };

  // 移动端卡片视图
  const renderMobileCards = () => (
    <div className="grid grid-cols-1 gap-4">
      {users.map((user) => (
        <div
          key={user.id}
          className="p-4 border rounded-lg shadow-sm dark:bg-gray-800 dark:border-gray-700"
        >
          <div className="flex justify-between items-center mb-2">
            <h3 className="font-semibold text-lg dark:text-white">
              {user.username}
            </h3>
            <Badge
              variant="solid"
              color={user.role === "admin" ? "danger" : "default"}
            >
              {user.role === "admin" ? "管理员" : "用户"}
            </Badge>
          </div>

          <div className="grid grid-cols-2 gap-2 mb-3">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">邮箱</p>
              <p className="dark:text-gray-200">{user.email}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                会员等级
              </p>
              <Badge
                variant="solid"
                color={
                  user.level === 4
                    ? "success"
                    : user.level === 3
                      ? "warning"
                      : user.level === 2
                        ? "primary"
                        : "default"
                }
              >
                {getLevelName(user.level)}
              </Badge>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">状态</p>
              <Badge
                variant="solid"
                color={user.status === 1 ? "success" : "danger"}
              >
                {user.status === 1 ? "正常" : "禁用"}
              </Badge>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                邮箱验证
              </p>
              <Badge
                variant="flat"
                color={user.email_verified ? "success" : "warning"}
              >
                {user.email_verified ? "已验证" : "未验证"}
              </Badge>
            </div>
          </div>

          <div className="mb-2">
            <p className="text-sm text-gray-500 dark:text-gray-400">注册时间</p>
            <p className="text-sm dark:text-gray-200">
              {formatDate(user.created_at)}
            </p>
          </div>

          <div className="mb-2">
            <p className="text-sm text-gray-500 dark:text-gray-400">最后登录</p>
            <p className="text-sm dark:text-gray-200">
              {formatDate(user.last_login_at)}
              {user.last_login_ip && (
                <span className="text-gray-500 ml-1">
                  IP: {user.last_login_ip}
                </span>
              )}
            </p>
          </div>

          <div className="flex justify-end gap-2 mt-2">
            <Button
              variant="bordered"
              size="sm"
              onClick={() => openLevelDialog(user)}
            >
              修改等级
            </Button>
            <Button
              variant="bordered"
              size="sm"
              color={user.status === 1 ? "danger" : "success"}
              onClick={() => openStatusDialog(user)}
            >
              {user.status === 1 ? "禁用" : "启用"}
            </Button>
          </div>
        </div>
      ))}
    </div>
  );

  // 桌面端表格视图
  const renderDesktopTable = () => (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead className="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              ID
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              用户名
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              邮箱
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              角色
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              会员等级
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              状态
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              邮箱验证
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              注册时间
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              最后登录
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              操作
            </th>
          </tr>
        </thead>
        <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
          {users.map((user) => (
            <tr
              key={user.id}
              className="hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {user.id}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {user.username}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {user.email}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <Badge
                  variant={user.role === "admin" ? "solid" : "flat"}
                  color={user.role === "admin" ? "danger" : "default"}
                >
                  {user.role === "admin" ? "管理员" : "用户"}
                </Badge>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <Badge
                  variant="solid"
                  color={
                    user.level === 4
                      ? "success"
                      : user.level === 3
                        ? "warning"
                        : user.level === 2
                          ? "primary"
                          : "default"
                  }
                >
                  {getLevelName(user.level)}
                </Badge>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <Badge
                  variant="solid"
                  color={user.status === 1 ? "success" : "danger"}
                >
                  {user.status === 1 ? "正常" : "禁用"}
                </Badge>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <Badge
                  variant="flat"
                  color={user.email_verified ? "success" : "warning"}
                >
                  {user.email_verified ? "已验证" : "未验证"}
                </Badge>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {formatDate(user.created_at)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <div className="text-xs">
                  <div>{formatDate(user.last_login_at)}</div>
                  {user.last_login_ip && (
                    <div className="text-gray-500 dark:text-gray-400">
                      IP: {user.last_login_ip}
                    </div>
                  )}
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <div className="flex gap-2">
                  <Button
                    variant="bordered"
                    size="sm"
                    onClick={() => openLevelDialog(user)}
                  >
                    修改等级
                  </Button>
                  <Button
                    variant="bordered"
                    size="sm"
                    color={user.status === 1 ? "danger" : "success"}
                    onClick={() => openStatusDialog(user)}
                  >
                    {user.status === 1 ? "禁用" : "启用"}
                  </Button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );

  return (
    <div className="container mx-auto py-6 px-4 md:px-6">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <h1 className="text-2xl font-bold dark:text-white">用户管理</h1>
      </div>

      {loading ? (
        <div className="text-center py-10 dark:text-gray-300">加载中...</div>
      ) : (
        <>
          {/* 根据屏幕大小显示不同视图 */}
          <div className="md:block hidden">{renderDesktopTable()}</div>
          <div className="md:hidden block">{renderMobileCards()}</div>
        </>
      )}

      {/* 修改等级对话框 */}
      <Dialog open={levelDialogOpen} onOpenChange={setLevelDialogOpen}>
        <DialogContent className="dark:bg-gray-800 dark:text-white">
          <DialogHeader>
            <DialogTitle>修改用户等级</DialogTitle>
            <DialogDescription className="dark:text-gray-400">
              为用户{" "}
              <span className="font-semibold">{selectedUser?.username}</span>{" "}
              设置新的等级
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="level" className="text-right md:block hidden">
                等级
              </Label>
              <Label htmlFor="level" className="md:hidden block">
                等级
              </Label>
              <div className="col-span-3">
                <select
                  id="level"
                  value={selectedLevel}
                  onChange={(e) => setSelectedLevel(e.target.value)}
                  className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                >
                  <option value="1">青铜会员 (最多1个分站)</option>
                  <option value="2">白银会员 (最多3个分站)</option>
                  <option value="3">黄金会员 (最多10个分站)</option>
                  <option value="4">钻石会员 (无限分站)</option>
                </select>
              </div>
            </div>

            <div className="col-span-4 text-sm text-gray-500 dark:text-gray-400 mt-2">
              <p>等级说明：</p>
              <ul className="list-disc pl-5 mt-1">
                <li>青铜会员：最多可创建1个分站</li>
                <li>白银会员：最多可创建3个分站，9折优惠</li>
                <li>黄金会员：最多可创建10个分站，8折优惠</li>
                <li>钻石会员：可创建无限数量分站，7折优惠</li>
              </ul>
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="bordered"
              onClick={() => setLevelDialogOpen(false)}
              disabled={isSubmitting}
            >
              取消
            </Button>
            <Button
              onClick={handleLevelChange}
              disabled={isSubmitting}
              color="primary"
            >
              {isSubmitting ? "提交中..." : "确认修改"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 修改状态对话框 */}
      <Dialog open={statusDialogOpen} onOpenChange={setStatusDialogOpen}>
        <DialogContent className="dark:bg-gray-800 dark:text-white">
          <DialogHeader>
            <DialogTitle>修改用户状态</DialogTitle>
            <DialogDescription className="dark:text-gray-400">
              确定要{selectedUser?.status === 1 ? "禁用" : "启用"}用户{" "}
              <span className="font-semibold">{selectedUser?.username}</span>{" "}
              吗？
              {selectedUser?.status === 1 && (
                <p className="text-red-500 mt-2">
                  禁用后该用户将无法登录系统，所有API请求将被拒绝。
                </p>
              )}
            </DialogDescription>
          </DialogHeader>

          <DialogFooter>
            <Button
              variant="bordered"
              onClick={() => setStatusDialogOpen(false)}
              disabled={isSubmitting}
            >
              取消
            </Button>
            <Button
              color={selectedUser?.status === 1 ? "danger" : "success"}
              onClick={handleStatusChange}
              disabled={isSubmitting}
            >
              {isSubmitting
                ? "提交中..."
                : `确认${selectedUser?.status === 1 ? "禁用" : "启用"}`}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
