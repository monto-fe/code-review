import Form from '@/pages/component/Form';
import { ITableFilter } from './data';

function Filters(props: ITableFilter) {
  const { items, size, handleSearch } = props;

  return <Form items={items} size={size} handleSearch={handleSearch} />;
}

export default Filters;
