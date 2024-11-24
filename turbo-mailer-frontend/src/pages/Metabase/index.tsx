import React, { useEffect } from 'react';

const Metabase: React.FC = () => {
  useEffect(() => {
    const currentHost = window.location.host;
    if (currentHost === 'localhost:8000' || /^(\d{1,3}\.){3}\d{1,3}(:\d+)?$/.test(currentHost)) {
      window.location.href = 'https://metabase.turbomx.org';
      return;
    }
    const hostParts = currentHost.split('.');
    hostParts[0] = 'metabase';
    const metabaseUrl = `http://${hostParts.join('.')}`;
    window.location.href = metabaseUrl;
  }, []);

  return <div>Metabase redirecting...</div>;
};

export default Metabase;
