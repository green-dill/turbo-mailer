declare namespace EmailTask {
  import NumberPoolListItem = NumberPool.NumberPoolListItem;
  type EmailTaskListItem = {
    id?: string;
    subject?: string;
    content_type?: string;
    content: File;
    receivers: File;
    content_files: RcFile[];
    receivers_files: RcFile[];
    metadata?: string;
    state?: string;
    pools?: TaskPool[];
    pool_ids?: number[];
    pools_weights?: number[];
    max_dispatch_per_hour?: number;
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
