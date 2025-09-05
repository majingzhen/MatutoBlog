import request from '@/utils/request'

// 获取附件列表（分页）
export function getAttachList(params) {
    return request({
        url: '/attach/page',
        method: 'get',
        params
    })
}

// 上传附件
export function uploadAttach(file) {
    const formData = new FormData()
    formData.append('file', file)
    
    return request({
        url: '/attach/upload',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })
}

// 删除附件
export function deleteAttach(id) {
    return request({
        url: `/attach/${id}`,
        method: 'delete'
    })
}