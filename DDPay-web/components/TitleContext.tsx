"use client";

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { usePathname } from "next/navigation";

// 定义标题上下文类型
interface TitleContextType {
  title: string;
  setTitle: (title: string) => void;
}

// 创建标题上下文
const TitleContext = createContext<TitleContextType>({
  title: "DDPay",
  setTitle: () => {},
});

// 创建自定义钩子以便组件使用
export const useTitle = () => useContext(TitleContext);

interface TitleProviderProps {
  children: ReactNode;
}

// 标题映射表
const pathTitleMap: Record<string, string> = {
  "/admin/dashboard": "控制台",
  "/admin/products": "商品管理",
  "/admin/store/categories": "商品分类",
  "/admin/orders": "订单列表",
  "/admin/orders/after-sale": "售后处理",
  "/admin/orders/shipping": "发货管理",
  "/admin/users": "用户列表",
  "/admin/users/levels": "会员等级",
  "/admin/site/settings": "站点设置",
  "/admin/site/pages": "页面管理",
  "/admin/payment/config": "接口配置",
  "/admin/payment/status": "接口状态",
  "/admin/finance/transactions": "交易记录",
  "/admin/finance/refunds": "退款管理",
  "/admin/settings/basic": "基础设置",
  "/admin/settings/security": "安全设置",
};

export const TitleProvider: React.FC<TitleProviderProps> = ({ children }) => {
  const [title, setTitle] = useState("DDPay");
  const pathname = usePathname();

  // 根据路径自动设置标题
  useEffect(() => {
    const defaultTitle = "DDPay";

    if (pathname) {
      // 检查是否有对应的标题
      const pathTitle = pathTitleMap[pathname];

      if (pathTitle) {
        setTitle(`${pathTitle}`);
        // 更新浏览器标题
        document.title = `${defaultTitle} - ${pathTitle}`;
      } else {
        setTitle(defaultTitle);
        document.title = defaultTitle;
      }
    }
  }, [pathname]);

  // 允许组件手动设置标题
  const updateTitle = (newTitle: string) => {
    setTitle(newTitle);
    document.title = newTitle;
  };

  return (
    <TitleContext.Provider value={{ title, setTitle: updateTitle }}>
      {children}
    </TitleContext.Provider>
  );
};
