import request from '@/utils/request'

/**
 * 用户登录
 * @param {string} account - 用户名
 * @param {string} password - 密码
 * @returns {Promise}
 */
export function login(account, password) {
    return request({
        url: '/login',
        method: 'post',
        data: {
            account,
            password
        }
    })
}

/**
 * 用户登出
 * @returns {Promise}
 */
export function logout() {
    return request({
        url: '/logout',
        method: 'post'
    })
}

/**
 * 获取当前用户信息
 * @returns {Promise}
 */
export function getUserInfo() {
    return request({
        url: '/profile',
        method: 'get'
    })
}
