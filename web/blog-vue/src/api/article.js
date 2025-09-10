import request from '@/utils/request'

// 获取文章列表
export function getArticleList(params) {
    return request({
        url: '/articles/page',
        method: 'get',
        params
    })
}

// 获取文章详情
export function getArticleById(id) {
    return request({
        url: `/articles/${id}`,
        method: 'get'
    })
}

// 创建文章
export function createArticle(data) {
    return request({
        url: '/articles/publish',
        method: 'post',
        data
    })
}

// 更新文章
export function updateArticle(id, data) {
    return request({
        url: `/articles/update/${id}`,
        method: 'put',
        data
    })
}

// 删除文章
export function deleteArticle(id) {
    return request({
        url: `/articles/${id}`,
        method: 'delete'
    })
}

// 获取分类列表
export function getCategoryList() {
    return request({
        url: '/categories/enable-list',
        method: 'get'
    })
}

// 获取标签列表
export function getTagList() {
    return request({
        url: '/tags/enable-list',
        method: 'get'
    })
}