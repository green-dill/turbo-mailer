import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { addEmailTask, removeEmailTask, queryEmailTask, updateEmailTask } from './service';
import {Button, Drawer, message, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {
  ModalForm, ProFormDateTimePicker, ProFormDigit, ProFormRadio, ProFormText, ProFormTextArea, ProFormUploadButton
} from "@ant-design/pro-components";

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
      await updateEmailTask(fields);
      message.success('编辑成功');
      return true;
    } catch (error) {
      message.error('编辑失败请重试！');
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
      title: '接收者',
      dataIndex: 'receivers',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
    },
    {
      title: '状态',
      dataIndex: 'state',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
      valueEnum: {
        '0': { text: '正常', status: 'blue' },
        '1': { text: '禁用', status: 'default' },
      },
    },
    {
      title: '号池',
      dataIndex: 'pools',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
    },
    {
      title: '发送频率',
      dataIndex: 'max_dispatch_pre_hour',
      hideInSearch: true,
      hideInForm: true,
      sorter: true,
    },
    {
      title: '调度时间',
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
            record.receivers = undefined;
            record.content = undefined;
            setCurrentRow(record);
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
        <ProFormUploadButton
          label='邮件内容'
          placeholder='请上传'
          tooltip='上传 HTML/TXT 文件，支持模板语法'
          name="content"
          fieldProps={{
            beforeUpload(file, fileList) {
              return false;
            },
          }}
        />
        <ProFormUploadButton
          label='收件人列表'
          placeholder='请上传'
          tooltip='上传 CSV/EXCEL 文件'
          name="recipients"
          accept={'.xlsx,.xls,.xlsm,.csv'}
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
          name="maxDispatchPerHour"
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
          name="schedule_at"
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
        <ProFormUploadButton
          label='邮件内容'
          placeholder='重新上传'
          tooltip='上传 HTML/TXT 文件，支持模板语法'
          name="content"
          fieldProps={{
            beforeUpload(file, fileList) {
              return false;
            },
          }}
        />
        <ProFormUploadButton
          label='收件人列表'
          placeholder='请上传'
          tooltip='重新上传 CSV/EXCEL 文件'
          name="recipients"
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
          rules={[
            {
              required: true,
              message: '发送时间不能为空',
            },
          ]}
          width="md"
          name="schedule_at"
        />
        <ProFormRadio.Group
          name="state"
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
