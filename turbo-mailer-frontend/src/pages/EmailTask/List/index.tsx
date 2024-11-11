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
import {Button, message, Modal, Popconfirm, Tag} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {
  ModalForm,
  ProFormText, ProFormUploadButton,
} from "@ant-design/pro-components";
import {querySimpleNumberPool} from "@/pages/NumberPool/List/service";
import moment from "moment";

const EmailTaskList: React.FC = () => {
  const [modalOpen, handleModalOpen] = useState<boolean>(false);
  const [testModalOpen, handleTestModalOpen] = useState<boolean>(false);
  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<EmailTask.EmailTaskListItem>();

  /**
   * 添加或修改邮件任务
   * @param formData 邮件任务
   */
  const handleSubmit = async (fields: EmailTask.EmailTaskListItem) => {
    try {
      fields.id ? await updateEmailTask(fields) : await addEmailTask(fields);
      message.success(`${fields.id ? '编辑' : '添加'}成功`);
      return true;
    } catch (error) {
      message.error(`${fields.id ? '编辑' : '添加'}失败请重试!`);
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
      tooltip: '支持模板语法',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '请输入',
          },
        ],
      }
    },
    {
      title: '内容类型',
      dataIndex: 'content_type',
      hideInSearch: true,
      valueType: "radio",
      formItemProps: {
        rules: [
          {
            required: true,
            message: '请输入',
          },
        ]
      },
      fieldProps: {
        options: [
          {
            label: 'html',
            value: 'text/html',
          },
          {
            label: 'text',
            value: 'text/plain',
          },
        ]
      }
    },
    {
      title: '邮件内容',
      dataIndex: 'content_files',
      hideInSearch: true,
      hideInTable: true,
      tooltip: '上传 HTML/TXT 文件，支持模板语法',
      formItemProps: {
        rules: [
          {
            required: currentRow === undefined,
            message: '请上传',
          },
        ],
        name: 'content_files',
        valuePropName: 'content_files',
        getValueFromEvent: e => {
          return e && e.fileList;
        }
      },
      renderFormItem: (_, { type, defaultRender }, form) => {
        if (type === 'form') {
          const content_type = form.getFieldValue('content_type');

          const helpContent = (
            content_type && <>
              需要帮助？
              {content_type === 'text/plain' && (
                <a href='/api/v1/tasks/content-template?type=text' target="_blank" rel="noopener noreferrer">
                  下载模板
                </a>
              )}
              {content_type === 'text/html' && (
                <a href='/api/v1/tasks/content-template?type=html' target="_blank" rel="noopener noreferrer">
                  下载模板
                </a>
              )}
            </>
          );

          // 根据 content_type 动态设置接受的文件类型
          let acceptTypes;
          if (content_type === 'text/html') {
            acceptTypes = '.html';
          } else if (content_type === 'text/plain') {
            acceptTypes = '.txt';
          } else {
            // 如果 content_type 不是 'text/html' 或 'text/plain'，则默认接受两种类型
            acceptTypes = '.html, .txt';
          }

          return (
            <ProFormUploadButton
              placeholder='请上传'
              tooltip='上传 HTML/TXT 文件，支持模板语法'
              help={helpContent}
              name="content_files"
              accept={acceptTypes}
              max={1}
              disabled={content_type === undefined}
              fieldProps={{
                beforeUpload(file, fileList) {
                  return false;
                },
              }}
            />
          );
        }
        return defaultRender(_);
      },
    },
    {
      title: '收件人列表',
      dataIndex: 'receivers_files',
      hideInSearch: true,
      hideInTable: true,
      tooltip: '上传 CSV/EXCEL 文件',
      formItemProps: {
        rules: [
          {
            required: currentRow === undefined,
            message: '请上传',
          },
        ],
        name: 'receivers_files',
        valuePropName: 'receivers_files',
        getValueFromEvent: e => {
          return e && e.fileList;
        }
      },
      renderFormItem: (_, { type, defaultRender }, form) => {
        if (type === 'form') {
          return (
            <ProFormUploadButton
              placeholder='请上传'
              tooltip='上传 CSV/EXCEL 文件'
              help={<>需要帮助？<a href="/api/v1/tasks/receivers-template" target="_blank" rel="noopener noreferrer">下载模板</a></>}
              name="receivers_files"
              accept={'.csv'}
              max={1}
              fieldProps={{
                beforeUpload(file, fileList) {
                  return false;
                },
              }}
            />
          );
        }
        return defaultRender(_);
      },
    },
    {
      title: '号池',
      dataIndex: 'pools',
      hideInSearch: true,
      hideInForm: true,
      render: (text, record, _, action) => {
        // 假设 pools 是一个数组，每个元素是一个对象，对象中有一个 name 属性
        if (Array.isArray(text)) {
          return (
            <>
              {text.map((pool, index) => (
                <Tag key={index}>{pool.pool.name}</Tag> // 使用 tag 标签包裹 pool.name，并添加 key 属性以避免警告
              ))}
            </>
          );
        }
        return '无'; // 如果 pools 不是数组，可以返回一个默认值
      },
    },
    {
      title: '号池',
      dataIndex: 'pool_ids',
      hideInSearch: true,
      hideInTable: true,
      valueType: 'select',
      request: querySimpleNumberPool,
      params: {current: 1, pageSize: 1000},
      fieldProps: {
        mode: 'multiple',
        disabled: currentRow !== undefined,
      },
      formItemProps: {
        rules: [
          {
            required: true,
            message: '请选择',
          },
        ],
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
      dataIndex: 'max_dispatch_per_hour',
      tooltip: '每小时最多发送邮件的数量',
      hideInSearch: true,
      sorter: true,
      valueType: 'digit',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '请输入',
          },
        ],
      },
    },
    {
      title: '计划时间',
      dataIndex: 'schedule_at',
      valueType: 'dateTime',
      tooltip: '计划发送邮件的时间',
      hideInSearch: true,
      sorter: true,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '请选择',
          },
        ],
      },
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
            handleModalOpen(true);
            record.pool_ids = record.pools && Object.values(record.pools.map((item) => item.pool_id));
            record.schedule_at = moment(record.schedule_at).format('YYYY-MM-DD HH:mm:ss')
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
          <Button type="primary" onClick={() => {
            handleModalOpen(true);
            setCurrentRow(undefined);
          }}>
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={queryEmailTask}
        columns={columns}
      />
      <Modal
        title={`${currentRow?.id ? '更新' : '新建'}任务s`}
        open={modalOpen}
        onCancel={() => {
          handleModalOpen(false);
          setCurrentRow(undefined);
        }}
        footer={null}
        destroyOnClose // 确保弹窗关闭时子组件被销毁
      >
        <ProTable<EmailTask.EmailTaskListItem, EmailTask.EmailTaskListItem>
          onSubmit={async (fields) => {
            fields.id = currentRow?.id;

            if (Array.isArray(fields.pool_ids)) {
              fields.pools_weights = fields.pool_ids.map(() => 1);
            }

            if (Array.isArray(fields.receivers_files)) {
              fields.receivers = fields.receivers_files[0].originFileObj
              fields.receivers_files = undefined;
            }

            if (Array.isArray(fields.content_files)) {
              fields.content = fields.content_files[0].originFileObj
              fields.content_files = undefined;
            }
            const success = await handleSubmit(fields);
            if (success) {
              handleModalOpen(false);
              setCurrentRow(undefined);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          rowKey="id"
          type="form"
          columns={columns}
          form={{
            initialValues: currentRow,
          }}
        />
      </Modal>
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
            {
              type: 'email',
              message: '电子邮箱格式不正确',
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
