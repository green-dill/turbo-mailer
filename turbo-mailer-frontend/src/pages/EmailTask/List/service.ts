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

export async function updateEmailTask(body: FormData, id: number, options?: { [key: string]: any }) {
  return request<Record<string, any>>(`/api/v1/task/${id}`, {
    method: 'POST',
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

export async function addEmailTask(body: FormData, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/v1/task', {
    method: 'POST',
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
