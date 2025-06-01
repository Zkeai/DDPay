import React, { ReactNode, useEffect } from "react";

import { useTitle } from "@/components/TitleContext";

interface PageContainerProps {
  children: ReactNode;
  title?: string;
  subtitle?: string;
  actions?: ReactNode;
  pageTitle?: string; // 用于设置浏览器标题
}

/**
 * 响应式页面容器组件
 * 为所有页面提供一致的响应式布局和样式
 */
const PageContainer: React.FC<PageContainerProps> = ({
  children,
  title,
  subtitle,
  actions,
  pageTitle,
}) => {
  const { setTitle } = useTitle();

  // 如果提供了pageTitle，则设置页面标题
  useEffect(() => {
    if (pageTitle) {
      setTitle(pageTitle);
    } else if (title) {
      // 如果没有提供pageTitle，则使用title
      setTitle(`DDPay - ${title}`);
    }
  }, [pageTitle, title, setTitle]);

  return (
    <div className="p-4 md:p-6">
      {(title || actions) && (
        <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-6 gap-4">
          <div>
            {title && (
              <h1 className="text-xl md:text-2xl font-bold text-gray-900 dark:text-white">
                {title}
              </h1>
            )}
            {subtitle && (
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {subtitle}
              </p>
            )}
          </div>
          {actions && <div className="flex items-center gap-2">{actions}</div>}
        </div>
      )}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700">
        {children}
      </div>
    </div>
  );
};

export default PageContainer;
