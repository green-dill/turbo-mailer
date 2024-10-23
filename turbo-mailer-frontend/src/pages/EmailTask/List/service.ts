import { request } from 'umi';
import {errorConfig} from "@/requestErrorConfig";

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
  const { ID } = params;
  return request<EmailTask.EmailTaskListItem>(`/api/v1/task/${ID}`, {
    method: 'GET',
    ...(options || {}),
  }).then(res => {
    return {
      data: res,
    }
  });
}

export async function updateEmailTask(body: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/v1/task', {
    method: 'PUT',
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
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function removeEmailTask(params: EmailTask.EmailTaskListItem, options?: { [key: string]: any }) {
  const { ID } = params;
  return request<Record<string, any>>(`/api/v1/task/${ID}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}
