import { request } from 'umi';

export async function queryEmailTask(params: API.PageParams, options?: { [key: string]: any }) {
  return request<EmailTask.EmailTaskList>('/email-task', {
    method: 'GET',
    params: params,
    ...(options || {}),
  });
}

export async function updateEmailTask(body: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/email-task', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function addEmailTask(body: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/email-task', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function removeEmailTask(params: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  const { id } = params;
  return request<Record<string, any>>(`/email-task/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
