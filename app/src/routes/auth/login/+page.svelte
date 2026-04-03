<script lang="ts">
  import { enhance } from "$app/forms";

  let { form } = $props();
</script>

<svelte:head>
  <title>ログイン - Prometheus</title>
</svelte:head>

<div class="flex min-h-[80vh] items-center justify-center">
  <div class="w-full max-w-sm">
    <h1 class="mb-2 text-center font-mono text-2xl font-bold text-text-main">Prometheus</h1>
    <p class="mb-8 text-center text-sm text-text-muted">アカウントにログイン</p>

    {#if form?.error}
      <div class="mb-4 rounded-lg border border-error/30 bg-error/10 px-4 py-3 text-sm text-error">
        {form.error}
      </div>
    {/if}

    <form method="POST" use:enhance class="space-y-4">
      <div>
        <label for="email" class="mb-1.5 block text-xs font-medium text-text-muted">メールアドレス</label>
        <input
          id="email"
          name="email"
          type="email"
          required
          class="w-full rounded-lg border border-border bg-bg-card px-3.5 py-2.5 text-sm text-text-main outline-none placeholder:text-text-dim focus:border-primary"
          placeholder="you@example.com"
        />
      </div>

      <div>
        <label for="password" class="mb-1.5 block text-xs font-medium text-text-muted">パスワード</label>
        <input
          id="password"
          name="password"
          type="password"
          required
          class="w-full rounded-lg border border-border bg-bg-card px-3.5 py-2.5 text-sm text-text-main outline-none placeholder:text-text-dim focus:border-primary"
          placeholder="パスワード"
        />
      </div>

      <button
        type="submit"
        class="w-full rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-primary-dark"
      >
        ログイン
      </button>
    </form>

    <p class="mt-6 text-center text-sm text-text-muted">
      アカウントをお持ちでない方は
      <a href="/auth/register" class="text-primary-light hover:underline">新規登録</a>
    </p>
  </div>
</div>
