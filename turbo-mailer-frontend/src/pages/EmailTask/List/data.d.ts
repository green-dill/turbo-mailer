declare namespace EmailTask {
  import NumberPoolListItem = NumberPool.NumberPoolListItem;
  type EmailTaskListItem = {
    id?: stirng;
    subject?: string;
    content_type?: string;
    content?: string;
    receivers?: [];
    metadata?: string;
    state?: string;
    pools?: TaskPool[];
    pool_ids?: number[];
    pools_weights?: number[];
    max_dispatch_pre_hour?: number;
    schedule_at?: string;
    last_dispatch_at?: Date;
    created_at?: Date;
    updated_at?: Date;
  };

  type TaskPool = {
    pool_id: number;
    pool_name: string;
    pool: NumberPoolListItem;
  }

  type EmailTaskList = {
    data?: EmailTaskListItem[];
    list?: EmailTaskListItem[];
    /** 列表的内容总数 */
    total?: number;
    currentPage?: number;
    totalPages?: number;
  };
}
