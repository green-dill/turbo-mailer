declare namespace NumberPool {
  type NumberPoolListItem = {
    id?: stirng;
    name?: string;
    description?: string;
    sender_count?: number;
    sender?: [];
    created_at?: Date;
    updated_at?: Date;
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
