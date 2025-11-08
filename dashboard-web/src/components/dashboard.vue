<template>
  <div class="dashboard">
    <header class="header">
      <div class="title-wrap">
        <h1>🏆 Linux Exam 排行榜</h1>
      </div>
      <div class="controls">
        <div class="stats">
          <div class="stat-card">
            <div class="stat-value">{{ totalPlayers }}</div>
            <div class="stat-label">参赛人数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ topScore }}</div>
            <div class="stat-label">最高分</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ lastUpdatedText }}</div>
            <div class="stat-label">最近更新时间</div>
          </div>
        </div>

        <button class="btn-refresh" @click="fetchData" title="手动刷新">刷新</button>
      </div>
    </header>

    <main class="main-content">
      <div class="chart-container">
        <div ref="chart" class="chart"></div>
      </div>
    </main>

    <footer class="footer">
      <p>Powered by <strong>Vue 3 + ECharts</strong></p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from "vue";
import * as echarts from "echarts";

type Player = { username: string; score: number };

const chart = ref<HTMLElement | null>(null);
let chartInstance: echarts.ECharts | null = null;
const players = ref<Player[]>([]);
const lastUpdated = ref<Date | null>(null);
let timer: number | null = null;

const totalPlayers = computed(() => players.value.length);
const topScore = computed(() => {
  if (players.value.length === 0) return "--";
  return Math.max(...players.value.map((p) => p.score));
});
const lastUpdatedText = computed(() =>
  lastUpdated.value ? lastUpdated.value.toLocaleString() : "--"
);

// 拉取数据
async function fetchData() {
  try {
    const resp = await fetch("/api/dashboard");
    const json = await resp.json();

    if (json.code !== 0) {
      console.error("获取排行榜失败:", json.message);
      return;
    }

  players.value = (json.data || []) as Player[];
    lastUpdated.value = new Date();
    // ensure chart renders after data updated
    renderChart();
  } catch (err) {
    console.error("请求出错:", err);
  }
}

// 渲染图表
function renderChart() {
  if (!chartInstance || players.value.length === 0) return;
  const sorted = [...players.value].sort((a, b) => b.score - a.score);
  // 备用颜色（用于第4名及以后，比较素雅）
  const mutedColor = "#DCDCDC";
  // 其他非奖牌位使用统一的素色
  // 更鲜艳的金/银/铜色
  const podiumColors = ["#ffb300", "#9aa6b2", "#cd7f32"]; // 🥇 🥈 🥉
  const medalIcons = ["🥇", "🥈", "🥉"];

  // 固定每条 bar 的高度（px），容器高度根据条目数增长
  const BAR_HEIGHT = 48; // 单个柱子高度
  const PADDING = 120; // 上下内边距 + label space
  const neededHeight = Math.max(420, sorted.length * BAR_HEIGHT + PADDING);

  if (chart.value) {
    (chart.value as HTMLElement).style.height = `${neededHeight}px`;
  }

  const option: echarts.EChartsOption = {
    backgroundColor: "transparent",
    // 初始与更新动画设置
    animationDuration: 600,
    animationEasing: 'cubicOut',
    animationDurationUpdate: 800,
    animationEasingUpdate: 'cubicOut',
    grid: { left: "10%", right: "10%", top: 40, bottom: 30 },
    xAxis: {
      type: "value",
      max: Math.max(...sorted.map((p) => p.score)) + 50,
      axisLabel: { color: "#444", fontSize: 14 },
      splitLine: { show: false },
    },
    yAxis: {
      type: "category",
      data: sorted.map(p => p.username),
      inverse: true,
      axisLabel: {
        color: "#222",
        margin: 12,
        rich: {
          rank: {
            width: 40,
            align: 'right',
            fontSize: 16,
            color: '#666',
            padding: [0, 8, 0, 0]
          },
          medal: {
            fontSize: 34,
            width: 40,
            align: 'center',
            padding: [0, 8, 0, 0]
          },
          name: {
            width: 100,
            align: 'left',
            fontSize: 16,
            color: '#2c3e50',
            fontWeight: 500
          }
        },
        formatter: function (value: string, idx: number) {
          if (idx < 3) {
            // 前三名：奖牌 + 用户名
            return `{medal|${medalIcons[idx]}} {name|${value}}`;
          }
          // 其他：数字序号 + 用户名
          return `{rank|${idx + 1}.} {name|${value}}`;
        }
      },
      axisTick: { show: false },
      axisLine: { show: false },
    },
    series: [
      {
        type: "bar",
        data: sorted.map((p, i) => ({
          value: p.score,
          itemStyle: {
            // 前三名使用鲜艳色并加徽章，其他名次使用素色
            color: i < 3 ? podiumColors[i] : mutedColor,
            // 四个角都圆角
            borderRadius: 8,
          },
        })),
        label: {
          show: true,
          position: "right",
          formatter: "{c} 分",
          color: "#333",
          fontWeight: "bold",
          fontSize: 16,
        },
        barWidth: 36, // 每个柱子固定高度（像素）
        barGap: '10%'
      },
    ],
  };

  // 等待浏览器应用高度变更后再 resize，确保 canvas 尺寸正确
  if (chart.value && chartInstance) {
    requestAnimationFrame(() => requestAnimationFrame(() => {
      chartInstance!.resize();
      chartInstance!.setOption(option as any, { notMerge: false, lazyUpdate: false });
    }));
  } else if (chartInstance) {
    chartInstance.setOption(option as any, { notMerge: false, lazyUpdate: false });
  }
}

onMounted(() => {
  chartInstance = echarts.init(chart.value);
  fetchData();

  // 自动轮询
  timer = setInterval(fetchData, 5000);

  window.addEventListener("resize", () => { if (chartInstance) chartInstance.resize(); });
});

onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
});
</script>

<style scoped>
.dashboard {
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  background-color: #f4f6f8; /* 纯色背景 */
  color: #333;
  text-align: center;
  font-family: "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 40px 0;
  box-sizing: border-box;
}

.header {
  margin-bottom: 20px;
}

.header h1 {
  font-size: 36px;
  margin-bottom: 6px;
  color: #222;
}

.subtitle {
  font-size: 16px;
  color: #666;
}

/* （已移除旧样式以避免与新的可伸缩样式冲突） */

.footer {
  margin-top: 20px;
  color: #777;
  font-size: 14px;
}

/* New styles for improved UI */
.header {
  width: 100%;
  max-width: 1200px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 18px;
}

.title-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-start;
}

.controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.stats {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.stat-card {
  background: linear-gradient(180deg, rgba(255,255,255,0.98), #fff);
  padding: 12px 16px;
  border-radius: 10px;
  min-width: 120px;
  box-shadow: 0 6px 18px rgba(20, 30, 60, 0.06);
  text-align: left;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #111827;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.btn-refresh {
  background: linear-gradient(90deg,#6366f1,#4f46e5);
  color: #fff;
  border: none;
  padding: 10px 14px;
  border-radius: 10px;
  cursor: pointer;
  box-shadow: 0 6px 16px rgba(79,70,229,0.18);
  transition: transform .12s ease, box-shadow .12s ease;
}
.btn-refresh:hover { transform: translateY(-2px); }

.main-content {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-top: 18px;
}

.chart-container {
  width: 90%;
  max-width: 1200px;
  background: linear-gradient(180deg,#ffffff, #fbfdff);
  border-radius: 14px;
  box-shadow: 0 8px 30px rgba(16,24,40,0.06);
  padding: 22px;
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
}

.chart {
  width: 100%;
  flex: 1 1 auto;
  min-height: 420px;
}

@media (max-width: 720px) {
  .controls { flex-direction: column; align-items: stretch; }
  .btn-refresh { align-self: flex-end; }
  .chart { height: 380px; }
  .stat-card { min-width: 100px; }
}
</style>
