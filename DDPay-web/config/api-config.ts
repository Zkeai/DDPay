/**
 * API配置
 */
export const ApiConfig = {
    // API基础URL
    baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:2900',

    // API路径前缀
    apiPrefix: process.env.NEXT_PUBLIC_API_PREFIX || '/api/v1',

    // 认证相关
    auth: {
        // 登录端点
        loginEndpoint: '/user/login',

        // 注册端点
        registerEndpoint: '/user/register',

        // 刷新令牌端点
        refreshTokenEndpoint: '/user/refresh-token',

        // 登出端点
        logoutEndpoint: '/user/logout',

        // OAuth端点
        githubLoginEndpoint: '/user/oauth/github',
        googleLoginEndpoint: '/user/oauth/google',
    },

    // 超时配置（毫秒）
    timeout: 10000,

    // 自动刷新令牌缓冲时间（毫秒），令牌过期前多长时间触发刷新
    refreshTokenBuffer: 5 * 60 * 1000, // 5分钟
};
