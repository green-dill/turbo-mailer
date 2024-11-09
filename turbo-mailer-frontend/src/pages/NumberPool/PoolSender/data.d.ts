declare namespace PoolSender {
  type PoolSenderListItem = {
    id?: number;
    pool_id?: string;
    from_name?: string;
    from_email?: string;
    reply_to?: string;
    domain?: string;
    created_at?: Date;
    updated_at?: Date;
  };

  type PoolSenderList = {
    data?: PoolSenderListItem[];
    list?: PoolSenderListItem[];
    /** 列表的内容总数 */
    total?: number;
    currentPage?: number;
    totalPages?: number;
  };
}
