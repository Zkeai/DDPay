import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { jwtDecode } from 'jwt-decode';

interface User {
    id: number;
    username: string;
    email: string;
    avatar: string;
    role: string;
}

interface AuthState {
    user: User | null;
    accessToken: string | null;
    refreshToken: string | null;
    expiresIn: number | null;
    isAuthenticated: boolean;

    // 登录方法
    login: (user: User, accessToken: string, refreshToken: string, expiresIn: number) => void;

    // 注册方法
    register: (user: User, accessToken: string, refreshToken: string, expiresIn: number) => void;

    // 登出方法
    logout: () => void;

    // 更新token
    updateTokens: (accessToken: string, refreshToken: string, expiresIn: number) => void;

    // 检查token是否过期
    isTokenExpired: () => boolean;

    // 获取token过期时间
    getTokenExpiration: () => number | null;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set, get) => ({
            user: null,
            accessToken: null,
            refreshToken: null,
            expiresIn: null,
            isAuthenticated: false,

            login: (user, accessToken, refreshToken, expiresIn) => {
                set({
                    user,
                    accessToken,
                    refreshToken,
                    expiresIn,
                    isAuthenticated: true,
                });
            },

            register: (user, accessToken, refreshToken, expiresIn) => {
                set({
                    user,
                    accessToken,
                    refreshToken,
                    expiresIn,
                    isAuthenticated: true,
                });
            },

            logout: () => {
                set({
                    user: null,
                    accessToken: null,
                    refreshToken: null,
                    expiresIn: null,
                    isAuthenticated: false,
                });
            },

            updateTokens: (accessToken, refreshToken, expiresIn) => {
                set({
                    accessToken,
                    refreshToken,
                    expiresIn,
                });
            },

            isTokenExpired: () => {
                const { accessToken } = get();

                if (!accessToken) return true;

                try {
                    const decoded: any = jwtDecode(accessToken);
                    const currentTime = Date.now() / 1000;

                    return decoded.exp < currentTime;
                } catch (error) {
                    return true;
                }
            },

            getTokenExpiration: () => {
                const { accessToken } = get();

                if (!accessToken) return null;

                try {
                    const decoded: any = jwtDecode(accessToken);

                    return decoded.exp;
                } catch (error) {
                    return null;
                }
            },
        }),
        {
            name: 'auth-storage', // 存储在localStorage中的键名
            partialize: (state) => ({
                user: state.user,
                accessToken: state.accessToken,
                refreshToken: state.refreshToken,
                expiresIn: state.expiresIn,
                isAuthenticated: state.isAuthenticated,
            }),
        }
    )
); 