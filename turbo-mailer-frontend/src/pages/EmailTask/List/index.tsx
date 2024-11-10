import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {
  addEmailTask,
  removeEmailTask,
  queryEmailTask,
  updateEmailTask,
  updateTest,
  startImmediately,
} from './service';
import {Button, message, Popconfirm, Tag} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {
  ModalForm,
  ProFormDateTimePicker,
  ProFormDigit,
  ProFormRadio,
  ProFormSelect,
  ProFormText,
  ProFormUploadButton
} from "@ant-design/pro-components";
import {querySimpleNumberPool} from "@/pages/NumberPool/List/service";

const EmailTaskList: React.FC = () => {
  /**
   * @en-US Pop-up window of new window
   * @zh-CN 新建窗口的弹窗
   *  */
  const [createModalOpen, handleModalOpen] = useState<boolean>(false);
  /**
   * @en-US The pop-up window of the distribution update window
   * @zh-CN 分布更新窗口的弹窗
   * */
  const [updateModalOpen, handleUpdateModalOpen] = useState<boolean>(false);
  const [testModalOpen, handleTestModalOpen] = useState<boolean>(false);
  const [contentType, setContentType] = useState<string>();

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<EmailTask.EmailTaskListItem>();

  /**
   * 添加邮件任务
   * @param fields 邮件任务
   */
  const handleAdd = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      await addEmailTask(fields);
      message.success('添加成功');
      return true;
    } catch (error) {
      message.error('添加失败请重试！');
      return false;
    }
  };

  /**
   * 修改邮件任务
   * @param fields 邮件任务
   */
  const handleUpdate = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      fields.id = currentRow?.id;
      // 检查 pool_ids 是否存在且为数组
      if (Array.isArray(fields.pool_ids)) {
        // 使用 map 方法创建一个新的数组，其长度与 pool_ids 相同，所有元素都为 1
        fields.pools_weights = fields.pool_ids.map(() => 1);
      }
      await updateEmailTask(fields);
      message.success('编辑成功');
      return true;
    } catch (error) {
      message.error('编辑失败请重试！');
      return false;
    }
  };

  /**
   * 修改邮件任务
   * @param fields 邮件任务
   */
  const handleTest = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      fields.id = currentRow?.id;
      await updateTest(fields);
      message.success('测试邮件发送成功');
      return true;
    } catch (error) {
      message.error('测试邮件发送失败请重试！');
      return false;
    }
  };

  /**
   * 立即执行任务
   * @param fields 邮件任务
   */
  const handleStartImmediately = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      await startImmediately(fields);
      message.success('执行成功');
      return true;
    } catch (error) {
      message.error('执行失败请重试！');
      return false;
    }
  };

  /**
   *  删除邮件任务
   * @param fields 邮件任务
   */
  const handleRemove = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      await removeEmailTask(fields);
      message.success('删除成功');
      return true;
    } catch (error) {
      message.error('删除失败请重试！');
      return false;
    }
  };

  const columns: ProColumns<EmailTask.EmailTaskListItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      hideInForm: true,
      hideInTable: true,
    },
    {
      title: '标题',
      dataIndex: 'subject',
    },
    {
      title: '类型',
      dataIndex: 'content_type',
      hideInSearch: true,
    },
    {
      title: '号池',
      dataIndex: 'pools',
      hideInSearch: true,
      render: (text, record, _, action) => {
        // 假设 pools 是一个数组，每个元素是一个对象，对象中有一个 name 属性
        if (Array.isArray(text)) {
          return text.map(pool => pool.pool.name).join(', '); // 将名称用逗号分隔
        }
        return '无'; // 如果 pools 不是数组，可以返回一个默认值
      },
    },
    {
      title: '状态',
      dataIndex: 'state',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
      valueEnum: {
        'pending': { text: '等待', status: 'Default' },
        'dispatched': { text: '已调度', status: 'Processing' },
        'finished': { text: '完成', status: 'Success' },
      },
    },
    {
      title: '发送频率',
      dataIndex: 'max_dispatch_pre_hour',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
    },
    {
      title: '计划时间',
      dataIndex: 'schedule_at',
      valueType: 'dateTime',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
    },
    {
      title: '上次调度时间',
      dataIndex: 'last_dispatch_at',
      valueType: 'dateTime',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'date',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
      defaultSortOrder: 'descend',
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'date',
      hideInSearch: true,
      hideInForm: true,
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
            handleUpdateModalOpen(true);
            record.pool_ids = record.pools && Object.values(record.pools.map((item) => item.pool_id));
            record.receivers = undefined;
            record.content = undefined;
            setCurrentRow(record);
          }}
        >
          编辑
        </Button>,
        <Button
          type="link"
          size="small"
          style={{padding: 0}}
          key='test'
          onClick={() => {
            handleTestModalOpen(true);
            setCurrentRow(record);
          }}
        >
          测试
        </Button>,
        <Popconfirm
          key='start'
          title="确定执行吗？"
          onConfirm={async () => {
            const success = await handleStartImmediately(record);
            if (success) {
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
        >
          <Button style={{padding: 0}} type="link" size="small">执行</Button>
        </Popconfirm>,
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
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<EmailTask.EmailTaskListItem, API.PageParams>
        headerTitle="任务列表"
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
        title="添加邮件任务"
        width="400px"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        onFinish={async (value) => {
          const success = await handleAdd(value as EmailTask.EmailTaskListItem);
          if (success) {
            handleModalOpen(false);
            setCurrentRow(undefined);
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
              message: '请输入邮件标题',
            },
          ]}
          width="md"
          name="subject"
        />
        <ProFormRadio.Group
          name="content_type"
          width="md"
          label="内容类型"
          rules={[
            {
              required: true,
              message: '请选择内容类型',
            },
          ]}
          options={[
            {
              label: 'HTML',
              value: 'text/html',
            },
            {
              label: 'TXT',
              value: 'text/plain',
            },
          ]}
          fieldProps={{
            onChange: (e) => {
              setContentType(e.target.value)
            },
          }}
        />
        <ProFormUploadButton
          label='邮件内容'
          placeholder='请上传'
          tooltip='上传 HTML/TXT 文件，支持模板语法'
          help={contentType && <>需要帮助？<a href={`/api/v1/tasks/content-template?type=${contentType == 'text/plain' ? 'text' : contentType == 'text/html' ? 'html' : undefined}`} target="_blank" rel="noopener noreferrer">下载模板</a></>}
          name="content"
          accept={'.html, .txt'}
          fieldProps={{
            beforeUpload(file, fileList) {
              return false;
            },
            disabled: contentType === undefined,
          }}
        />
        <ProFormUploadButton
          label='收件人列表'
          placeholder='请上传'
          tooltip='上传 CSV/EXCEL 文件'
          help={<>需要帮助？<a href="/api/v1/tasks/recipients-template" target="_blank" rel="noopener noreferrer">下载模板</a></>}
          name="recipients"
          accept={'.csv'}
          fieldProps={{
            beforeUpload(file, fileList) {
              return false;
            },
          }}
        />
        <ProFormDigit
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
          name="max_dispatch_per_hour"
        />
        <ProFormDateTimePicker
          label='发送时间'
          placeholder='请选择'
          tooltip='邮件开始发送的时间'
          width="md"
          name="schedule_at"
        />
        <ProFormSelect
          label='号池'
          placeholder='请选择'
          mode="multiple"
          name="pools"
          allowClear
          width="md"
          request={querySimpleNumberPool}
          params={{current: 1, pageSize: 1000}}
        />
      </ModalForm>
      <ModalForm
        title="编辑邮件任务"
        width="400px"
        open={updateModalOpen}
        onOpenChange={handleUpdateModalOpen}
        initialValues={currentRow} // 设置初始值
        onFinish={async (value) => {
          const success = await handleUpdate(value as EmailTask.EmailTaskListItem);
          if (success) {
            handleUpdateModalOpen(false);
            setCurrentRow(undefined);
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
              message: '请输入邮件标题',
            },
          ]}
          width="md"
          name="subject"
        />
        <ProFormRadio.Group
          name="content_type"
          width="md"
          label="内容类型"
          rules={[
            {
              required: true,
              message: '请选择内容类型',
            },
          ]}
          options={[
            {
              label: 'HTML',
              value: 'text/html',
            },
            {
              label: 'TXT',
              value: 'text/plain',
            },
          ]}
        />
        <ProFormUploadButton
          label='邮件内容'
          placeholder='请上传'
          tooltip='上传 HTML/TXT 文件，支持模板语法'
          help={currentRow?.content_type && <>需要帮助？<a href={`/api/v1/tasks/content-template?type=${currentRow?.content_type == 'text/plain' ? 'text' : currentRow.content_type == 'text/html' ? 'html' : undefined}`} target="_blank" rel="noopener noreferrer">下载模板</a></>}
          name="content"
          accept={'.html, .txt'}
          fieldProps={{
            beforeUpload(file, fileList) {
              return false;
            },
            disabled: contentType === undefined,
          }}
        />
        <ProFormUploadButton
          label='收件人列表'
          placeholder='请上传'
          tooltip='重新上传 CSV/EXCEL 文件'
          help={<>需要帮助？<a href="/api/v1/tasks/receivers-template" target="_blank" rel="noopener noreferrer">下载模板</a></>}
          name="recipients"
          accept={'.csv'}
          fieldProps={{
            beforeUpload(file, fileList) {
              return false;
            },
          }}
        />
        <ProFormDigit
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
          name="max_dispatch_pre_hour"
        />
        <ProFormDateTimePicker
          label='发送时间'
          placeholder='请选择'
          tooltip='邮件开始发送的时间'
          width="md"
          name="schedule_at"
        />
        <ProFormSelect
          label='号池'
          placeholder='请选择'
          mode="multiple"
          name='pool_ids'
          allowClear
          width="md"
          request={querySimpleNumberPool}
          params={{current: 1, pageSize: 1000}}
        />
      </ModalForm>
      <ModalForm
        title="发送测试邮件"
        width="400px"
        open={testModalOpen}
        onOpenChange={handleTestModalOpen}
        initialValues={currentRow} // 设置初始值
        onFinish={async (value) => {
          const success = await handleTest(value as EmailTask.EmailTaskListItem);
          if (success) {
            handleTestModalOpen(false);
            setCurrentRow(undefined);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          label='邮件地址'
          placeholder='请输入'
          tooltip='接收测试邮件的地址'
          rules={[
            {
              required: true,
              message: '请输入邮件地址',
            },
          ]}
          width="md"
          name="email"
        />
      </ModalForm>
    </PageContainer>
  );
};

export default EmailTaskList;
