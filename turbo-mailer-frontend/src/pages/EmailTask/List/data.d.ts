declare namespace EmailTask {
  type EmailTaskListItem = {
    ID?: stirng;
    subject?: string;
    context_type?: string;
    content?: string;
    receivers?: [];
    state?: string;
    pools?: [];
    max_dispatch_pre_hour?: number;
    schedule_at?: Date;
    last_dispatch_at?: Date;
    CreatedAt?: Date;
    UpdatedAt?: Date;
  };

  type EmailTaskList = {
    data?: EmailTaskListItem[];
    list?: EmailTaskListItem[];
    /** 列表的内容总数 */
    total?: number;
    currentPage?: number;
    totalPages?: number;
  };
}
