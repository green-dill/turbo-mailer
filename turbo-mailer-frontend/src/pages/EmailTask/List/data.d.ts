declare namespace EmailTask {
  type EmailTaskListItem = {
    id?: number;
    title?: string;
    content?: string;
    recipients?: string;
    numberPoolIds?: number[];
    sendInterval?: string;
    taskStartedAt?: Date;
    status?: number;
    createdAt?: Date;
    updatedAt?: Date;
  };

  type EmailTaskList = {
    data?: EmailTaskListItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };
}
