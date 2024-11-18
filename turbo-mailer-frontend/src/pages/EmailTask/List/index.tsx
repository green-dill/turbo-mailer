import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable, {TableDropdown} from '@ant-design/pro-table';
import {
  addEmailTask,
  removeEmailTask,
  queryEmailTask,
  updateEmailTask,
  updateTest,
  startImmediately,
} from './service';
import {Button, message, Modal, Space, Tag} from "antd";
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
      const formData = new FormData();

      formData.append('pools_weights', fields.pool_ids?.map(id => 1).join(',') || '');
      formData.append('pools', fields.pool_ids?.join(',') || '');

      fields.pools = undefined;
      fields.pool_ids = undefined;

      for (const key in fields) {
        if (fields[key as keyof EmailTask.EmailTaskListItem] !== undefined) {
          const value = fields[key as keyof EmailTask.EmailTaskListItem];
          if (value instanceof Blob || value instanceof File) {
            // 如果值是 Blob 或 File 类型，直接添加
            formData.append(key, value as Blob);
          } else {
            // 如果值是其他类型，转换为字符串然后添加
            formData.append(key, String(value));
          }
        }
      }

      fields.id ? await updateEmailTask(formData, fields.id) : await addEmailTask(formData);
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
      ellipsis: true,
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
            label: 'text/html',
            value: 'text/html',
          },
          {
            label: 'text/plain',
            value: 'text/plain',
          },
        ]
      }
    },
    {
      title: '邮件内容',
      dataIndex: 'content',
      hideInSearch: true,
      hideInTable: true,
      formItemProps: {
        rules: [
          {
            required: currentRow === undefined,
            message: '请上传',
          },
        ],
        name: 'content',
        valuePropName: 'content',
        getValueFromEvent: e => {
          return e.fileList[0]?.originFileObj;
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
                  下载TXT模板
                </a>
              )}
              {content_type === 'text/html' && (
                <a href='/api/v1/tasks/content-template?type=html' target="_blank" rel="noopener noreferrer">
                  下载HTML模板
                </a>
              )}
            </>
          );

          // 根据 content_type 动态设置接受的文件类型
          let acceptTypes;
          if (content_type === 'text/html') {
            acceptTypes = '.html, .eml';
          } else if (content_type === 'text/plain') {
            acceptTypes = '.txt';
          } else {
            // 如果 content_type 不是 'text/html' 或 'text/plain'，则默认接受两种类型
            acceptTypes = '.html, .txt';
          }

          return (
            <ProFormUploadButton
              placeholder='请上传'
              tooltip='上传 HTML/TXT/EML 文件，支持模板语法'
              help={helpContent}
              name="content"
              accept={acceptTypes}
              max={1}
              disabled={content_type === undefined}
              fieldProps={{
                beforeUpload: () => {
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
      dataIndex: 'receivers',
      hideInSearch: true,
      hideInTable: true,
      formItemProps: {
        rules: [
          {
            required: currentRow === undefined,
            message: '请上传',
          },
        ],
        name: 'receivers',
        valuePropName: 'receivers',
        getValueFromEvent: e => {
          return e.fileList[0]?.originFileObj;
        }
      },
      renderFormItem: (_, { type, defaultRender }, form) => {
        if (type === 'form') {
          return (
            <ProFormUploadButton
              placeholder='请上传'
              tooltip='上传 CSV/TXT 文件'
              help={<>需要帮助？<a href="/api/v1/tasks/receivers-template" target="_blank" rel="noopener noreferrer">下载CSV模板</a></>}
              name="receivers"
              accept={'.csv, .txt'}
              max={1}
              fieldProps={{
                beforeUpload: () => {
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
      render: (_, record) => (
        <Space>
          {record?.pools?.map(({ pool }) => (
            <Tag key={pool.id}>
              {pool.name}
            </Tag>
          ))}
        </Space>
      ),
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
      render: (text, record: EmailTask.EmailTaskListItem, _, action) => [
        <Button
          type="link"
          size="small"
          style={{padding: 0}}
          key='edit'
          disabled={record?.state === 'finished'}
          onClick={() => {
            handleModalOpen(true);
            record.pool_ids = record.pools && Object.values(record.pools.map((item) => item.pool_id));
            record.schedule_at = moment(record.schedule_at).format('YYYY-MM-DD HH:mm:ss');
            record.content = undefined;
            record.receivers = undefined;
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
          disabled={record?.state === 'finished'}
          onClick={() => {
            handleTestModalOpen(true);
            setCurrentRow(record);
          }}
        >
          测试
        </Button>,
        <TableDropdown
          key="actionGroup"
          menus={[
            {
              key: 'start',
              name: '执行',
              disabled: record?.state === 'finished',
              onClick: async () => {
                const success = await handleStartImmediately(record);
                if (success) {
                  if (actionRef.current) {
                    actionRef.current.reload();
                  }
                }
            }},
            {
              key: 'delete',
              name: '删除',
              disabled: record?.state === 'dispatched',
              onClick: (e) => {
                Modal.confirm({
                  title: '确定删除吗？',
                  content: '此操作将永久删除该记录，是否继续？',
                  okText: '确定',
                  cancelText: '取消',
                  onOk: async () => {
                    const success = await handleRemove(record);
                    if (success) {
                      if (actionRef.current) {
                        actionRef.current.reload();
                      }
                    }
                  },
                });
              },
            },
          ]}
        />,
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
        columnsState={{
          persistenceKey: 'email-task-list',
          persistenceType: 'localStorage',
          defaultValue: {
            created_at: {show: false},
            updated_at: {show: false},
            option: {fixed: 'right', disable: true},
          },
        }}
        columns={columns}
      />
      <Modal
        title={`${currentRow?.id ? '更新' : '新建'}任务`}
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
