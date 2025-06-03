import { useAuthStore } from '@/store/auth';

interface AfetchOptions extends RequestInit {
    baseUrl?: string;
    skipAuth?: boolean;
    refreshOnUnauthorized?: boolean;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:2900';

/**
 * 验证JWT令牌是否有效
 * @returns 令牌是否有效
 */
export const validateToken = (): boolean => {
    // 服务器端渲染时返回true，避免水合错误
    if (typeof window === 'undefined') {
        return true;
    }

    const authStore = useAuthStore.getState();

    // 检查是否已登录
    if (!authStore.isAuthenticated || !authStore.accessToken) {
        return false;
    }

    // 检查令牌是否过期
    return !authStore.isTokenExpired();
};

/**
 * 封装的fetch函数，自动携带JWT令牌
 * @param url 请求地址
 * @param options 请求选项
 * @returns Promise
 */
export const afetch = async <T = any>(
    url: string,
    options: AfetchOptions = {}
): Promise<T> => {
    const {
        baseUrl = API_BASE_URL,
        skipAuth = false,
        refreshOnUnauthorized = true,
        headers = {},
        ...rest
    } = options;

    // 获取认证状态
    const authStore = useAuthStore.getState();
    const { accessToken, refreshToken, updateTokens } = authStore;

    // 如果需要认证但没有登录或令牌已过期，则抛出错误
    if (!skipAuth && !validateToken()) {
        // 尝试使用刷新令牌获取新的访问令牌
        if (refreshToken && refreshOnUnauthorized) {
            try {
                const refreshData = await refreshAccessToken(refreshToken);

                if (refreshData.code === 200 && refreshData.data) {
                    const { access_token, refresh_token, expires_in } = refreshData.data;

                    updateTokens(access_token, refresh_token, expires_in);
                } else {
                    authStore.logout();
                    throw new Error('会话已过期，请重新登录');
                }
            } catch (error) {
                authStore.logout();
                throw new Error('会话已过期，请重新登录');
            }
        } else {
            authStore.logout();
            throw new Error('未登录或会话已过期，请重新登录');
        }
    }

    // 如果需要认证且有token，则添加到请求头
    const authHeaders: HeadersInit = {};

    if (!skipAuth && accessToken) {
        authHeaders['Authorization'] = `Bearer ${accessToken}`;
    }

    // 合并请求头
    const mergedHeaders = {
        'Content-Type': 'application/json',
        ...authHeaders,
        ...headers,
    };

    try {
        // 发送请求
        const response = await fetch(`${baseUrl}${url}`, {
            ...rest,
            headers: mergedHeaders,
        });

        // 如果返回401未授权，且有刷新令牌，尝试刷新令牌
        if (response.status === 401 && refreshToken && refreshOnUnauthorized) {
            try {
                // 尝试刷新令牌
                const refreshData = await refreshAccessToken(refreshToken);

                // 更新令牌
                if (refreshData.code === 200 && refreshData.data) {
                    const { access_token, refresh_token, expires_in } = refreshData.data;

                    updateTokens(access_token, refresh_token, expires_in);

                    // 使用新令牌重试原始请求
                    const retryResponse = await fetch(`${baseUrl}${url}`, {
                        ...rest,
                        headers: {
                            ...mergedHeaders,
                            'Authorization': `Bearer ${access_token}`,
                        },
                    });

                    // 处理重试响应
                    const data = await retryResponse.json();

                    return data as T;
                }

                // 如果刷新令牌失败，则登出
                authStore.logout();
                throw new Error('会话已过期，请重新登录');
            } catch (refreshError) {
                authStore.logout();
                throw new Error('会话已过期，请重新登录');
            }
        }

        // 检查响应状态
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));

            throw new Error(errorData.msg || `请求失败: ${response.status}`);
        }

        // 解析响应数据
        const data = await response.json();

        return data as T;
    } catch (error) {
        throw error;
    }
};

/**
 * 刷新访问令牌
 * @param refreshToken 刷新令牌
 * @returns 刷新结果
 */
export const refreshAccessToken = async (refreshToken: string) => {
    const response = await fetch(`${API_BASE_URL}/api/v1/user/refresh-token`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
        throw new Error('刷新令牌失败');
    }

    return await response.json();
};

/**
 * 发送登录请求
 * @param email 邮箱
 * @param password 密码
 * @returns 登录结果
 */
export const login = async (email: string, password: string) => {
    const response = await afetch<{
        code: number;
        msg: string;
        data: {
            user: any;
            access_token: string;
            refresh_token: string;
            expires_in: number;
        };
    }>('/api/v1/user/login', {
        method: 'POST',
        skipAuth: true,
        body: JSON.stringify({ email, password }),
    });

    return response;
};

/**
 * 发送注册请求
 * @param email 邮箱
 * @param password 密码
 * @param username 用户名
 * @param code 验证码
 * @returns 注册结果
 */
export const register = async (email: string, password: string, username: string, code: string) => {
    const response = await afetch<{
        code: number;
        msg: string;
        data: {
            user: any;
            access_token: string;
            refresh_token: string;
            expires_in: number;
        };
    }>('/api/v1/user/register', {
        method: 'POST',
        skipAuth: true,
        body: JSON.stringify({ email, password, username, code }),
    });

    return response;
};

/**
 * 发送验证码
 * @param email 邮箱
 * @param type 类型（register或reset_password）
 * @returns 发送结果
 */
export const sendVerificationCode = async (email: string, type: 'register' | 'reset_password') => {
    const response = await afetch<{
        code: number;
        msg: string;
    }>('/api/v1/user/send-code', {
        method: 'POST',
        skipAuth: true,
        body: JSON.stringify({ email, type }),
    });

    return response;
};

/**
 * 重置密码
 * @param email 邮箱
 * @param code 验证码
 * @param password 新密码
 * @returns 重置结果
 */
export const resetPassword = async (email: string, code: string, password: string) => {
    const response = await afetch<{
        code: number;
        msg: string;
    }>('/api/v1/user/reset-password', {
        method: 'POST',
        skipAuth: true,
        body: JSON.stringify({ email, code, new_password: password }),
    });

    return response;
};

/**
 * 检查邮箱是否已存在
 * @param email 邮箱
 * @returns 检查结果
 */
export const checkEmailExists = async (email: string) => {
    try {
        const response = await afetch<{
            code: number;
            msg: string;
            data: {
                exists: boolean;
            };
        }>(`/api/v1/user/check-email?email=${encodeURIComponent(email)}`, {
            method: 'GET',
            skipAuth: true,
        });

        return {
            exists: response.data.exists,
            message: response.msg
        };
    } catch (error: any) {
        // 如果是网络错误等情况，则抛出错误
        throw error;
    }
};

/**
 * 注销登录 - 清除本地token和服务器Redis中的token
 * @returns 注销结果
 */
export const logout = async () => {
    try {
        // 调用后端API清除Redis中的token
        const response = await afetch<{
            code: number;
            msg: string;
        }>('/api/v1/user/logout', {
            method: 'POST',
        });

        // 无论API调用是否成功，都清除本地存储的token
        useAuthStore.getState().logout();

        return response;
    } catch (error) {
        // 即使API调用失败，也要清除本地存储的token
        useAuthStore.getState().logout();
        throw error;
    }
}; 