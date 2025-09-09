import request from '@/utils/request'

// 获取评论列表（分页）
export function getCommentList(params) {
    return request({
        url: '/comments/page',
        method: 'get',
        params
    })
}

// 更新评论状态
export function updateCommentStatus(id, status) {
    return request({
        url: `/comments/${id}/status`,
        method: 'put',
        data: { status }
    })
}

// 删除评论
export function deleteComment(id) {
    return request({
        url: `/comments/${id}`,
        method: 'delete'
    })
}

// 批量审核评论
export function batchReviewComments(ids, status) {
    return request({
        url: '/comments/batch-review',
        method: 'post',
        data: { ids, status }
    })
}

// 获取评论详情
export function getCommentById(id) {
    return request({
        url: `/comments/${id}`,
        method: 'get'
    })
}