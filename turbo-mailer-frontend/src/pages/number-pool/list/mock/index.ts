import Mock from 'mockjs';
import qs from 'query-string';
import dayjs from 'dayjs';
import setupMock from '@/utils/setupMock';

const { list } = Mock.mock({
  'list|100': [
    {
      id: /[0-9]{8}[-][0-9]{4}/,
      name: () =>
        Mock.Random.pick([
          '号池名称一',
          '号池名称二',
          '号池名称三',
        ]),
      'count|0-2000': 0,
      'createdAt|1-60': 0,
      'status|0-1': 0,
    },
  ],
});

const filterData = (
  rest: {
    id?: string;
    name?: string;
    'createdAt[]'?: string[];
    'status[]'?: string;
  } = {}
) => {
  const {
    id,
    name,
    'createdAt[]': createdAt,
    'status[]': status,
  } = rest;
  if (id) {
    return list.filter((item) => item.id === id);
  }
  let result = [...list];
  if (name) {
    result = result.filter((item) => {
      return (item.name as string).toLowerCase().includes(name.toLowerCase());
    });
  }
  if (createdAt && createdAt.length === 2) {
    const [begin, end] = createdAt;
    result = result.filter((item) => {
      const time = dayjs()
        .subtract(item.createdAt, 'days')
        .format('YYYY-MM-DD HH:mm:ss');
      return (
        !dayjs(time).isBefore(dayjs(begin)) && !dayjs(time).isAfter(dayjs(end))
      );
    });
  }

  if (status && status.length) {
    result = result.filter((item) => status.includes(item.status.toString()));
  }

  return result;
};

setupMock({
  setup: () => {
    Mock.mock(new RegExp('/api/list'), (params) => {
      const {
        page = 1,
        pageSize = 10,
        ...rest
      } = qs.parseUrl(params.url).query;
      const p = page as number;
      const ps = pageSize as number;

      const result = filterData(rest);
      return {
        list: result.slice((p - 1) * ps, p * ps),
        total: result.length,
      };
    });
  },
});
