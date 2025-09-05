import request from '@/utils/request'

// 获取文章列表
export function getArticleList(params) {
    return request({
        url: '/articles/page',
        method: 'get',
        params
    })
}

// 删除文章
export function deleteArticle(id) {
    return request({
        url: `/articles/${id}`,
        method: 'delete'
    })
}