declare namespace Dashboard {
  type Stats = {
    poolCount?: number;
    topPools?: TopPool[];
    taskCount?: number;
    taskStateCount?: TaskStateCount;
  };

  type TopPool = {
    name?: string;
    senderCount?: number;
  };

  type TaskStateCount = {
    pending?: number;
    dispatched?: number;
    finished?: number;
  };
}
