<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">Links</h1>
        <p class="page-subtitle">{{ total }} short links</p>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">New Link</el-button>
    </div>

    <!-- Filters -->
    <el-card class="filter-card">
      <div class="filter-row">
        <el-input
          v-model="keyword"
          placeholder="Search by code, title, or URL…"
          :prefix-icon="Search"
          clearable
          style="width: 300px"
          @input="debouncedFetch"
          @clear="fetchLinks"
        />
        <el-select
          v-model="statusFilter"
          placeholder="All statuses"
          clearable
          style="width: 150px"
          @change="fetchLinks"
        >
          <el-option label="Active" :value="1" />
          <el-option label="Disabled" :value="0" />
        </el-select>
      </div>
    </el-card>

    <!-- Table -->
    <el-card class="table-card">
      <el-table
        v-loading="loading"
        :data="links"
        row-key="id"
        style="width: 100%"
        class="clickable-table"
        @row-click="goDetail"
      >
        <el-table-column label="Code" width="140">
          <template #default="{ row }">
            <span class="code-cell mono">{{ row.code }}</span>
          </template>
        </el-table-column>

        <el-table-column label="Destination URL" min-width="220">
          <template #default="{ row }">
            <el-tooltip :content="row.origin_url" placement="top" :show-after="400">
              <a
                :href="row.origin_url"
                target="_blank"
                rel="noopener"
                class="url-cell truncate"
                @click.stop
              >{{ row.origin_url }}</a>
            </el-tooltip>
          </template>
        </el-table-column>

        <el-table-column label="Title" min-width="140">
          <template #default="{ row }">
            <span class="text-secondary">{{ row.title || '—' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="Status" width="100">
          <template #default="{ row }">
            <LinkStatusTag :status="row.status" :expire-at="row.expire_at" />
          </template>
        </el-table-column>

        <el-table-column label="Expires" width="120">
          <template #default="{ row }">
            <span class="text-secondary">{{ row.expire_at ? formatDate(row.expire_at) : '—' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="Created" width="110">
          <template #default="{ row }">
            <span class="text-secondary">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="Actions" width="140" align="right">
          <template #default="{ row }">
            <div class="row-actions" @click.stop>
              <el-tooltip content="Edit" placement="top">
                <button class="action-btn" @click="openEdit(row)">
                  <el-icon :size="16"><Edit /></el-icon>
                </button>
              </el-tooltip>
              <el-tooltip content="Stats" placement="top">
                <button class="action-btn" @click="goStats(row)">
                  <el-icon :size="16"><DataLine /></el-icon>
                </button>
              </el-tooltip>
              <el-popconfirm
                title="Delete this link?"
                confirm-button-text="Delete"
                confirm-button-type="danger"
                @confirm="handleDelete(row)"
              >
                <template #reference>
                  <button class="action-btn action-btn--danger">
                    <el-icon :size="16"><Delete /></el-icon>
                  </button>
                </template>
              </el-popconfirm>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > pageSize"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @change="fetchLinks"
      />
    </el-card>

    <LinkFormDrawer
      v-model="drawerVisible"
      :link="editingLink"
      @saved="fetchLinks"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Search, Edit, Delete, DataLine } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { listLinks, deleteLink } from '@/api/links'
import LinkStatusTag from '@/components/LinkStatusTag.vue'
import LinkFormDrawer from '@/components/LinkFormDrawer.vue'
import dayjs from 'dayjs'

const router = useRouter()

const links = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const statusFilter = ref(null)
const loading = ref(false)
const drawerVisible = ref(false)
const editingLink = ref(null)

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchLinks, 350)
}

async function fetchLinks() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (keyword.value) params.keyword = keyword.value
    if (statusFilter.value !== null) params.status = statusFilter.value
    const data = await listLinks(params)
    links.value = data.items ?? []
    total.value = data.total ?? 0
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingLink.value = null
  drawerVisible.value = true
}

function openEdit(row) {
  editingLink.value = row
  drawerVisible.value = true
}

function goDetail(row) {
  router.push({ name: 'link-detail', params: { id: row.id } })
}

function goStats(row) {
  router.push({ name: 'link-stats', params: { id: row.id } })
}

async function handleDelete(row) {
  await deleteLink(row.id)
  ElMessage.success('Link deleted')
  fetchLinks()
}

function formatDate(iso) {
  return dayjs(iso).format('MMM D, YYYY HH:mm')
}

onMounted(fetchLinks)
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-title {
  font-size: 26px;
  font-weight: 700;
  color: var(--color-label);
  letter-spacing: -0.03em;
  line-height: 1.2;
}

.page-subtitle {
  font-size: 13px;
  color: var(--color-label-tertiary);
  margin-top: 4px;
}

.filter-card {
  margin-bottom: 16px;
}

.filter-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.filter-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-card :deep(.el-table) {
  border-radius: 0;
}

.clickable-table :deep(.el-table__row) {
  cursor: pointer;
}

.table-card :deep(.el-pagination) {
  padding: 16px 20px;
  border-top: 1px solid var(--color-separator);
}

.code-cell {
  color: var(--color-accent);
  font-size: 13px;
}

.url-cell {
  color: var(--color-label-secondary);
  font-size: 13px;
  text-decoration: none;
  display: block;
  max-width: 220px;
}

.url-cell:hover {
  color: var(--color-accent);
}

.text-secondary {
  color: var(--color-label-secondary);
  font-size: 13px;
}

.row-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-label-secondary);
  cursor: pointer;
  transition: background var(--duration-fast) var(--ease-out),
              color var(--duration-fast) var(--ease-out);
}

.action-btn:hover {
  background: rgba(0, 0, 0, 0.06);
  color: var(--color-label);
}

.action-btn--danger:hover {
  background: rgba(255, 59, 48, 0.08);
  color: var(--color-destructive);
}
</style>
