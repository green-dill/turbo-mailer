import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { addEmailTask, removeEmailTask, queryEmailTask, updateEmailTask, } from './service';
import {Button, message, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {ModalForm, ProFormDateTimePicker, ProFormRadio, ProFormText, ProFormTextArea} from "@ant-design/pro-components";
import {ProFormUploadButton} from "@ant-design/pro-form";

const EmailTaskList: React.FC = () => {
  const [createModalOpen, handleModalOpen] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();

  /**
   * 添加或修改邮件任务
   * @param fields 邮件任务
   */
  const handleSubmit = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      const response = fields.id ? await updateEmailTask(fields) : await addEmailTask(fields);
      if (response && response.code === 200) {
        message.success('编辑成功');
      } else {
        message.error('编辑失败请重试！');
      }
      return true;
    } catch (error) {
      message.error('添加失败请重试！');
      return false;
    }
  };

  /**
   *  删除邮件任务
   * @param fields 邮件任务
   */
  const handleRemove = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      const response = await removeEmailTask(fields);
      if (response && response.code === 200) {
        message.success('删除成功');
      }
      return true;
    } catch (error) {
      message.error('删除失败请重试！');
      return false;
    }
  };

  const columns: ProColumns<EmailTask.EmailTaskListItem>[] = [
    // id?: number;
    //     title?: string;
    //     recipients?: string;
    //     numberPoolId?: number[];
    //     sendInterval?: string;
    //     taskStartedAt?: Date;
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      hideInTable: true,
    },
    {
      title: '邮件标题',
      dataIndex: 'title',
      tooltip: '支持模板语法',
    },
    {
      title: '邮件内容地址',
      dataIndex: 'content',
      tooltip: 'HTML/TET 文件，支持模板语法',
      hideInSearch: true,
    },
    {
      title: '收件人数量',
      dataIndex: 'recipients',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '号池列表',
      dataIndex: 'numberPoolIds',
      hideInSearch: true,
    },
    {
      title: '发送频率',
      dataIndex: 'sendInterval',
      hideInSearch: true,
    },
    {
      title: '任务开始时间',
      dataIndex: 'taskStartedAt',
      hideInSearch: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueEnum: {
        0: {
          text: '正常',
          status: 'Success',
        },
        1: {
          text: '禁用',
          status: 'Error',
        },
      },
      filters: true,
      onFilter: true,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      valueType: 'date',
      hideInSearch: true,
      sorter: true,
      defaultSortOrder: 'descend',
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      valueType: 'date',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record: EmailTask.EmailTaskListItem) => [
        <Button
          type="link"
          size="small"
          style={{padding: 0}}
          key='edit'
          onClick={() => {
            handleModalOpen(true);
          }}
        >
          编辑
        </Button>,
        <Popconfirm
          key='delete'
          title="确定删除吗？"
          onConfirm={async () => {
            const success = await handleRemove(record);
            if (success) {
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
        >
          <Button style={{padding: 0}} type="link" size="small">删除</Button>
        </Popconfirm>,
        <Popconfirm
          key='test'
          title="输入测试邮件地址"
        >
          <Button style={{padding: 0}} type="link" size="small">测试</Button>
        </Popconfirm>,
        <Popconfirm
          key='start'
          title="确定开始发送吗？"
          onConfirm={async () => {
            const success = await handleRemove(record);
            if (success) {
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
        >
          <Button style={{padding: 0}} type="link" size="small">开始</Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<EmailTask.EmailTaskListItem, API.PageParams>
        headerTitle="邮件任务列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button type="primary" onClick={() => handleModalOpen(true)}>
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={queryEmailTask}
        columns={columns}
      />
      <ModalForm
        title='创建邮件任务'
        width="400px"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        onFinish={async (value) => {
          const success = await handleSubmit(value as NumberPool.NumberPoolListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          label='邮件标题'
          placeholder='请输入'
          tooltip='支持模板语法'
          rules={[
            {
              required: true,
              message: '号池邮件标题不能为空',
            },
          ]}
          width="md"
          name="name"
        />
        <ProFormUploadButton
          label='邮件内容'
          placeholder='请上传'
          tooltip='上传 HTML/TXT 文件，支持模板语法'
          name="content"
        />
        <ProFormUploadButton
          label='收件人列表'
          placeholder='请上传'
          tooltip='上传 CSV/EXCEL 文件'
          name="recipients"
        />
        <ProFormText
          label='发送频率'
          placeholder='请输入'
          tooltip='每小时发送的邮件数量'
          rules={[
            {
              required: true,
              message: '发送频率不能为空',
            },
          ]}
          width="md"
          name="sendInterval"
        />
        <ProFormDateTimePicker
          label='发送时间'
          placeholder='请选择'
          tooltip='邮件开始发送的时间'
          rules={[
            {
              required: true,
              message: '发送时间不能为空',
            },
          ]}
          width="md"
          name="taskStartedAt"
        />
        <ProFormRadio.Group
          name="status"
          label="状态"
          rules={[
            {
              required: true,
              message: '请选择状态',
            },
          ]}
          options={[
            {
              value: '0',
              label: '正常',
            },
            {
              value: '1',
              label: '禁用',
            },
          ]}
        />
      </ModalForm>
    </PageContainer>
  );
};

export default EmailTaskList;
