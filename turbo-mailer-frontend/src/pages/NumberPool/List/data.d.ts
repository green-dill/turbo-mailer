declare namespace NumberPool {
  type NumberPoolListItem = {
    id?: number;
    name?: string;
    description?: string;
    quantity?: number;
    status?: number;
    createdAt?: Date;
    updatedAt?: Date;
  };

  type NumberPoolList = {
    data?: NumberPoolListItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };
}
