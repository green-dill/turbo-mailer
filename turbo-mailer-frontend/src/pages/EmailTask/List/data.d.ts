declare namespace EmailTask {
  type EmailTaskListItem = {
    id?: stirng;
    subject?: string;
    context_type?: string;
    content?: string;
    receivers?: [];
    state?: string;
    pools?: [];
    maxDispatchPerHour?: number;
    schedule_at?: Date;
    last_dispatch_at?: Date;
    created_at?: Date;
    updated_at?: Date;
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
