import { PageContainer } from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import {Card, Table, theme} from 'antd';
import React from 'react';

const EmailServiceGuide: React.FC = () => {
  const { token } = theme.useToken();
  const { initialState } = useModel('@@initialState');

  const dataSource = [
    {
      key: '1',
      domain: 'example.com',
      type: 'TXT',
      value: 'v=spf1 mx a ip4:192.0.2.1 a:spf.protection.example.com -all',
    },
    {
      key: '2',
      domain: 'example.net',
      type: 'TXT',
      value: 'google-site-verification=abcdefghijklmnopqrstuvwxyz1234567890',
    },
    {
      key: '1',
      domain: 'example.org',
      type: 'TXT',
      value: 'dmarc=v=DMARC1; p=none;',
    },
  ];

  const columns = [
    {
      title: '域名',
      dataIndex: 'domain',
      key: 'domain',
    },
    {
      title: '记录类型',
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: '记录值',
      dataIndex: 'value',
      key: 'value',
    },
  ];

  return (
    <PageContainer>
      <Card
        style={{
          borderRadius: 8,
        }}
        styles={{
          body: {
            backgroundImage:
              initialState?.settings?.navTheme === 'realDark'
                ? 'background-image: linear-gradient(75deg, #1A1B1F 0%, #191C1F 100%)'
                : 'background-image: linear-gradient(75deg, #FBFDFF 0%, #F5F7FF 100%)',
          },
        }}
      >
        <div
          style={{
            backgroundPosition: '100% -30%',
            backgroundRepeat: 'no-repeat',
            backgroundSize: '274px auto',
            backgroundImage:
              "url('https://gw.alipayobjects.com/mdn/rms_a9745b/afts/img/A*BuFmQqsB2iAAAAAAAAAAAAAAARQnAQ')",
          }}
        >
          <div
            style={{
              fontSize: '20px',
              color: token.colorTextHeading,
            }}
          >
            邮件发送系统 - DNS 配置指南
          </div>
          <p
            style={{
              fontSize: '14px',
              color: token.colorTextSecondary,
              lineHeight: '22px',
              marginTop: 16,
              marginBottom: 32,
              width: '65%',
            }}
          >
            确保您的域名正确配置以优化邮件发送服务。
          </p>

          <Table dataSource={dataSource} columns={columns} />
    </div>
      </Card>
    </PageContainer>
  );
};

export default EmailServiceGuide;
