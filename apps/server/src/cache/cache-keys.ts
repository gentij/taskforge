export const cacheKeys = {
  workflowGet: (id: string) => `workflow:get:${id}`,
  workflowVersionGet: (workflowId: string, version: number) =>
    `workflowVersion:get:${workflowId}:${version}`,
};
