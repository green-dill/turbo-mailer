import { request } from 'umi';

export async function queryEmailTask(params: API.PageParams, options?: { [key: string]: any }) {
  return request<EmailTask.EmailTaskList>('/api/v1/task', {
    method: 'GET',
    params,
    ...(options || {}),
  }).then(res => {
    return {
      total: res.total,
      data: res.list,
      success: true,
    }
  });
}

export async function queryEmailTaskById(params: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  const { id } = params;
  return request<EmailTask.EmailTaskListItem>(`/api/v1/task/${id}`, {
    method: 'GET',
    ...(options || {}),
  }).then(res => {
    return {
      data: res,
    }
  });
}

export async function updateEmailTask(body: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  // const { id } = body;
  //
  // // 检查 pool_ids 是否存在且为数组
  // if (Array.isArray(body.pool_ids)) {
  //   // 使用 map 方法创建一个新的数组，其长度与 pool_ids 相同，所有元素都为 1
  //   body.pools_weights = body.pool_ids.map(() => 1);
  // }


  return request<Record<string, any>>(`/api/v1/task/${body.id}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: body,
    ...(options || {}),
  });
}

export async function updateTest(body: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  const { id } = body;
  return request<Record<string, any>>(`/api/v1/task/${id}/test`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function addEmailTask(body: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/v1/task', {
    method: 'POST',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: body,
    ...(options || {}),
  });
}

export async function startImmediately(params: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  const { id } = params;
  return request<Record<string, any>>(`/api/v1/task/${id}/start-immediately`, {
    method: 'POST',
    ...(options || {}),
  });
}

export async function removeEmailTask(params: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  const { id } = params;
  return request<Record<string, any>>(`/api/v1/task/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
