<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { Search, Filter, Eye, EyeOff } from "lucide-svelte";
  import {
    forceSimulation,
    forceLink,
    forceManyBody,
    forceCenter,
    forceCollide,
    forceX,
    forceY,
    type SimulationNodeDatum,
    type SimulationLinkDatum,
  } from "d3-force";
  import { select } from "d3-selection";
  import { zoom, zoomIdentity } from "d3-zoom";
  import { drag } from "d3-drag";

  interface GraphNode extends SimulationNodeDatum {
    id: string;
    path: string;
    title: string;
    tags: string[];
  }

  interface GraphLink extends SimulationLinkDatum<GraphNode> {
    slug: string;
  }

  let container: HTMLDivElement;
  let svgEl: SVGSVGElement | undefined = $state(undefined);
  let allNodes: GraphNode[] = $state([]);
  let allLinks: GraphLink[] = $state([]);
  let loading = $state(true);

  let searchQuery = $state("");
  let selectedTags = $state<Set<string>>(new Set());
  let allTags = $state<string[]>([]);
  let showFilters = $state(false);
  let showHulls = $state(true);
  let showTagLabels = $state(false);

  const TAG_COLORS = [
    { fill: "rgba(99, 102, 241, 0.07)", stroke: "rgba(99, 102, 241, 0.25)" },
    { fill: "rgba(14, 165, 233, 0.07)", stroke: "rgba(14, 165, 233, 0.25)" },
    { fill: "rgba(34, 197, 94, 0.07)", stroke: "rgba(34, 197, 94, 0.25)" },
    { fill: "rgba(236, 72, 153, 0.07)", stroke: "rgba(236, 72, 153, 0.25)" },
    { fill: "rgba(245, 158, 11, 0.07)", stroke: "rgba(245, 158, 11, 0.25)" },
    { fill: "rgba(168, 85, 247, 0.07)", stroke: "rgba(168, 85, 247, 0.25)" },
    { fill: "rgba(6, 182, 212, 0.07)", stroke: "rgba(6, 182, 212, 0.25)" },
    { fill: "rgba(239, 68, 68, 0.07)", stroke: "rgba(239, 68, 68, 0.25)" },
  ];

  // Stable tag-to-color index
  let tagColorIndex = $state<Map<string, number>>(new Map());

  function getFilteredData(): { nodes: GraphNode[]; links: GraphLink[] } {
    let nodes = [...allNodes];
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      nodes = nodes.filter((n) => n.title.toLowerCase().includes(q) || n.tags.some((t) => t.includes(q)));
    }
    if (selectedTags.size > 0) {
      nodes = nodes.filter((n) => n.tags.some((t) => selectedTags.has(t)));
    }
    const nodeIds = new Set(nodes.map((n) => n.id));
    const links = allLinks.filter((l) => {
      const src = typeof l.source === "string" ? l.source : (l.source as GraphNode).id;
      const tgt = typeof l.target === "string" ? l.target : (l.target as GraphNode).id;
      return nodeIds.has(src) && nodeIds.has(tgt);
    });
    return { nodes, links };
  }

  let currentSim: ReturnType<typeof forceSimulation> | null = null;

  onMount(async () => {
    await loadGraph();
    loading = false;
    if (allNodes.length > 0) {
      requestAnimationFrame(() => rebuildGraph());
    }
  });

  async function loadGraph() {
    try {
      const res = await fetch("/api/graph");
      const data = await res.json();
      allNodes = data.nodes;
      allLinks = data.links;

      const tags = new Set<string>();
      for (const n of allNodes) for (const t of n.tags) tags.add(t);
      allTags = Array.from(tags).sort();

      const idx = new Map<string, number>();
      let i = 0;
      for (const t of allTags) { idx.set(t, i++); }
      tagColorIndex = idx;
    } catch {
      allNodes = [];
      allLinks = [];
    }
  }

  function rebuildGraph() {
    if (!svgEl || !container || allNodes.length === 0) return;
    if (currentSim) currentSim.stop();
    const { nodes, links } = getFilteredData();
    for (const n of nodes) { n.x = undefined; n.y = undefined; n.vx = undefined; n.vy = undefined; }
    renderGraph(nodes, links);
  }

  // Debounced rebuild for search input
  let searchTimer: ReturnType<typeof setTimeout> | undefined;
  function onSearchInput() {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => rebuildGraph(), 200);
  }

  function toggleTag(tag: string) {
    const next = new Set(selectedTags);
    if (next.has(tag)) next.delete(tag);
    else next.add(tag);
    selectedTags = next;
    rebuildGraph();
  }

  function computeHull(points: { x: number; y: number }[], padding: number): string {
    if (points.length === 0) return "";

    const cx = points.reduce((s, p) => s + p.x, 0) / points.length;
    const cy = points.reduce((s, p) => s + p.y, 0) / points.length;

    if (points.length === 1) {
      return `M ${cx - padding} ${cy} A ${padding} ${padding} 0 1 1 ${cx + padding} ${cy} A ${padding} ${padding} 0 1 1 ${cx - padding} ${cy}`;
    }

    const sorted = [...points].sort((a, b) => Math.atan2(a.y - cy, a.x - cx) - Math.atan2(b.y - cy, b.x - cx));

    const expanded = sorted.map((p) => {
      const dx = p.x - cx;
      const dy = p.y - cy;
      const dist = Math.sqrt(dx * dx + dy * dy);
      if (dist < 1) return { x: cx + padding, y: cy };
      return { x: cx + (dx / dist) * (dist + padding), y: cy + (dy / dist) * (dist + padding) };
    });

    if (expanded.length === 2) {
      const [a, b] = expanded;
      const dx = b.x - a.x;
      const dy = b.y - a.y;
      const nx = -dy * 0.35;
      const ny = dx * 0.35;
      const mx = (a.x + b.x) / 2;
      const my = (a.y + b.y) / 2;
      return `M ${a.x} ${a.y} Q ${mx + nx} ${my + ny} ${b.x} ${b.y} Q ${mx - nx} ${my - ny} ${a.x} ${a.y}`;
    }

    let path = "";
    for (let i = 0; i < expanded.length; i++) {
      const curr = expanded[i];
      const next = expanded[(i + 1) % expanded.length];
      const mx = (curr.x + next.x) / 2;
      const my = (curr.y + next.y) / 2;
      if (i === 0) path = `M ${mx} ${my}`;
      const nn = expanded[(i + 2) % expanded.length];
      path += ` Q ${next.x} ${next.y} ${(next.x + nn.x) / 2} ${(next.y + nn.y) / 2}`;
    }
    return path + " Z";
  }

  function renderGraph(nodes: GraphNode[], links: GraphLink[]) {
    if (!svgEl || !container) return;

    const width = container.clientWidth;
    const height = container.clientHeight;

    const svg = select(svgEl).attr("width", width).attr("height", height);
    svg.selectAll("*").remove();
    const g = svg.append("g");

    const initialTransform = zoomIdentity.translate(width / 2, height / 2);

    const zoomBehavior = zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.1, 4])
      .on("zoom", (event) => g.attr("transform", event.transform));

    svg.call(zoomBehavior);
    svg.call(zoomBehavior.transform, initialTransform);

    // Apply transform immediately so first frame is centered
    g.attr("transform", initialTransform.toString());

    const linkCounts = new Map<string, number>();
    for (const l of links) {
      const s = typeof l.source === "string" ? l.source : (l.source as GraphNode).id;
      const t = typeof l.target === "string" ? l.target : (l.target as GraphNode).id;
      linkCounts.set(s, (linkCounts.get(s) ?? 0) + 1);
      linkCounts.set(t, (linkCounts.get(t) ?? 0) + 1);
    }

    const tagMap = new Map<string, GraphNode[]>();
    for (const node of nodes) {
      for (const tag of node.tags) {
        if (!tagMap.has(tag)) tagMap.set(tag, []);
        tagMap.get(tag)!.push(node);
      }
    }
    const tagEntries = Array.from(tagMap.entries());

    // Spread out based on node count
    const spread = Math.max(100, nodes.length * 15);

    const simulation = forceSimulation<GraphNode>(nodes)
      .force("link", forceLink<GraphNode, GraphLink>(links).id((d) => d.id).distance(80).strength(0.5))
      .force("charge", forceManyBody().strength(-150))
      .force("center", forceCenter(0, 0).strength(0.04))
      .force("collide", forceCollide(28))
      .force("x", forceX(0).strength(0.02))
      .force("y", forceY(0).strength(0.02));

    currentSim = simulation;

    // Hulls
    const hullGroup = g.append("g");
    let hulls: any;
    let tagLabelsEl: any;

    if (showHulls) {
      hulls = hullGroup.selectAll("path").data(tagEntries).join("path")
        .attr("fill", ([tag]) => TAG_COLORS[(tagColorIndex.get(tag) ?? 0) % TAG_COLORS.length].fill)
        .attr("stroke", ([tag]) => TAG_COLORS[(tagColorIndex.get(tag) ?? 0) % TAG_COLORS.length].stroke)
        .attr("stroke-width", 1)
        .attr("stroke-dasharray", "4 3");
    }

    if (showTagLabels) {
      tagLabelsEl = hullGroup.selectAll("text").data(tagEntries).join("text")
        .attr("fill", "var(--color-text-dim)")
        .attr("font-size", "9px")
        .attr("font-family", "var(--font-sans)")
        .attr("text-anchor", "middle")
        .attr("pointer-events", "none")
        .attr("opacity", 0.6)
        .text(([tag]) => `#${tag}`);
    }

    // Links
    const linkEl = g.append("g").selectAll("line").data(links).join("line")
      .attr("stroke", "var(--color-text-dim)")
      .attr("stroke-width", 1.5)
      .attr("stroke-opacity", 0.3);

    // Nodes
    const nodeEl = g.append("g").selectAll<SVGGElement, GraphNode>("g").data(nodes).join("g")
      .attr("cursor", "pointer")
      .call(
        drag<SVGGElement, GraphNode>()
          .on("start", (event, d) => { if (!event.active) simulation.alphaTarget(0.3).restart(); d.fx = d.x; d.fy = d.y; })
          .on("drag", (event, d) => { d.fx = event.x; d.fy = event.y; })
          .on("end", (event, d) => { if (!event.active) simulation.alphaTarget(0); d.fx = null; d.fy = null; }),
      );

    nodeEl.append("circle")
      .attr("r", (d) => (linkCounts.get(d.id) ?? 0) > 0 ? Math.max(5, Math.min(10, 5 + (linkCounts.get(d.id) ?? 0) * 1.2)) : 4)
      .attr("fill", (d) => (linkCounts.get(d.id) ?? 0) > 0 ? "var(--color-primary)" : "var(--color-text-dim)")
      .attr("stroke", "var(--color-bg-dark)")
      .attr("stroke-width", 1.5);

    nodeEl.append("text")
      .text((d) => d.title)
      .attr("x", 0)
      .attr("y", (d) => -((linkCounts.get(d.id) ?? 0) > 0 ? Math.max(5, 5 + (linkCounts.get(d.id) ?? 0) * 1.2) : 4) - 5)
      .attr("text-anchor", "middle")
      .attr("fill", "var(--color-text-muted)")
      .attr("font-size", "10px")
      .attr("font-family", "var(--font-sans)")
      .attr("pointer-events", "none");

    nodeEl
      .on("mouseenter", function (_, d) {
        select(this).select("circle").attr("fill", "var(--color-primary-light)").attr("r", (linkCounts.get(d.id) ?? 0) > 0 ? Math.max(5, 5 + (linkCounts.get(d.id) ?? 0) * 1.2) + 2 : 6);
        select(this).select("text").attr("fill", "var(--color-text-main)").attr("font-weight", "600");
        linkEl.attr("stroke-opacity", (l) => (l.source as GraphNode).id === d.id || (l.target as GraphNode).id === d.id ? 0.8 : 0.05);
      })
      .on("mouseleave", function (_, d) {
        select(this).select("circle").attr("fill", (linkCounts.get(d.id) ?? 0) > 0 ? "var(--color-primary)" : "var(--color-text-dim)").attr("r", (linkCounts.get(d.id) ?? 0) > 0 ? Math.max(5, 5 + (linkCounts.get(d.id) ?? 0) * 1.2) : 4);
        select(this).select("text").attr("fill", "var(--color-text-muted)").attr("font-weight", "normal");
        linkEl.attr("stroke-opacity", 0.3);
      })
      .on("click", (_, d) => goto(`/note/${d.path}`));

    simulation.on("tick", () => {
      linkEl
        .attr("x1", (d) => (d.source as GraphNode).x ?? 0).attr("y1", (d) => (d.source as GraphNode).y ?? 0)
        .attr("x2", (d) => (d.target as GraphNode).x ?? 0).attr("y2", (d) => (d.target as GraphNode).y ?? 0);
      nodeEl.attr("transform", (d) => `translate(${d.x ?? 0},${d.y ?? 0})`);

      if (showHulls && hulls) {
        hulls.attr("d", ([, gn]: [string, GraphNode[]]) => {
          const pts = gn.filter((n) => n.x != null).map((n) => ({ x: n.x!, y: n.y! }));
          return computeHull(pts, 40);
        });
      }
      if (showTagLabels && tagLabelsEl) {
        tagLabelsEl
          .attr("x", ([, gn]: [string, GraphNode[]]) => { const ps = gn.filter((n) => n.x != null); return ps.length > 0 ? ps.reduce((s, n) => s + n.x!, 0) / ps.length : 0; })
          .attr("y", ([, gn]: [string, GraphNode[]]) => { const ps = gn.filter((n) => n.y != null); return ps.length > 0 ? Math.min(...ps.map((n) => n.y!)) - 48 : 0; });
      }
    });
  }
</script>

<div bind:this={container} class="relative h-full w-full bg-bg-dark">
  {#if loading}
    <div class="flex h-full items-center justify-center">
      <span class="text-sm text-text-muted">グラフを読み込み中...</span>
    </div>
  {:else if allNodes.length === 0}
    <div class="flex h-full items-center justify-center">
      <div class="text-center">
        <p class="text-text-muted">ノートがありません</p>
        <p class="mt-1.5 text-sm text-text-dim">ノートを作成するとグラフに表示されます</p>
      </div>
    </div>
  {:else}
    <svg bind:this={svgEl} class="h-full w-full"></svg>

    <!-- Controls -->
    <div class="absolute right-3 top-3 flex flex-col gap-1.5">
      <div class="flex items-center gap-1.5 rounded-lg border border-border bg-bg-card/90 px-2.5 py-1.5 backdrop-blur-sm">
        <Search size={12} class="text-text-dim" />
        <input
          bind:value={searchQuery}
          oninput={onSearchInput}
          placeholder="検索..."
          class="w-28 bg-transparent text-xs text-text-main outline-none placeholder:text-text-dim"
        />
      </div>

      <button
        onclick={() => (showFilters = !showFilters)}
        class="flex items-center gap-1.5 rounded-lg border border-border bg-bg-card/90 px-2.5 py-1.5 text-xs text-text-muted backdrop-blur-sm hover:text-text-main"
      >
        <Filter size={12} />
        タグ
        {#if selectedTags.size > 0}
          <span class="rounded-full bg-primary px-1.5 text-[10px] text-white">{selectedTags.size}</span>
        {/if}
      </button>

      {#if showFilters}
        <div class="max-h-48 overflow-y-auto rounded-lg border border-border bg-bg-card/95 p-1.5 backdrop-blur-sm">
          {#if selectedTags.size > 0}
            <button onclick={() => { selectedTags = new Set(); rebuildGraph(); }} class="w-full rounded px-2 py-1 text-left text-[10px] text-text-dim hover:bg-bg-card-hover">
              リセット
            </button>
          {/if}
          {#each allTags as tag}
            <button
              onclick={() => toggleTag(tag)}
              class="flex w-full items-center gap-2 rounded px-2 py-1 text-left text-[11px] hover:bg-bg-card-hover {selectedTags.has(tag) ? 'text-text-main' : 'text-text-muted'}"
            >
              <span class="h-2 w-2 rounded-sm" style="background-color: {TAG_COLORS[(tagColorIndex.get(tag) ?? 0) % TAG_COLORS.length].stroke}"></span>
              {tag}
              {#if selectedTags.has(tag)}<span class="ml-auto text-primary text-[10px]">✓</span>{/if}
            </button>
          {/each}
        </div>
      {/if}

      <!-- Display toggles -->
      <div class="flex gap-1">
        <button
          onclick={() => { showHulls = !showHulls; rebuildGraph(); }}
          class="flex items-center gap-1 rounded-lg border border-border bg-bg-card/90 px-2 py-1 text-[10px] backdrop-blur-sm {showHulls ? 'text-primary-light' : 'text-text-dim'}"
          title="集合の表示"
        >
          {#if showHulls}<Eye size={10} />{:else}<EyeOff size={10} />{/if}
          集合
        </button>
        <button
          onclick={() => { showTagLabels = !showTagLabels; rebuildGraph(); }}
          class="flex items-center gap-1 rounded-lg border border-border bg-bg-card/90 px-2 py-1 text-[10px] backdrop-blur-sm {showTagLabels ? 'text-primary-light' : 'text-text-dim'}"
          title="タグ名の表示"
        >
          {#if showTagLabels}<Eye size={10} />{:else}<EyeOff size={10} />{/if}
          ラベル
        </button>
      </div>
    </div>

    <div class="absolute bottom-3 right-3 text-[10px] text-text-dim">
      {getFilteredData().nodes.length} ノート / {getFilteredData().links.length} リンク
    </div>
  {/if}
</div>
