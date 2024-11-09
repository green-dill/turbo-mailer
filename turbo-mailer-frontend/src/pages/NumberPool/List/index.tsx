import React, {useRef, useState} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import { addNumberPool, removeNumberPool, queryNumberPool, updateNumberPool, } from './service';
import {Button, Drawer, message, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import {
  ModalForm, ProFormText, ProFormTextArea,
} from "@ant-design/pro-components";
import {Link} from "@@/exports";
import {useParams} from "react-router";

const NumberPoolList: React.FC = () => {
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
  const [currentRow, setCurrentRow] = useState<NumberPool.NumberPoolListItem>();

  /**
   * 添加号池
   * @param fields 号池
   */
  const handleAdd = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      await addNumberPool(fields);
      message.success('添加成功');
      return true;
    } catch (error) {
      message.error('添加失败请重试！');
      return false;
    }
  };

  /**
   * 修改号池
   * @param fields 号池
   */
  const handleUpdate = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      await updateNumberPool(fields);
      message.success('编辑成功');
      return true;
    } catch (error) {
      message.error('编辑失败请重试！');
      return false;
    }
  };

  /**
   *  删除号池
   * @param fields 号池
   */
  const handleRemove = async (fields: NumberPool.NumberPoolListItem) => {
    try {
      await removeNumberPool(fields);
      message.success('删除成功');
      return true;
    } catch (error) {
      message.error('删除失败请重试！');
      return false;
    }
  };

  const columns: ProColumns<NumberPool.NumberPoolListItem>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      hideInForm: true,
      hideInTable: true,
    },
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      hideInSearch: true,
    },
    {
      title: '邮箱数量',
      dataIndex: 'sender_count',
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
      render: (_, record: NumberPool.NumberPoolListItem) => [
        <Link key="setting" to={`/number-pool/list/${record.id}`}>
          配置
        </Link>,
        <Button
          type="link"
          size="small"
          style={{padding: 0}}
          key='edit'
          onClick={() => {
            handleUpdateModalOpen(true);
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
      <ProTable<NumberPool.NumberPoolListItem, API.PageParams>
        headerTitle="号池列表"
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
        request={queryNumberPool}
        columns={columns}
      />
      <ModalForm
        title="添加号池"
        width="400px"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        onFinish={async (value) => {
          const success = await handleAdd(value as NumberPool.NumberPoolListItem);
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
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="name"
          label="号池名称"
        />
        <ProFormTextArea name="description" label="号池描述" />
      </ModalForm>
      <ModalForm
        title="编辑号池"
        width="400px"
        open={updateModalOpen}
        onOpenChange={handleUpdateModalOpen}
        initialValues={currentRow} // 设置初始值
        onFinish={async (value) => {
          const success = await handleUpdate(value as NumberPool.NumberPoolListItem);
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
          rules={[
            {
              required: true,
              message: '请输入',
            },
          ]}
          name="name"
          label="号池名称"
        />
        <ProFormTextArea name="description" label="号池描述" />
      </ModalForm>
    </PageContainer>
  );
};

export default NumberPoolList;
