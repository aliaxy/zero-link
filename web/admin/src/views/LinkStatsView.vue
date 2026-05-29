<template>
  <div>
    <div class="back-row">
      <el-button text :icon="ArrowLeft" @click="router.push('/links')">Links</el-button>
    </div>

    <div class="page-header">
      <div>
        <h1 class="page-title">
          <span class="mono">{{ linkCode }}</span> — Stats
        </h1>
        <p v-if="linkDescription" class="page-subtitle">{{ linkDescription }}</p>
        <p v-else class="page-subtitle">Daily page views and unique visitors</p>
      </div>

      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="to"
        start-placeholder="From"
        end-placeholder="To"
        format="YYYY-MM-DD"
        value-format="YYYY-MM-DD"
        :disabled-date="disabledDate"
        @change="fetchStats"
      />
    </div>

    <!-- Summary cards -->
    <div class="stat-cards">
      <div class="stat-card">
        <span class="stat-label">Total Views</span>
        <span class="stat-value">{{ totalPV.toLocaleString() }}</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">Unique Visitors</span>
        <span class="stat-value">{{ totalUV.toLocaleString() }}</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">Peak Day</span>
        <span class="stat-value">{{ peakDay }}</span>
        <span class="stat-meta">{{ peakPV.toLocaleString() }} views</span>
      </div>
    </div>

    <!-- Chart -->
    <el-card class="chart-card" v-loading="loading">
      <StatsChart :items="items" />
    </el-card>

    <!-- Data table -->
    <el-card class="data-card" v-if="items.length">
      <el-table :data="items" style="width: 100%" size="small">
        <el-table-column prop="stat_date" label="Date" width="140" />
        <el-table-column prop="pv" label="Page Views" />
        <el-table-column prop="uv" label="Unique Visitors" />
      </el-table>
    </el-card>

    <div v-else-if="!loading" class="empty-state">
      <p>No visit data for this period.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getLink, getLinkStats } from '@/api/links'
import StatsChart from '@/components/StatsChart.vue'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const linkCode = ref('')
const linkDescription = ref('')
const items = ref([])
const loading = ref(false)

const defaultFrom = dayjs().subtract(29, 'day').format('YYYY-MM-DD')
const defaultTo = dayjs().format('YYYY-MM-DD')
const dateRange = ref([defaultFrom, defaultTo])

const totalPV = computed(() => items.value.reduce((s, i) => s + i.pv, 0))
const totalUV = computed(() => items.value.reduce((s, i) => s + i.uv, 0))

const peakItem = computed(() =>
  items.value.reduce((max, i) => (i.pv > (max?.pv ?? -1) ? i : max), null),
)
const peakDay = computed(() => peakItem.value?.stat_date ?? '—')
const peakPV = computed(() => peakItem.value?.pv ?? 0)

function disabledDate(d) {
  const from = dateRange.value?.[0]
  if (!from) return false
  return dayjs(d).diff(dayjs(from), 'day') > 90
}

async function fetchStats() {
  loading.value = true
  try {
    const [from, to] = dateRange.value ?? [defaultFrom, defaultTo]
    const data = await getLinkStats(route.params.id, { from, to })
    items.value = data.items ?? []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const link = await getLink(route.params.id).catch(() => null)
  linkCode.value = link?.code ?? route.params.id
  linkDescription.value = link?.description ?? ''
  await fetchStats()
})
</script>

<style scoped>
.back-row {
  margin-bottom: 20px;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
  gap: 16px;
  flex-wrap: wrap;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-label);
  letter-spacing: -0.03em;
}

.page-title .mono {
  color: var(--color-accent);
  font-size: 22px;
}

.page-subtitle {
  font-size: 13px;
  color: var(--color-label-tertiary);
  margin-top: 4px;
}

.stat-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-card);
  border: 1px solid var(--color-separator);
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-label-tertiary);
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-label);
  letter-spacing: -0.03em;
  line-height: 1.1;
}

.stat-meta {
  font-size: 12px;
  color: var(--color-label-secondary);
}

.chart-card {
  margin-bottom: 16px;
}

.chart-card :deep(.el-card__body) {
  padding: 20px;
}

.data-card :deep(.el-card__body) {
  padding: 0;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--color-label-secondary);
  font-size: 14px;
}

@media (max-width: 700px) {
  .stat-cards {
    grid-template-columns: 1fr;
  }
}
</style>
