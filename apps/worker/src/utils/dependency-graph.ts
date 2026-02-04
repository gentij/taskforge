interface StepDefinition {
  key: string;
  dependsOn?: string[];
  request?: Record<string, unknown>;
  input?: Record<string, unknown>;
}

interface DependencyGraph {
  dependencies: Map<string, string[]>;
  executionBatches: string[][];
}

export function buildDependencyGraph(steps: StepDefinition[]): DependencyGraph {
  const dependencies = new Map<string, string[]>();
  const allStepKeys = new Set(steps.map((s) => s.key));
  const inDegree = new Map<string, number>();
  const graph = new Map<string, string[]>();

  // Initialize
  for (const step of steps) {
    dependencies.set(step.key, [...(step.dependsOn || [])]);
    inDegree.set(step.key, 0);
    graph.set(step.key, []);
  }

  // Detect template-based dependencies
  const templatePattern = /\{\{steps\.([a-zA-Z0-9_-]+)\./g;

  for (const step of steps) {
    if (!step.request) continue;

    const requestStr = JSON.stringify(step.request);
    let match;

    while ((match = templatePattern.exec(requestStr)) !== null) {
      const referencedStep = match[1];

      // Skip self-references
      if (referencedStep === step.key) continue;

      // Only add if the referenced step exists in our workflow
      if (allStepKeys.has(referencedStep)) {
        const existing = dependencies.get(step.key) || [];
        if (!existing.includes(referencedStep)) {
          existing.push(referencedStep);
          dependencies.set(step.key, existing);
        }
      }
    }
  }

  // Build graph for topological sort
  for (const [stepKey, deps] of dependencies) {
    for (const dep of deps) {
      if (allStepKeys.has(dep)) {
        const children = graph.get(dep) || [];
        children.push(stepKey);
        graph.set(dep, children);
        inDegree.set(stepKey, (inDegree.get(stepKey) || 0) + 1);
      }
    }
  }

  // Kahn's algorithm for topological sort (batched for parallel execution)
  const queue: string[] = [];
  for (const [key, degree] of inDegree) {
    if (degree === 0) {
      queue.push(key);
    }
  }

  const executionBatches: string[][] = [];

  while (queue.length > 0) {
    const batch = [...queue];
    executionBatches.push(batch);
    queue.length = 0;

    for (const stepKey of batch) {
      const children = graph.get(stepKey) || [];
      for (const child of children) {
        const newDegree = (inDegree.get(child) || 0) - 1;
        inDegree.set(child, newDegree);
        if (newDegree === 0) {
          queue.push(child);
        }
      }
    }
  }

  // Check for cycles
  if (executionBatches.length !== steps.length) {
    throw new Error(
      `Dependency cycle detected. Steps: ${steps.map((s) => s.key).join(', ')}`
    );
  }

  return { dependencies, executionBatches };
}

export function getStepDependencies(
  stepKey: string,
  graph: DependencyGraph
): string[] {
  return graph.dependencies.get(stepKey) || [];
}