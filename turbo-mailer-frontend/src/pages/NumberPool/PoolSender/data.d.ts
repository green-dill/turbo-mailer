declare namespace PoolSender {
  type PoolSenderListItem = {
    ID?: number;
    pool_id?: string;
    from_name?: string;
    from_email?: string;
    reply_to?: string;
    domain?: string;
    CreatedAt?: Date;
    UpdatedAt?: Date;
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
