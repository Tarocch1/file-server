import React, { useCallback, useState, useEffect } from 'react';
import { useMount, useHash, useLatest } from 'react-use';
import { Breadcrumb, Table, message } from 'antd';
import { HomeOutlined, FolderTwoTone, FileTwoTone } from '@ant-design/icons';
import dayjs from 'dayjs';
import { getList as getListService } from './service';

function App() {
  const [hash, setHash] = useHash();
  const latestHash = useLatest(hash);
  const [list, setList] = useState([]);
  const [loading, setLoading] = useState(false);
  useMount(() => {
    setHash('#/');
  });
  useEffect(getList, [hash]);
  const renderBreadcrumbItem = useCallback(
    function () {
      const paths = hash.split('/').slice(0, -1);
      let href = '';
      return paths.map(function (path, i) {
        href += `${path}/`;
        return (
          <Breadcrumb.Item key={i} href={href}>
            {path === '#' ? <HomeOutlined /> : path}
          </Breadcrumb.Item>
        );
      });
    },
    [hash],
  );
  function getList() {
    const path = latestHash.current.replace('#', '.');
    setLoading(true);
    getListService(path).then(function (res) {
      if (!res.erred) {
        setList(res.data);
      } else {
        message.error(res.message);
      }
      setLoading(false);
    });
  }
  const columns = [
    {
      title: '',
      dataIndex: 'isDir',
      render: isDir => (isDir ? <FolderTwoTone /> : <FileTwoTone />),
      width: 14,
    },
    {
      title: '名称',
      dataIndex: 'name',
      render: (name, record) => {
        if (record.isDir) {
          return <a href={`${hash}${name}/`}>{name}</a>;
        }
        return name;
      },
    },
    {
      title: '修改日期',
      dataIndex: 'time',
      render: time => dayjs(time * 1000).format('YYYY-MM-DD HH:mm:ss'),
      width: 200,
    },
    {
      title: '大小',
      dataIndex: 'size',
      width: 100,
    },
    {
      title: '操作',
      render: record => {},
    },
  ];
  return (
    <div className="wrap">
      <Breadcrumb>{renderBreadcrumbItem()}</Breadcrumb>
      <Table
        style={{ marginTop: 16 }}
        loading={loading}
        columns={columns}
        dataSource={list}
        pagination={false}
        bordered
        size="small"
      />
    </div>
  );
}

export default App;
