"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-hot-toast";
import { Button } from "@heroui/button";
import { Badge } from "@heroui/badge";
import { useTheme } from "next-themes";

import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
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

interface MembershipLevel {
  id: number;
  name: string;
  level: number;
  icon: string;
  price: number;
  description: string;
  discount_rate: number;
  max_subsites: number;
  custom_service_access: boolean;
  vip_group_access: boolean;
  priority: number;
  created_at: string;
  updated_at: string;
  benefits: MembershipBenefit[];
}

interface MembershipBenefit {
  id: number;
  level_id: number;
  title: string;
  description: string;
  icon: string;
}

export default function MembershipLevelsPage() {
  const [levels, setLevels] = useState<MembershipLevel[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedLevel, setSelectedLevel] = useState<MembershipLevel | null>(
    null
  );
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const router = useRouter();
  const { user, isAuthenticated, isTokenExpired } = useAuthStore();
  const { theme } = useTheme();

  // 编辑表单状态
  const [formData, setFormData] = useState({
    id: 0,
    name: "",
    level: 1,
    icon: "",
    price: 0,
    description: "",
    discount_rate: 1.0,
    max_subsites: 1,
    custom_service_access: false,
    vip_group_access: false,
    priority: 1,
  });

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

    fetchMembershipLevels();
  }, [isAuthenticated, user, router]);

  const fetchMembershipLevels = async () => {
    setLoading(true);
    try {
      const res = await afetch<{
        code: number;
        msg: string;
        data: {
          levels: MembershipLevel[];
        };
      }>("/api/v1/membership/levels", {});

      if (res.code !== 200) {
        toast.error(res.msg || "获取会员等级列表失败");

        return;
      }

      setLevels(res.data.levels);
    } catch (error) {
      toast.error("获取会员等级列表失败");
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value, type } = e.target as HTMLInputElement;

    if (type === "checkbox") {
      const target = e.target as HTMLInputElement;

      setFormData({
        ...formData,
        [name]: target.checked,
      });
    } else if (type === "number") {
      setFormData({
        ...formData,
        [name]: parseFloat(value),
      });
    } else {
      setFormData({
        ...formData,
        [name]: value,
      });
    }
  };

  const openEditDialog = (level: MembershipLevel) => {
    setSelectedLevel(level);
    setFormData({
      id: level.id,
      name: level.name,
      level: level.level,
      icon: level.icon,
      price: level.price,
      description: level.description,
      discount_rate: level.discount_rate,
      max_subsites: level.max_subsites,
      custom_service_access: level.custom_service_access,
      vip_group_access: level.vip_group_access,
      priority: level.priority,
    });
    setEditDialogOpen(true);
  };

  const openDeleteDialog = (level: MembershipLevel) => {
    setSelectedLevel(level);
    setDeleteDialogOpen(true);
  };

  const openCreateDialog = () => {
    setFormData({
      id: 0,
      name: "",
      level:
        levels.length > 0 ? Math.max(...levels.map((l) => l.level)) + 1 : 1,
      icon: "",
      price: 0,
      description: "",
      discount_rate: 1.0,
      max_subsites: 1,
      custom_service_access: false,
      vip_group_access: false,
      priority:
        levels.length > 0 ? Math.max(...levels.map((l) => l.priority)) + 1 : 1,
    });
    setCreateDialogOpen(true);
  };

  const handleCreateLevel = async () => {
    try {
      setIsSubmitting(true);
      const res = await afetch<{
        code: number;
        msg: string;
      }>("/api/v1/membership/admin/level", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${useAuthStore.getState().accessToken}`,
        },
        body: JSON.stringify(formData),
      });

      if (res.code === 401) {
        toast.error("登录已过期，请重新登录");
        router.replace("/admin", { scroll: false });

        return;
      }

      if (res.code !== 200) {
        toast.error(res.msg || "创建会员等级失败");

        return;
      }

      toast.success("会员等级创建成功");
      setCreateDialogOpen(false);
      fetchMembershipLevels(); // 刷新列表
    } catch (error) {
      toast.error("创建会员等级失败");
      console.error("创建会员等级出错:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleUpdateLevel = async () => {
    if (!selectedLevel) return;

    try {
      setIsSubmitting(true);
      const res = await afetch<{
        code: number;
        msg: string;
      }>("/api/v1/membership/admin/level", {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${useAuthStore.getState().accessToken}`,
        },
        body: JSON.stringify(formData),
      });

      if (res.code === 401) {
        toast.error("登录已过期，请重新登录");
        router.replace("/admin", { scroll: false });

        return;
      }

      if (res.code !== 200) {
        toast.error(res.msg || "更新会员等级失败");

        return;
      }

      toast.success("会员等级更新成功");
      setEditDialogOpen(false);
      fetchMembershipLevels(); // 刷新列表
    } catch (error) {
      toast.error("更新会员等级失败");
      console.error("更新会员等级出错:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteLevel = async () => {
    if (!selectedLevel) return;

    try {
      setIsSubmitting(true);
      const res = await afetch<{
        code: number;
        msg: string;
      }>(`/api/v1/membership/admin/level?id=${selectedLevel.id}`, {
        method: "DELETE",
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
        toast.error(res.msg || "删除会员等级失败");

        return;
      }

      toast.success("会员等级删除成功");
      setDeleteDialogOpen(false);
      fetchMembershipLevels(); // 刷新列表
    } catch (error) {
      toast.error("删除会员等级失败");
      console.error("删除会员等级出错:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return "未设置";

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

  const renderLevelForm = () => (
    <div className="grid gap-4 py-4">
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="name" className="text-right md:block hidden">
          等级名称
        </Label>
        <Label htmlFor="name" className="md:hidden block">
          等级名称
        </Label>
        <Input
          id="name"
          name="name"
          value={formData.name}
          onChange={handleInputChange}
          className="col-span-3 md:col-span-3"
          required
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="level" className="text-right">
          等级值
        </Label>
        <Input
          id="level"
          name="level"
          type="number"
          min="1"
          value={formData.level}
          onChange={handleInputChange}
          className="col-span-3"
          required
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="icon" className="text-right">
          图标URL
        </Label>
        <Input
          id="icon"
          name="icon"
          value={formData.icon}
          onChange={handleInputChange}
          className="col-span-3"
          placeholder="/assets/membership/bronze.png"
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="price" className="text-right">
          价格
        </Label>
        <Input
          id="price"
          name="price"
          type="number"
          min="0"
          step="0.01"
          value={formData.price}
          onChange={handleInputChange}
          className="col-span-3"
          required
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="description" className="text-right">
          描述
        </Label>
        <Input
          id="description"
          name="description"
          value={formData.description}
          onChange={handleInputChange}
          className="col-span-3"
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="discount_rate" className="text-right">
          折扣率
        </Label>
        <Input
          id="discount_rate"
          name="discount_rate"
          type="number"
          min="0"
          max="1"
          step="0.1"
          value={formData.discount_rate}
          onChange={handleInputChange}
          className="col-span-3"
          required
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="max_subsites" className="text-right">
          最大分站数
        </Label>
        <Input
          id="max_subsites"
          name="max_subsites"
          type="number"
          min="-1"
          value={formData.max_subsites}
          onChange={handleInputChange}
          className="col-span-3"
          required
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="priority" className="text-right">
          优先级
        </Label>
        <Input
          id="priority"
          name="priority"
          type="number"
          min="1"
          value={formData.priority}
          onChange={handleInputChange}
          className="col-span-3"
          required
        />
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label className="text-right">专属客服权限</Label>
        <div className="col-span-3 flex items-center space-x-2">
          <input
            type="checkbox"
            id="custom_service_access"
            name="custom_service_access"
            checked={formData.custom_service_access}
            onChange={handleInputChange}
            className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
          />
          <Label htmlFor="custom_service_access">启用</Label>
        </div>
      </div>

      <div className="grid grid-cols-4 items-center gap-4">
        <Label className="text-right">VIP群权限</Label>
        <div className="col-span-3 flex items-center space-x-2">
          <input
            type="checkbox"
            id="vip_group_access"
            name="vip_group_access"
            checked={formData.vip_group_access}
            onChange={handleInputChange}
            className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
          />
          <Label htmlFor="vip_group_access">启用</Label>
        </div>
      </div>
    </div>
  );

  // 移动端卡片视图
  const renderMobileCards = () => (
    <div className="grid grid-cols-1 gap-4">
      {levels.map((level) => (
        <div
          key={level.id}
          className="p-4 border rounded-lg shadow-sm dark:bg-gray-800 dark:border-gray-700"
        >
          <div className="flex justify-between items-center mb-2">
            <h3 className="font-semibold text-lg dark:text-white">
              {level.name}
            </h3>
            <Badge
              variant="solid"
              color={
                level.level > 3
                  ? "success"
                  : level.level > 2
                    ? "warning"
                    : level.level > 1
                      ? "primary"
                      : "default"
              }
            >
              等级 {level.level}
            </Badge>
          </div>

          <div className="grid grid-cols-2 gap-2 mb-3">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">价格</p>
              <p className="dark:text-gray-200">¥{level.price.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">折扣率</p>
              <p className="dark:text-gray-200">
                {(level.discount_rate * 10).toFixed(1)}折
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                最大分站数
              </p>
              <p className="dark:text-gray-200">
                {level.max_subsites === -1 ? "无限" : level.max_subsites}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                创建时间
              </p>
              <p className="text-sm dark:text-gray-200">
                {formatDate(level.created_at)}
              </p>
            </div>
          </div>

          <div className="mb-2">
            <p className="text-sm text-gray-500 dark:text-gray-400">特殊权限</p>
            <div className="flex flex-wrap gap-1 mt-1">
              {level.custom_service_access && (
                <Badge variant="flat" color="primary">
                  专属客服
                </Badge>
              )}
              {level.vip_group_access && (
                <Badge variant="flat" color="success">
                  VIP群
                </Badge>
              )}
            </div>
          </div>

          <div className="flex justify-end gap-2 mt-2">
            <Button
              variant="bordered"
              size="sm"
              onClick={() => openEditDialog(level)}
            >
              编辑
            </Button>
            <Button
              variant="bordered"
              size="sm"
              color="danger"
              onClick={() => openDeleteDialog(level)}
            >
              删除
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
              等级
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              名称
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              价格
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              折扣率
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              最大分站数
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              特殊权限
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              创建时间
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              操作
            </th>
          </tr>
        </thead>
        <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
          {levels.map((level) => (
            <tr
              key={level.id}
              className="hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {level.id}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <Badge
                  variant="solid"
                  color={
                    level.level > 3
                      ? "success"
                      : level.level > 2
                        ? "warning"
                        : level.level > 1
                          ? "primary"
                          : "default"
                  }
                >
                  {level.level}
                </Badge>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {level.name}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                ¥{level.price.toFixed(2)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {(level.discount_rate * 10).toFixed(1)}折
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {level.max_subsites === -1 ? "无限" : level.max_subsites}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <div className="flex flex-col gap-1">
                  {level.custom_service_access && (
                    <Badge variant="flat" color="primary">
                      专属客服
                    </Badge>
                  )}
                  {level.vip_group_access && (
                    <Badge variant="flat" color="success">
                      VIP群
                    </Badge>
                  )}
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                {formatDate(level.created_at)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-gray-200">
                <div className="flex gap-2">
                  <Button
                    variant="bordered"
                    size="sm"
                    onClick={() => openEditDialog(level)}
                  >
                    编辑
                  </Button>
                  <Button
                    variant="bordered"
                    size="sm"
                    color="danger"
                    onClick={() => openDeleteDialog(level)}
                  >
                    删除
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
        <h1 className="text-2xl font-bold dark:text-white">会员等级管理</h1>
        <div className="flex items-center gap-2">
          <Button onClick={openCreateDialog} color="primary">
            新增会员等级
          </Button>
        </div>
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

      {/* 对话框部分保持不变，但添加暗色主题支持 */}
      <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
        <DialogContent className="dark:bg-gray-800 dark:text-white">
          <DialogHeader>
            <DialogTitle>创建会员等级</DialogTitle>
            <DialogDescription className="dark:text-gray-400">
              填写会员等级信息，创建新的会员等级
            </DialogDescription>
          </DialogHeader>

          {renderLevelForm()}

          <DialogFooter>
            <Button
              variant="bordered"
              onClick={() => setCreateDialogOpen(false)}
              disabled={isSubmitting}
            >
              取消
            </Button>
            <Button
              onClick={handleCreateLevel}
              disabled={isSubmitting}
              color="primary"
            >
              {isSubmitting ? "提交中..." : "创建"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 编辑会员等级对话框 */}
      <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
        <DialogContent className="dark:bg-gray-800 dark:text-white">
          <DialogHeader>
            <DialogTitle>编辑会员等级</DialogTitle>
            <DialogDescription className="dark:text-gray-400">
              修改会员等级{" "}
              <span className="font-semibold">{selectedLevel?.name}</span>{" "}
              的信息
            </DialogDescription>
          </DialogHeader>

          {renderLevelForm()}

          <DialogFooter>
            <Button
              variant="bordered"
              onClick={() => setEditDialogOpen(false)}
              disabled={isSubmitting}
            >
              取消
            </Button>
            <Button
              onClick={handleUpdateLevel}
              disabled={isSubmitting}
              color="primary"
            >
              {isSubmitting ? "提交中..." : "保存"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 删除会员等级对话框 */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent className="dark:bg-gray-800 dark:text-white">
          <DialogHeader>
            <DialogTitle>删除会员等级</DialogTitle>
            <DialogDescription className="dark:text-gray-400">
              确定要删除会员等级{" "}
              <span className="font-semibold">{selectedLevel?.name}</span> 吗？
              <p className="text-red-500 mt-2">
                此操作不可撤销，删除后可能会影响已使用该等级的用户。
              </p>
            </DialogDescription>
          </DialogHeader>

          <DialogFooter>
            <Button
              variant="bordered"
              onClick={() => setDeleteDialogOpen(false)}
              disabled={isSubmitting}
            >
              取消
            </Button>
            <Button
              color="danger"
              onClick={handleDeleteLevel}
              disabled={isSubmitting}
            >
              {isSubmitting ? "提交中..." : "确认删除"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
