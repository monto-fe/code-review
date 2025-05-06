import { Button } from 'antd';
import { queryList as queryUserList } from '@/pages/acl/user/List/service';

export default function TestAPIPermission() {
  const fetchAPI = () => {
    queryUserList()
      .then(() => {
        alert('调用成功');
      })
      .catch(() => {
        alert('失败了');
      });
  };

  return (
    <div>
      获取用户列表只有admin用户才调用哦~~：
      <Button onClick={fetchAPI}>点我获取用户列表</Button>
    </div>
  );
}
