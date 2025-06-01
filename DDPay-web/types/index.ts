import { SVGProps, ReactNode } from "react";

export type IconSvgProps = SVGProps<SVGSVGElement> & {
  size?: number;
};

export interface SidebarSubItem {
  title: string;
  href: string;
  icon?: ReactNode;
}

export interface SidebarItem {
  title: string;
  icon?: ReactNode;
  children?: SidebarSubItem[];
  href?: string;
}
