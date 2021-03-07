import React, { useCallback, useState, useEffect, useRef } from 'react';
import { useMount, useHash, useLatest } from 'react-use';
import {
  Row,
  Col,
  Breadcrumb,
  Button,
  Table,
  Divider,
  Popconfirm,
  Drawer,
  Typography,
  Progress,
  message,
} from 'antd';
import { HomeOutlined, FolderTwoTone, FileTwoTone } from '@ant-design/icons';
import dayjs from 'dayjs';
import { getList as getListService, remove as removeService } from './service';
import Uploader from './uploader';

function App() {
  const inputEl = useRef(null);
  const [hash, setHash] = useHash();
  const latestHash = useLatest(hash);
  const [list, setList] = useState([]);
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadList, setUploadList] = useState([]);
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
  function download(file) {
    const path = `${latestHash.current.replace('#', '.')}${file.name}`;
    window.open(`/api/download?path=${path}`, '_blank');
  }
  function remove(file) {
    const path = `${latestHash.current.replace('#', '.')}${file.name}`;
    removeService(path).then(function (res) {
      if (!res.erred) {
        message.success('删除成功');
        getList();
      } else {
        message.error(res.message);
      }
    });
  }
  function openFileSelector() {
    if (inputEl.current) {
      inputEl.current.click();
    }
  }
  function startUpload() {
    const path = `${latestHash.current.replace('#', '.')}`;
    const files = inputEl.current.files;
    console.log(files);
    if (files.length === 0) return;
    setUploading(true);
    [...files].forEach(function (file, i) {
      setUploadList(function (_list) {
        const list = [..._list];
        const id = list.length;
        const uploader = new Uploader({
          file,
          path: `${path}${file.name}`,
          onSuccess: function () {
            setUploadList(function (_list) {
              const list = [..._list];
              list[id].percent = 100;
              list[id].status = 'success';
              if (list[id + 1]) {
                list[id + 1].uploader.start();
              } else {
                setUploading(false);
                getList();
              }
              return list;
            });
          },
          onError: function (msg) {
            message.error(msg);
            setUploadList(function (_list) {
              const list = [..._list];
              list[id].status = 'exception';
              if (list[id + 1]) {
                list[id + 1].uploader.start();
              } else {
                setUploading(false);
                getList();
              }
              return list;
            });
          },
          onPregress: function (e) {
            setUploadList(function (_list) {
              const list = [..._list];
              list[id].percent = (e.loaded * 100) / file.size;
              return list;
            });
          },
        });
        if (i === 0) {
          uploader.start();
        }
        list.push({
          name: file.name,
          percent: 0,
          status: 'active',
          uploader,
        });
        return list;
      });
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
      render: record => (
        <React.Fragment>
          <Popconfirm
            title={`是否要删除该${record.isDir ? '文件夹' : '文件'}？`}
            onConfirm={function () {
              remove(record);
            }}
          >
            <a>删除</a>
          </Popconfirm>
          {!record.isDir && (
            <React.Fragment>
              <Divider type="vertical" />
              <a
                onClick={function () {
                  download(record);
                }}
              >
                下载
              </a>
            </React.Fragment>
          )}
        </React.Fragment>
      ),
      width: 100,
    },
  ];
  return (
    <div className="wrap">
      <Row align="middle" gutter={16}>
        <Col style={{ flexGrow: 1 }}>
          <Breadcrumb>{renderBreadcrumbItem()}</Breadcrumb>
        </Col>
        <Col style={{ flexShrink: 0 }}>
          <Button onClick={openFileSelector}>上传文件</Button>
        </Col>
      </Row>

      <Table
        style={{ marginTop: 16 }}
        rowKey="name"
        loading={loading}
        columns={columns}
        dataSource={list}
        pagination={false}
        bordered
        size="small"
      />
      <input
        ref={inputEl}
        type="file"
        style={{ display: 'none' }}
        multiple
        onChange={startUpload}
      />
      <Drawer
        title="上传列表"
        visible={uploading}
        maskClosable={false}
        closable={false}
        width={400}
      >
        {uploadList.map((item, i) => (
          <div key={i} style={{ marginBottom: 16 }}>
            <Typography.Text
              style={{ width: '100%' }}
              ellipsis
              title={item.name}
            >
              {item.name}
            </Typography.Text>
            <Progress percent={item.percent} status={item.status} />
          </div>
        ))}
      </Drawer>
    </div>
  );
}

export default App;
