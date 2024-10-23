declare namespace NumberPool {
  type NumberPoolListItem = {
    ID?: number;
    name?: string;
    description?: string;
    sender_count?: number;
    sender?: [];
    CreatedAt?: Date;
    UpdatedAt?: Date;
  };

  type NumberPoolList = {
    data?: NumberPoolListItem[];
    list?: NumberPoolListItem[];
    /** 列表的内容总数 */
    total?: number;
    currentPage?: number;
    totalPages?: number;
  };
}
