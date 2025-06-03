"use client";
import React, { useState, useEffect, useRef } from "react";
import { usePathname, useRouter } from "next/navigation";
import {
  ShoppingBagIcon,
  ShoppingCartIcon,
  UserGroupIcon,
  Cog6ToothIcon,
  CreditCardIcon,
  ChartBarIcon,
  TagIcon,
  ClipboardDocumentListIcon,
  TruckIcon,
  DocumentTextIcon,
  HomeIcon,
  ArrowPathIcon,
  BuildingStorefrontIcon,
  UserIcon,
  WrenchScrewdriverIcon,
  ComputerDesktopIcon,
  KeyIcon,
  XMarkIcon,
  Bars3Icon,
  ChevronDownIcon,
  ChevronRightIcon,
  ArrowRightOnRectangleIcon,
  ClockIcon,
} from "@heroicons/react/24/outline";
import Link from "next/link";
import { Image } from "@heroui/image";

import { SidebarItem } from "@/types";
import Sidebar from "@/components/layout/sidebar";
import TopNav from "@/components/layout/topnav";
import { logout, validateToken } from "@/lib/afetch";
import { useAuthStore } from "@/store/auth";

interface AdminLayoutProps {
  children: React.ReactNode;
}

const AdminLayout = ({ children }: AdminLayoutProps) => {
  const [collapsed, setCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const userMenuRef = useRef<HTMLDivElement>(null);
  const pathname = usePathname();
  const router = useRouter();
  const [currentRoute, setCurrentRoute] = useState("首页");
  const [isMobile, setIsMobile] = useState(false);
  const [openMenus, setOpenMenus] = useState<{ [key: string]: boolean }>({});
  const { user } = useAuthStore();
  const [isValidToken, setIsValidToken] = useState(true);

  // 验证JWT令牌有效性
  useEffect(() => {
    const checkTokenValidity = () => {
      const isValid = validateToken();
      setIsValidToken(isValid);

      // 如果令牌无效且当前不在登录页面，则重定向到登录页
      if (!isValid && !pathname.includes("/admin/login")) {
        router.push("/admin");
      }
    };

    // 初次加载检查
    checkTokenValidity();

    // 设置定时器，定期检查令牌有效性
    const tokenCheckInterval = setInterval(checkTokenValidity, 60000); // 每分钟检查一次

    return () => {
      clearInterval(tokenCheckInterval);
    };
  }, [pathname, router]);

  // 所有菜单项
  const sidebarItems: SidebarItem[] = [
    {
      title: "控制台",
      icon: <HomeIcon className="w-5 h-5" />,
      href: "/admin/dashboard",
    },
    {
      title: "店铺管理",
      icon: <BuildingStorefrontIcon className="w-5 h-5" />,
      children: [
        {
          title: "商品分类",
          href: "/admin/store/categories",
          icon: <TagIcon className="w-4 h-4" />,
        },
        {
          title: "商品管理",
          href: "/admin/products",
          icon: <ShoppingBagIcon className="w-4 h-4" />,
        },
      ],
    },
    {
      title: "订单管理",
      icon: <ShoppingCartIcon className="w-5 h-5" />,
      children: [
        {
          title: "订单列表",
          href: "/admin/orders",
          icon: <ClipboardDocumentListIcon className="w-4 h-4" />,
        },
        {
          title: "售后处理",
          href: "/admin/orders/after-sale",
          icon: <ArrowPathIcon className="w-4 h-4" />,
        },
        {
          title: "发货管理",
          href: "/admin/orders/shipping",
          icon: <TruckIcon className="w-4 h-4" />,
        },
      ],
    },
    {
      title: "用户管理",
      icon: <UserGroupIcon className="w-5 h-5" />,
      children: [
        {
          title: "用户列表",
          href: "/admin/users",
          icon: <UserIcon className="w-4 h-4" />,
        },
        {
          title: "会员等级",
          href: "/admin/users/levels",
          icon: <TagIcon className="w-4 h-4" />,
        },
      ],
    },
    {
      title: "站点管理",
      icon: <ComputerDesktopIcon className="w-5 h-5" />,
      children: [
        {
          title: "站点设置",
          href: "/admin/site/settings",
          icon: <Cog6ToothIcon className="w-4 h-4" />,
        },
        {
          title: "页面管理",
          href: "/admin/site/pages",
          icon: <DocumentTextIcon className="w-4 h-4" />,
        },
      ],
    },
    {
      title: "支付接口",
      icon: <CreditCardIcon className="w-5 h-5" />,
      children: [
        {
          title: "接口配置",
          href: "/admin/payment/config",
          icon: <WrenchScrewdriverIcon className="w-4 h-4" />,
        },
        {
          title: "接口状态",
          href: "/admin/payment/status",
          icon: <ChartBarIcon className="w-4 h-4" />,
        },
      ],
    },
    {
      title: "支付订单",
      icon: <CreditCardIcon className="w-5 h-5" />,
      children: [
        {
          title: "交易记录",
          href: "/admin/finance/transactions",
          icon: <ClipboardDocumentListIcon className="w-4 h-4" />,
        },
        {
          title: "退款管理",
          href: "/admin/finance/refunds",
          icon: <ArrowPathIcon className="w-4 h-4" />,
        },
      ],
    },
    {
      title: "系统设置",
      icon: <Cog6ToothIcon className="w-5 h-5" />,
      children: [
        {
          title: "基础设置",
          href: "/admin/settings/basic",
          icon: <Cog6ToothIcon className="w-4 h-4" />,
        },
        {
          title: "安全设置",
          href: "/admin/settings/security",
          icon: <KeyIcon className="w-4 h-4" />,
        },
      ],
    },
  ];

  // 菜单分类
  const categories = [
    {
      name: "MAIN",
      items: [sidebarItems[0]], // 控制台
    },
    {
      name: "SHOP",
      items: [sidebarItems[1], sidebarItems[2]], // 店铺管理、订单管理
    },
    {
      name: "USER",
      items: [sidebarItems[3], sidebarItems[4]], // 用户管理、站点管理
    },
    {
      name: "PAY",
      items: [sidebarItems[5], sidebarItems[6]], // 支付接口、支付订单
    },
    {
      name: "CONFIG",
      items: [sidebarItems[7]], // 系统设置
    },
  ];

  // 检测设备尺寸
  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 768);
      // 在移动设备上默认折叠侧边栏
      if (window.innerWidth < 768) {
        setCollapsed(true);
      } else if (window.innerWidth >= 1024) {
        // 在大屏上默认展开侧边栏
        setCollapsed(false);
      }
    };

    // 初始检测
    handleResize();

    window.addEventListener("resize", handleResize);

    return () => window.removeEventListener("resize", handleResize);
  }, []);

  // 初始化时根据当前路径展开相应菜单
  useEffect(() => {
    const initialOpenMenus: { [key: string]: boolean } = {};

    sidebarItems.forEach((item) => {
      if (item.children) {
        // 如果当前路径匹配子菜单项，则展开父菜单
        const isActive = item.children.some(
          (subItem) =>
            pathname === subItem.href || pathname.startsWith(subItem.href + "/")
        );

        if (isActive) {
          initialOpenMenus[item.title] = true;
        }
      }
    });

    setOpenMenus(initialOpenMenus);
  }, [pathname]);

  // 点击外部关闭用户菜单
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        userMenuRef.current &&
        !userMenuRef.current.contains(event.target as Node)
      ) {
        setUserMenuOpen(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  useEffect(() => {
    // 根据路径设置当前路由名称
    const findMatchingRoute = () => {
      // 先检查直接链接的菜单项
      for (const item of sidebarItems) {
        if (item.href === pathname) {
          setCurrentRoute(item.title);

          return;
        }

        // 检查子菜单项
        if (item.children) {
          const subItem = item.children.find((sub) => sub.href === pathname);

          if (subItem) {
            setCurrentRoute(subItem.title);

            return;
          }
        }
      }

      // 如果没有精确匹配，查找路径前缀匹配
      for (const item of sidebarItems) {
        if (item.children) {
          const subItem = item.children.find((sub) =>
            pathname.startsWith(sub.href)
          );

          if (subItem) {
            setCurrentRoute(subItem.title);

            return;
          }
        }
      }

      // 默认值
      setCurrentRoute("首页");
    };

    findMatchingRoute();

    // 在移动设备上，导航后自动关闭菜单
    if (isMobile) {
      setMobileMenuOpen(false);
    }
  }, [pathname, isMobile]);

  const handleLogout = async () => {
    try {
      await logout();
      router.push("/admin");
    } catch (error) {
      console.error("注销失败:", error);
      // 即使失败也强制清除本地令牌并跳转
      localStorage.removeItem("auth-storage");
      router.push("/admin");
    }
  };

  const handleToggleCollapse = (isCollapsed: boolean) => {
    setCollapsed(isCollapsed);
  };

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  const toggleUserMenu = () => {
    setUserMenuOpen(!userMenuOpen);
  };

  const toggleMenu = (title: string) => {
    setOpenMenus((prevOpenMenus) => ({
      ...prevOpenMenus,
      [title]: !prevOpenMenus[title],
    }));
  };

  const isActive = (href: string) => {
    return pathname === href || pathname.startsWith(href + "/");
  };

  // 如果令牌无效，不渲染主布局
  if (!isValidToken) {
    return <div>{children}</div>;
  }

  return (
    <div className="flex flex-col md:flex-row h-screen overflow-hidden">
      {/* 移动端导航栏 */}
      <div className="md:hidden bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 px-4 py-3 flex items-center justify-between">
        <div className="flex items-center">
          <button
            aria-controls="mobile-menu"
            aria-expanded={mobileMenuOpen}
            aria-label={mobileMenuOpen ? "关闭菜单" : "打开菜单"}
            className="text-gray-600 dark:text-gray-300 focus:outline-none"
            onClick={toggleMobileMenu}
          >
            {mobileMenuOpen ? (
              <XMarkIcon className="h-6 w-6" />
            ) : (
              <Bars3Icon className="h-6 w-6" />
            )}
          </button>
          <span className="ml-4 font-bold text-blue-600 dark:text-blue-400">
            DDPay
          </span>
        </div>

        {/* 移动端右侧用户菜单 */}
        <div className="relative" ref={userMenuRef}>
          <button
            onClick={toggleUserMenu}
            className="flex items-center"
            aria-expanded={userMenuOpen}
            aria-haspopup="true"
          >
            <div className="w-8 h-8 rounded-full overflow-hidden">
              <Image
                alt="用户头像"
                className="w-full h-full object-cover"
                src={user?.avatar || "http://42.51.0.159:4399/favicon.ico"}
              />
            </div>
          </button>

          {/* 移动端用户下拉菜单 */}
          {userMenuOpen && (
            <div className="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 z-50">
              <div className="py-1">
                <div className="px-4 py-2 border-b border-gray-200 dark:border-gray-700">
                  <p className="text-sm font-medium text-gray-900 dark:text-white">
                    {user?.username || "管理员"}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {currentRoute}
                  </p>
                </div>
                <Link
                  href="/admin/profile/logs"
                  className="flex items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                  onClick={() => setUserMenuOpen(false)}
                >
                  <ClockIcon className="w-5 h-5 mr-3 text-gray-500 dark:text-gray-400" />
                  <span>登录日志</span>
                </Link>
                <Link
                  href="/admin/profile/security"
                  className="flex items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                  onClick={() => setUserMenuOpen(false)}
                >
                  <KeyIcon className="w-5 h-5 mr-3 text-gray-500 dark:text-gray-400" />
                  <span>安全设置</span>
                </Link>
                <div className="border-t border-gray-200 dark:border-gray-700" />
                <button
                  className="flex items-center w-full px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700"
                  onClick={handleLogout}
                >
                  <ArrowRightOnRectangleIcon className="w-5 h-5 mr-3" />
                  <span>注销登录</span>
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* 移动端侧边栏 - 覆盖模式 */}
      {mobileMenuOpen && (
        <div className="fixed inset-0 z-40 md:hidden">
          {/* 背景遮罩 */}
          <div
            aria-hidden="true"
            className="fixed inset-0 bg-gray-600 bg-opacity-75 transition-opacity"
            onClick={toggleMobileMenu}
          />

          {/* 侧边栏内容 */}
          <div className="fixed inset-y-0 left-0 flex max-w-full">
            <div
              className="relative flex w-full max-w-xs flex-col overflow-y-auto bg-white dark:bg-gray-900 pb-4 shadow-xl"
              id="mobile-menu"
            >
              {/* 头部 */}
              <div className="sticky top-0 z-10 bg-white dark:bg-gray-900">
                <div className="flex items-center justify-between px-4 py-4">
                  <div className="font-bold text-xl text-blue-600 dark:text-blue-400">
                    DDPay
                  </div>
                  <button
                    aria-label="关闭菜单"
                    className="inline-flex items-center justify-center rounded-md text-gray-400 hover:text-gray-500 focus:outline-none"
                    type="button"
                    onClick={toggleMobileMenu}
                  >
                    <XMarkIcon aria-hidden="true" className="h-6 w-6" />
                  </button>
                </div>

                {/* 用户头像 */}
                <div className="flex items-center px-4 py-3 border-b border-gray-200 dark:border-gray-700">
                  <div className="w-10 h-10 rounded-full overflow-hidden mr-3">
                    <Image
                      alt="用户头像"
                      className="w-full h-full object-cover"
                      src={
                        user?.avatar || "http://42.51.0.159:4399/favicon.ico"
                      }
                    />
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-medium text-gray-900 dark:text-white">
                      {user?.username || "管理员"}
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {user?.role || "管理员账户"}
                    </p>
                  </div>
                </div>
              </div>

              {/* 移动端菜单项 */}
              <nav className="flex-1 overflow-y-auto px-2 py-4">
                {categories.map((category) => (
                  <div key={category.name} className="mb-6">
                    <div className="mb-2 px-3">
                      <span className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
                        {category.name}
                      </span>
                    </div>
                    <div className="space-y-1">
                      {category.items.map((item) => (
                        <div key={item.title} className="rounded-md">
                          {item.href ? (
                            <Link
                              className={`flex items-center px-3 py-2 text-sm font-medium rounded-md ${
                                isActive(item.href)
                                  ? "bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
                                  : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800"
                              }`}
                              href={item.href}
                              onClick={toggleMobileMenu}
                            >
                              <span className="mr-3 flex-shrink-0">
                                {item.icon}
                              </span>
                              <span>{item.title}</span>
                            </Link>
                          ) : (
                            <div className="space-y-1">
                              <button
                                className={`flex items-center justify-between w-full px-3 py-2 text-sm font-medium rounded-md ${
                                  item.children?.some((sub) =>
                                    isActive(sub.href)
                                  )
                                    ? "bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
                                    : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800"
                                }`}
                                onClick={() => toggleMenu(item.title)}
                                aria-expanded={openMenus[item.title]}
                              >
                                <span className="flex items-center">
                                  <span className="mr-3 flex-shrink-0">
                                    {item.icon}
                                  </span>
                                  <span>{item.title}</span>
                                </span>
                                {openMenus[item.title] ? (
                                  <ChevronDownIcon className="w-4 h-4 flex-shrink-0" />
                                ) : (
                                  <ChevronRightIcon className="w-4 h-4 flex-shrink-0" />
                                )}
                              </button>

                              {openMenus[item.title] && item.children && (
                                <div className="pl-10 pr-2 py-1 space-y-1 bg-gray-50 dark:bg-gray-800/50 rounded-md mt-1 mb-1">
                                  {item.children.map((subItem) => (
                                    <Link
                                      key={subItem.href}
                                      className={`flex items-center py-2 px-3 text-sm rounded-md ${
                                        isActive(subItem.href)
                                          ? "bg-blue-50 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
                                          : "text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800/80"
                                      }`}
                                      href={subItem.href}
                                      onClick={toggleMobileMenu}
                                    >
                                      {subItem.icon && (
                                        <span className="mr-2 flex-shrink-0">
                                          {subItem.icon}
                                        </span>
                                      )}
                                      <span>{subItem.title}</span>
                                    </Link>
                                  ))}
                                </div>
                              )}
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </nav>

              {/* 登出按钮 */}
              <div className="border-t border-gray-200 dark:border-gray-700 p-4 mt-auto">
                <button
                  className="flex items-center w-full px-3 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-red-500 dark:hover:text-red-400 rounded-md hover:bg-gray-100 dark:hover:bg-gray-800"
                  onClick={handleLogout}
                >
                  <ArrowRightOnRectangleIcon className="w-5 h-5 mr-3 flex-shrink-0" />
                  退出登录
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* 桌面端侧边栏 */}
      <div className={`hidden md:block ${collapsed ? "md:w-0" : "md:w-64"}`}>
        <Sidebar
          avatarUrl={user?.avatar || "http://42.51.0.159:4399/favicon.ico"}
          categories={categories}
          collapsed={collapsed}
          items={sidebarItems}
          userName={user?.username || "管理员"}
          onLogout={handleLogout}
          openMenus={openMenus}
          onToggleMenu={toggleMenu}
        />
      </div>

      {/* 主内容区域 */}
      <div
        className={`flex flex-1 flex-col overflow-hidden ${collapsed ? "md:ml-0" : ""}`}
      >
        {/* 桌面端顶部导航 */}
        <div className="hidden md:block">
          <TopNav
            collapsed={collapsed}
            currentRoute={currentRoute}
            onToggleCollapse={handleToggleCollapse}
            userName={user?.username || "管理员"}
            avatarUrl={user?.avatar || "http://42.51.0.159:4399/favicon.ico"}
            onLogout={handleLogout}
          />
        </div>

        {/* 主内容 */}
        <main className="flex-1 overflow-auto">
          <div className="mx-auto container">{children}</div>
        </main>
      </div>
    </div>
  );
};

export default AdminLayout;
