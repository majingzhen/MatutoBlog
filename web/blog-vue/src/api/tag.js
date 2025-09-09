import request from '@/utils/request'

// 获取标签列表（分页）
export function getTagList(params) {
    return request({
        url: '/tags/page',
        method: 'get',
        params
    })
}

// 获取所有标签（不分页）
export function getAllTags() {
    return request({
        url: '/tags',
        method: 'get'
    })
}

// 创建标签
export function createTag(data) {
    return request({
        url: '/tags/create',
        method: 'post',
        data
    })
}

// 更新标签
export function updateTag(id, data) {
    return request({
        url: `/tags/${id}`,
        method: 'put',
        data
    })
}

// 删除标签
export function deleteTag(id) {
    return request({
        url: `/tags/${id}`,
        method: 'delete'
    })
}

// 获取标签详情
export function getTagById(id) {
    return request({
        url: `/tags/${id}`,
        method: 'get'
    })
}