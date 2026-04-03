<script lang="ts">
  import { CheckSquare, Circle, CheckCircle2, ArrowRight } from "lucide-svelte";

  interface Task {
    id: number;
    noteId: string;
    content: string;
    completed: boolean;
    lineNumber: number;
    dueDate: string | null;
    notePath: string;
    noteTitle: string;
  }

  let tasks: Task[] = $state([]);
  let filter = $state<"all" | "pending" | "completed">("pending");
  let loading = $state(true);

  $effect(() => {
    loadTasks();
  });

  async function loadTasks() {
    loading = true;
    try {
      const res = await fetch(`/api/tasks?filter=${filter}`);
      tasks = await res.json();
    } catch {
      tasks = [];
    } finally {
      loading = false;
    }
  }

  let groupedTasks = $derived(() => {
    const groups = new Map<string, { title: string; path: string; tasks: Task[] }>();
    for (const task of tasks) {
      if (!groups.has(task.notePath)) {
        groups.set(task.notePath, {
          title: task.noteTitle,
          path: task.notePath,
          tasks: [],
        });
      }
      groups.get(task.notePath)!.tasks.push(task);
    }
    return Array.from(groups.values());
  });

  async function toggleTask(task: Task) {
    const newCompleted = !task.completed;
    // Optimistic update
    task.completed = newCompleted;
    tasks = [...tasks];

    try {
      await fetch(`/api/tasks/${task.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ completed: newCompleted }),
      });
    } catch {
      // Revert on error
      task.completed = !newCompleted;
      tasks = [...tasks];
    }
  }
</script>

<svelte:head>
  <title>タスク - Prometheus</title>
</svelte:head>

<div class="mx-auto max-w-2xl py-4">
  <div class="mb-5 flex items-center justify-between">
    <h1 class="flex items-center gap-2 font-mono text-2xl font-bold text-text-main">
      <CheckSquare size={22} />
      タスク
    </h1>
    <div class="flex rounded-lg border border-border text-xs">
      {#each [
        { value: "pending", label: "未完了" },
        { value: "all", label: "すべて" },
        { value: "completed", label: "完了" },
      ] as tab}
        <button
          onclick={() => (filter = tab.value as typeof filter)}
          class="px-3 py-1.5 transition-colors first:rounded-l-lg last:rounded-r-lg {filter === tab.value
            ? 'bg-primary text-white'
            : 'text-text-muted hover:text-text-main'}"
        >
          {tab.label}
        </button>
      {/each}
    </div>
  </div>

  {#if loading}
    <p class="text-sm text-text-muted">読み込み中...</p>
  {:else if tasks.length === 0}
    <div class="rounded-xl border border-border bg-bg-card p-8 text-center">
      <p class="text-sm text-text-muted">
        {filter === "pending" ? "未完了のタスクはありません" : filter === "completed" ? "完了したタスクはありません" : "タスクがありません"}
      </p>
      <p class="mt-1.5 text-xs text-text-dim">
        ノート内で <code class="rounded bg-bg-dark px-1.5 py-0.5 text-accent">- [ ]</code> と書くとタスクとして認識されます
      </p>
    </div>
  {:else}
    <div class="space-y-5">
      {#each groupedTasks() as group}
        <div>
          <a
            href="/note/{group.path}"
            class="mb-2 inline-flex items-center gap-1.5 text-sm font-medium text-primary-light hover:underline"
          >
            {group.title}
            <ArrowRight size={12} />
          </a>
          <div class="space-y-1">
            {#each group.tasks as task}
              <div
                class="flex items-start gap-3 rounded-lg border border-border bg-bg-card px-4 py-2.5 transition-colors hover:bg-bg-card-hover"
              >
                <button
                  onclick={() => toggleTask(task)}
                  class="mt-0.5 shrink-0 transition-colors hover:text-primary"
                  title={task.completed ? "未完了に戻す" : "完了にする"}
                >
                  {#if task.completed}
                    <CheckCircle2 size={15} class="text-success" />
                  {:else}
                    <Circle size={15} class="text-text-dim" />
                  {/if}
                </button>
                <span class="flex-1 text-sm {task.completed ? 'text-text-dim line-through' : 'text-text-main'}">
                  {task.content}
                </span>
                {#if task.dueDate}
                  <span class="shrink-0 text-xs text-warning">{task.dueDate}</span>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/each}
    </div>

    <p class="mt-4 text-xs text-text-dim">
      {tasks.length}件のタスク
    </p>
  {/if}
</div>
