export const cacheKeys = {
  workflowList: () => 'workflow:list',
  workflowGet: (id: string) => `workflow:get:${id}`,
  workflowVersionList: (workflowId: string) =>
    `workflowVersion:list:${workflowId}`,
  workflowVersionGet: (workflowId: string, version: number) =>
    `workflowVersion:get:${workflowId}:${version}`,
};
