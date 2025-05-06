import Permission from '@/components/Permission';

export default function TestPagePermission() {
  return (
    <div>
      <Permission role={'admin'}>我是admin专属</Permission>
      <p>其余内容</p>
    </div>
  );
}
