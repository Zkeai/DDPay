"use client";

import React from "react";
import { Button } from "@heroui/button";
import { PlusIcon } from "@heroicons/react/24/outline";

import PageContainer from "@/components/layout/PageContainer";

const CategoriesPage = () => {
  return (
    <PageContainer
      title="商品分类"
      subtitle="管理店铺中的商品分类"
      actions={
        <Button variant="solid" size="sm">
          <PlusIcon className="w-4 h-4 mr-1" />
          添加分类
        </Button>
      }
    >
      <div className="p-6">
        <p className="text-gray-700 dark:text-gray-300">
          商品分类页面内容将在这里显示。您可以添加、编辑和删除商品分类。
        </p>
      </div>
    </PageContainer>
  );
};

export default CategoriesPage;
