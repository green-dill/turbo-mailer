import { GithubOutlined } from '@ant-design/icons';
import { DefaultFooter } from '@ant-design/pro-components';
import React from 'react';

const Footer: React.FC = () => {
  return (
    <DefaultFooter
      style={{
        background: 'none',
      }}
      links={[
        {
          key: 'github',
          title: <GithubOutlined />,
          href: 'https://github.com/green-dill/turbo-mailer',
          blankTarget: true,
        },
        {
          key: 'Turbo Mailer',
          title: 'Turbo Mailer',
          href: 'https://github.com/green-dill/turbo-mailer',
          blankTarget: true,
        },
      ]}
    />
  );
};

export default Footer;
