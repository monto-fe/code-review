import React from "react";
import { Tag } from "antd";
import { CheckCircleOutlined, SyncOutlined, CloseCircleOutlined } from "@ant-design/icons";

type StatusTagProps = {
  status: 1 | 2 | 3;
};

const StatusTag: React.FC<StatusTagProps> = ({ status }) => {
  switch (status) {
    case 3:
      return (
        <Tag icon={<CheckCircleOutlined />} color="success">
          success
        </Tag>
      );
    case 2:
      return (
        <Tag icon={<SyncOutlined spin />} color="processing">
          processing
        </Tag>
      );
    case 1:
      return (
        <Tag icon={<CloseCircleOutlined />} color="error">
          error
        </Tag>
      );
    default:
      return null;
  }
};

export default StatusTag;