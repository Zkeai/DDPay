"use client";

import React from "react";
import { Button } from "@heroui/button";
import { ArrowDownTrayIcon } from "@heroicons/react/24/outline";

import PageContainer from "@/components/layout/PageContainer";

const OrdersPage = () => {
  return (
    <PageContainer
      title="订单列表"
      subtitle="查看和管理所有订单"
      actions={
        <Button variant="solid" size="sm">
          <ArrowDownTrayIcon className="w-4 h-4 mr-1" />
          导出数据
        </Button>
      }
    >
      <div className="p-6">
        <p className="text-gray-700 dark:text-gray-300">
          订单列表页面内容将在这里显示。您可以查看、处理和管理所有订单。
        </p>
      </div>
    </PageContainer>
  );
};

export default OrdersPage;
